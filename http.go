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

// FIXME: we could also incorporate
// https://github.com/kdkumawat/golang/blob/main/http-retry/http/retry-client.go

package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type loggingTransport struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("ABCDEF0123456789")

func getid() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)

	// just requred for debugging
	id := getid()
	slog.Debug("REQUEST", "id", id, "uri", req.URL, "host", req.Host)
	slog.Debug("RESPONSE", "id", id, "status", resp.StatusCode, "contentlength", resp.ContentLength)

	if len(os.Getenv("DEBUGHTTP")) > 0 {
		fmt.Println("DEBUGHTTP Request ===>")
		bytes, _ := httputil.DumpRequestOut(req, true)
		fmt.Printf("%s\n", bytes)

		fmt.Println("<=== DEBUGHTTP Response")
		for header, value := range resp.Header {
			fmt.Printf("%s: %s\n", header, value)
		}
		fmt.Printf("Status: %s %s\nContent-Length: %d\n\n\n", resp.Proto, resp.Status, resp.ContentLength)

	}

	return resp, err
}
