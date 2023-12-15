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
	"os"
	"strings"

	"net/http"

	"astuart.co/goq"
)

type Index struct {
	Links []string `goquery:".text-module-begin a,[href]"`
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
func Start(uid string, dir string) error {
	client := &http.Client{}
	adlinks := []string{}

	baseuri := Baseuri + Listuri + "?userId=" + uid
	page := 1
	uri := baseuri

	slog.Info("fetching ad pages", "user", uid)

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

	for _, adlink := range adlinks {
		err := Scrape(Baseuri+adlink, dir)
		if err != nil {
			return err
		}
	}

	return nil
}

type Ad struct {
	Title  string `goquery:"h1"`
	Slug   string
	Id     string
	Text   string   `goquery:"p#viewad-description-text,html"`
	Images []string `goquery:".galleryimage-element img,[src]"`
	Price  string   `goquery:"h2#viewad-price"`
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

// scrape an ad. uri is the full uri of the ad, dir is the basedir
func Scrape(uri string, dir string) error {
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
	slog.Debug("extracted ad listing", "ad", ad)

	// prepare output dir
	dir = dir + "/" + ad.Slug
	err = Mkdir(dir)
	if err != nil {
		return err
	}

	// write ad file
	listingfile := strings.Join([]string{dir, "Adlisting.txt"}, "/")
	f, err := os.Create(listingfile)
	if err != nil {
		return err
	}

	ad.Text = strings.ReplaceAll(ad.Text, "<br/>", "\n")
	_, err = fmt.Fprintf(f, "Title: %s\nPrice: %s\nId: %s\nBody:\n\n%s\n",
		ad.Title, ad.Price, ad.Id, ad.Text)
	if err != nil {
		return err
	}
	slog.Info("wrote ad listing", "listingfile", listingfile)

	// fetch images
	img := 1
	for _, imguri := range ad.Images {
		file := fmt.Sprintf("%s/%d.jpg", dir, img)
		err := Getimage(imguri, file)
		if err != nil {
			return err
		}
		slog.Info("wrote ad image", "image", file)

		img++
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

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
