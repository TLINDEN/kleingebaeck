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

	// fmt.Println(uri)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// extract links from  all ad listing pages (that  is: use pagination)
// and scrape every page
func Start(uid string, dir string) error {
	client := &http.Client{}
	ads := []string{}

	baseuri := Baseuri + Listuri + "?userId=" + uid
	page := 1
	uri := baseuri

	for {
		var index Index
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

		for _, href := range index.Links {
			ads = append(ads, href)
			fmt.Println(href)
		}

		page++
		uri = baseuri + "&pageNum=" + fmt.Sprintf("%d", page)
	}

	for _, ad := range ads {
		err := Scrape(ad, dir)
		if err != nil {
			return err
		}
	}

	return nil
}

type Ad struct {
	Title  string   `goquery:"h1"`
	Text   string   `goquery:"p#viewad-description-text,html"`
	Images []string `goquery:".galleryimage-element img,[src]"`
	Price  string   `goquery:"h2#viewad-price"`
}

func Scrape(link string, dir string) error {
	client := &http.Client{}
	uri := Baseuri + link
	slurp := strings.Split(uri, "/")[1]

	var ad Ad
	body, err := Get(uri, client)
	if err != nil {
		return err
	}
	defer body.Close()

	err = goq.NewDecoder(body).Decode(&ad)
	if err != nil {
		return err
	}

	f, err := os.Create(strings.Join([]string{dir, slurp, "Anzeige.txt"}, "/"))
	if err != nil {
		return err
	}

	ad.Text = strings.ReplaceAll(ad.Text, "<br/>", "\n")
	_, err = fmt.Fprintf(f, "Title: %s\nPrice: %s\n\n%s", ad.Title, ad.Price, ad.Text)
	if err != nil {
		return err
	}

	img := 1
	for _, imguri := range ad.Images {
		file := fmt.Sprintf("%s/%d.jpg", dir, img)
		err := Getimage(imguri, file)
		if err != nil {
			return err
		}

		img++
	}

	return nil
}

// fetch an image
func Getimage(uri, fileName string) error {
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
