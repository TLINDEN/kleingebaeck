/*
Copyright © 2023 Thomas von Dein

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
	"sync"

	"astuart.co/goq"
)

type Index struct {
	Links []string `goquery:".text-module-begin a,[href]"`
}

type Ad struct {
	Title     string `goquery:"h1"`
	Slug      string
	Id        string
	Condition string
	Category  string
	Price     string   `goquery:"h2#viewad-price"`
	Created   string   `goquery:"#viewad-extra-info,text"`
	Text      string   `goquery:"p#viewad-description-text,html"`
	Images    []string `goquery:".galleryimage-element img,[src]"`
	Meta      []string `goquery:".addetailslist--detail--value,text"`
}

func (ad *Ad) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", ad.Title),
		slog.String("price", ad.Price),
		slog.String("id", ad.Id),
		slog.Int("imagecount", len(ad.Images)),
		slog.Int("bodysize", len(ad.Text)),
	)
}

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

	return res.Body, nil
}

// extract links from  all ad listing pages (that  is: use pagination)
// and scrape every page
func Start(conf *Config) error {
	client := &http.Client{}
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
		err := Scrape(Baseuri+adlink, conf.Outdir, conf.Template)
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
func Scrape(uri string, dir string, template string) error {
	client := &http.Client{}
	ad := &Ad{}

	// extract slug and id from uri
	uriparts := strings.Split(uri, "/")
	if len(uriparts) < 6 {
		return errors.New("invalid uri")
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
	if len(ad.Meta) == 2 {
		ad.Category = ad.Meta[0]
		ad.Condition = ad.Meta[1]
	}
	slog.Debug("extracted ad listing", "ad", ad)

	// write listing
	err = WriteAd(dir, ad, template)
	if err != nil {
		return err
	}

	return ScrapeImages(dir, ad)
}

func ScrapeImages(dir string, ad *Ad) error {
	// fetch images
	img := 1
	var wg sync.WaitGroup
	wg.Add(len(ad.Images))
	failure := make(chan string)

	for _, imguri := range ad.Images {
		imguri := imguri
		file := filepath.Join(dir, ad.Slug, fmt.Sprintf("%d.jpg", img))
		go func() {
			defer wg.Done()
			err := Getimage(imguri, file)
			if err != nil {
				failure <- err.Error()
				return
			}
			slog.Info("wrote ad image", "image", file)
		}()
		img++
	}

	close(failure)
	wg.Wait()
	goterr := <-failure

	if goterr != "" {
		return errors.New(goterr)
	}

	return nil
}

// fetch an image
func Getimage(uri, fileName string) error {
	slog.Debug("fetching ad image", "uri", uri)
	response, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	err = WriteImage(fileName, response.Body)
	if err != nil {
		return err
	}

	return nil
}
