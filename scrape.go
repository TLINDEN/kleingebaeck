/*
Copyright © 2023-2024 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"astuart.co/goq"
	"golang.org/x/sync/errgroup"
)

// extract links from  all ad listing pages (that  is: use pagination)
// and scrape every page
func ScrapeUser(fetch *Fetcher) error {
	adlinks := []string{}

	baseuri := fmt.Sprintf("%s%s?userId=%d", Baseuri, Listuri, fetch.Config.User)
	page := 1
	uri := baseuri

	slog.Info("fetching ad pages", "user", fetch.Config.User)

	for {
		var index Index

		slog.Debug("fetching page", "uri", uri)

		body, err := fetch.Get(uri)
		if err != nil {
			return err
		}
		defer body.Close()

		err = goq.NewDecoder(body).Decode(&index)
		if err != nil {
			return fmt.Errorf("failed to goquery decode HTML index body: %w", err)
		}

		if len(index.Links) == 0 {
			break
		}

		slog.Debug("extracted ad links", "count", len(index.Links))

		for _, href := range index.Links {
			adlinks = append(adlinks, href)
			slog.Debug("ad link", "href", href)
		}

		page++
		uri = baseuri + "&pageNum=" + strconv.Itoa(page)
	}

	for index, adlink := range adlinks {
		err := ScrapeAd(fetch, Baseuri+adlink)
		if err != nil {
			return err
		}

		if fetch.Config.Limit > 0 && index == fetch.Config.Limit-1 {
			break
		}
	}

	return nil
}

// scrape an ad. uri is the full uri of the ad, dir is the basedir
func ScrapeAd(fetch *Fetcher, uri string) error {
	advertisement := &Ad{}

	// extract slug and id from uri
	uriparts := strings.Split(uri, "/")
	if len(uriparts) < SlugURIPartNum {
		return fmt.Errorf("invalid uri: %s", uri)
	}

	advertisement.Slug = uriparts[4]
	advertisement.ID = uriparts[5]

	// get the ad
	slog.Debug("fetching ad page", "uri", uri)

	body, err := fetch.Get(uri)
	if err != nil {
		return err
	}
	defer body.Close()

	// extract ad contents with goquery/goq
	err = goq.NewDecoder(body).Decode(&advertisement)
	if err != nil {
		return fmt.Errorf("failed to goquery decode HTML ad body: %w", err)
	}

	if len(advertisement.CategoryTree) > 0 {
		advertisement.Category = strings.Join(advertisement.CategoryTree, " => ")
	}

	if advertisement.Incomplete() {
		slog.Debug("got ad", "ad", advertisement)

		return fmt.Errorf("could not extract ad data from page, got empty struct")
	}

	advertisement.CalculateExpire()

	// prepare ad dir name
	addir, err := AdDirName(fetch.Config, advertisement)
	if err != nil {
		return err
	}

	proceed := CheckAdVisited(fetch.Config, addir)
	if !proceed {
		return nil
	}

	// write listing
	err = WriteAd(fetch.Config, advertisement, addir)
	if err != nil {
		return err
	}

	// tell the user
	slog.Debug("extracted ad listing", "ad", advertisement)

	// stats
	fetch.Config.IncrAds()

	// register for later checks
	DirsVisited[addir] = 1

	return ScrapeImages(fetch, advertisement, addir)
}

func ScrapeImages(fetch *Fetcher, advertisement *Ad, addir string) error {
	// fetch images
	img := 1
	adpath := filepath.Join(fetch.Config.Outdir, addir)

	// scan existing images, if any
	cache, err := ReadImages(adpath, fetch.Config.ForceDownload)
	if err != nil {
		return err
	}

	egroup := new(errgroup.Group)

	for _, imguri := range advertisement.Images {
		imguri := imguri

		// we append the suffix later in NewImage() based on image format
		basefilename := filepath.Join(adpath, fmt.Sprintf("%d", img))

		egroup.Go(func() error {
			// wait a little

			throttle := GetThrottleTime()
			time.Sleep(throttle)

			body, err := fetch.Getimage(imguri)
			if err != nil {
				return err
			}

			buf := new(bytes.Buffer)

			_, err = buf.ReadFrom(body)
			if err != nil {
				return fmt.Errorf("failed to read from image buffer: %w", err)
			}

			reader := bytes.NewReader(buf.Bytes())

			image, err := NewImage(reader, basefilename, imguri)
			if err != nil {
				return err
			}

			err = image.CalcHash()
			if err != nil {
				return err
			}

			if !fetch.Config.ForceDownload {
				if image.SimilarExists(cache) {
					slog.Debug("similar image exists, not written", "uri", image.URI)

					return nil
				}
			}

			_, err = reader.Seek(0, 0)
			if err != nil {
				return fmt.Errorf("failed to seek(0) on image reader: %w", err)
			}

			err = WriteImage(image.Filename, reader)
			if err != nil {
				return err
			}

			slog.Debug("wrote image", "image", image, "size", buf.Len(), "throttle", throttle)

			return nil
		})
		img++
	}

	if err := egroup.Wait(); err != nil {
		return fmt.Errorf("failed to finalize error waitgroup: %w", err)
	}

	fetch.Config.IncrImgs(len(advertisement.Images))

	return nil
}
