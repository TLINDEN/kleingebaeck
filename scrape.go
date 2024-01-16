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
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"astuart.co/goq"
	"golang.org/x/sync/errgroup"
)

// fetch some web page content
func Get(uri string, client *http.Client) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", Useragent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	slog.Debug("response", "code", res.StatusCode, "status",
		res.Status, "size", res.ContentLength)

	if res.StatusCode != 200 {
		return nil, errors.New("could not get page via HTTP")
	}

	return res.Body, nil
}

// extract links from  all ad listing pages (that  is: use pagination)
// and scrape every page
func ScrapeUser(conf *Config, client *http.Client) error {
	adlinks := []string{}

	baseuri := fmt.Sprintf("%s%s?userId=%d", Baseuri, Listuri, conf.User)
	page := 1
	uri := baseuri

	slog.Info("fetching ad pages", "user", conf.User)

	for {
		var index Index
		slog.Debug("fetching page", "uri", uri)
		body, err := Get(uri, client)
		if err != nil {
			return err
		}
		defer body.Close()

		err = goq.NewDecoder(body).Decode(&index)
		if err != nil {
			return err
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
		uri = baseuri + "&pageNum=" + fmt.Sprintf("%d", page)
	}

	for i, adlink := range adlinks {
		err := ScrapeAd(conf, Baseuri+adlink, client)
		if err != nil {
			return err
		}

		if conf.Limit > 0 && i == conf.Limit-1 {
			break
		}
	}

	return nil
}

// scrape an ad. uri is the full uri of the ad, dir is the basedir
func ScrapeAd(c *Config, uri string, client *http.Client) error {
	ad := &Ad{}

	// extract slug and id from uri
	uriparts := strings.Split(uri, "/")
	if len(uriparts) < 6 {
		return errors.New("invalid uri: " + uri)
	}
	ad.Slug = uriparts[4]
	ad.Id = uriparts[5]

	// get the ad
	slog.Debug("fetching ad page", "uri", uri)
	body, err := Get(uri, client)
	if err != nil {
		return err
	}
	defer body.Close()

	// extract ad contents with goquery/goq
	err = goq.NewDecoder(body).Decode(&ad)
	if err != nil {
		return err
	}

	if len(ad.CategoryTree) > 0 {
		ad.Category = strings.Join(ad.CategoryTree, " => ")
	}

	if ad.Incomplete() {
		slog.Debug("got ad", "ad", ad)
		return errors.New("could not extract ad data from page, got empty struct")
	}

	slog.Debug("extracted ad listing", "ad", ad)

	// write listing
	addir, err := WriteAd(c, ad)
	if err != nil {
		return err
	}

	c.IncrAds()

	return ScrapeImages(c, ad, addir, client)
}

func ScrapeImages(c *Config, ad *Ad, addir string, client *http.Client) error {
	// fetch images
	img := 1
	g := new(errgroup.Group)

	for _, imguri := range ad.Images {
		imguri := imguri
		file := filepath.Join(c.Outdir, addir, fmt.Sprintf("%d.jpg", img))
		g.Go(func() error {
			err := Getimage(c, imguri, file, client)
			if err != nil {
				return err
			}

			return nil
		})
		img++
	}

	if err := g.Wait(); err != nil {
		return err
	}

	c.IncrImgs(len(ad.Images))

	return nil
}

// fetch an image
func Getimage(c *Config, uri, fileName string, client *http.Client) error {
	slog.Debug("fetching ad image", "uri", uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		if c.IgnoreErrors {
			slog.Info("Failed to download image, error ignored", "error", err.Error())
		}
		return err
	}

	req.Header.Set("User-Agent", Useragent)

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("could not get image via HTTP")
	}

	err = WriteImage(fileName, response.Body)
	if err != nil {
		return err
	}

	slog.Info("wrote ad image", "image", fileName)
	return nil
}
