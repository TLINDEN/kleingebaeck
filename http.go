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
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// I add an artificial "ID" to each HTTP request and the corresponding
// respose for  debugging purposes  so that  the pair  of them  can be
// easier associated in debug output
var letters = []rune("ABCDEF0123456789")

const IDLEN int = 8

// retry after HTTP 50x errors or err!=nil
const RetryCount = 3

func getid() string {
	b := make([]rune, IDLEN)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// used to inject debug log and implement retries
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

// the actual logging transport with retries
func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// just required for debugging
	requestid := getid()

	// clone the request body, put into request on retry
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	slog.Debug("REQUEST", "id", requestid, "uri", req.URL, "host", req.Host)

	// first try
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err == nil {
		slog.Debug("RESPONSE", "id", requestid, "status", resp.StatusCode,
			"contentlength", resp.ContentLength)
	}

	// enter retry check and loop, if first req were successful, leave loop immediately
	retries := 0
	for shouldRetry(err, resp) && retries < RetryCount {
		time.Sleep(backoff(retries))

		// consume any response to reuse the connection.
		drainBody(resp)

		// clone the request body again
		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// actual retry
		resp, err = http.DefaultTransport.RoundTrip(req)

		if err == nil {
			slog.Debug("RESPONSE", "id", requestid, "status", resp.StatusCode,
				"contentlength", resp.ContentLength, "retry", retries)
		}

		retries++
	}

	if err != nil {
		return resp, fmt.Errorf("failed to get HTTP response for %s: %w", req.URL, err)
	}

	return resp, nil
}
