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
	"io"
	"io/ioutil"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var letters = []rune("ABCDEF0123456789")

func getid() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

const RetryCount = 3

type loggingTransport struct{}

// escalating timeout, $retry^2 seconds
func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

// only retry in case of errors or certain non 200 HTTP codes
func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	if resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		return true
	}

	return false
}

// Body needs to be drained, otherwise we can't reuse the http.Response
func drainBody(resp *http.Response) {
	if resp != nil {
		if resp.Body != nil {
			_, err := io.Copy(io.Discard, resp.Body)
			if err != nil {
				// unable to copy data? uff!
				panic(err)
			}
			resp.Body.Close()
		}
	}
}

// our logging transport with retries
func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// just requred for debugging
	id := getid()

	// clone the request body, put into request on retry
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	slog.Debug("REQUEST", "id", id, "uri", req.URL, "host", req.Host)

	// first try
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err == nil {
		slog.Debug("RESPONSE", "id", id, "status", resp.StatusCode,
			"contentlength", resp.ContentLength)
	}

	// enter retry check and loop, if first req were successfull, leave loop immediately
	retries := 0
	for shouldRetry(err, resp) && retries < RetryCount {
		time.Sleep(backoff(retries))

		// consume any response to reuse the connection.
		drainBody(resp)

		// clone the request body again
		if req.Body != nil {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// actual retry
		resp, err = http.DefaultTransport.RoundTrip(req)

		if err == nil {
			slog.Debug("RESPONSE", "id", id, "status", resp.StatusCode,
				"contentlength", resp.ContentLength, "retry", retries)
		}

		retries++
	}

	return resp, err
}
