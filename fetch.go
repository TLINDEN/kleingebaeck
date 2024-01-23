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
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// convenient wrapper to fetch some web content
type Fetcher struct {
	Config    *Config
	Client    *http.Client
	Useragent string // FIXME: make configurable
	Cookies   []*http.Cookie
}

func NewFetcher(c *Config) *Fetcher {
	jar, err := cookiejar.New(nil)
	if err != nil {
		// cannot return error here, FIXME
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	return &Fetcher{
		Client: &http.Client{
			Transport: &loggingTransport{}, // implemented in http.go
			Jar:       jar,
		},
		Useragent: Useragent, // default in config.go
		Config:    c,
		Cookies:   []*http.Cookie{},
	}
}

func (f *Fetcher) Get(uri string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", f.Useragent)

	if len(f.Cookies) > 0 {
		uriobj, _ := url.Parse(Baseuri)
		slog.Debug("have cookies, sending them",
			"sample-cookie-name", f.Cookies[0].Name,
			"sample-cookie-expire", f.Cookies[0].Expires,
		)
		f.Client.Jar.SetCookies(uriobj, f.Cookies)
	}

	res, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("could not get page via HTTP")
	}

	f.Cookies = res.Cookies()

	return res.Body, nil
}

// fetch an image
func (f *Fetcher) Getimage(uri string) (io.ReadCloser, error) {
	slog.Debug("fetching ad image", "uri", uri)
	body, err := f.Get(uri)
	if err != nil {
		if f.Config.IgnoreErrors {
			slog.Info("Failed to download image, error ignored", "error", err.Error())
			return nil, nil
		}
		return nil, err
	}

	return body, nil
}
