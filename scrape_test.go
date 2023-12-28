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
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
)

type Adsource struct {
	uri     string
	content string
}

func IoRead(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func InitSources(conf *Config) []Adsource {
	ads := []Adsource{
		{
			uri:     fmt.Sprintf("%s%s?userId=%d", Baseuri, Listuri, conf.User),
			content: IoRead("t/adlist1.html"),
		},
		{
			uri:     fmt.Sprintf("%s%s?userId=%d&pageNum=2", Baseuri, Listuri, conf.User),
			content: IoRead("t/adlist2.html"),
		},
		{
			uri:     fmt.Sprintf("%s%s?userId=%d&pageNum=3", Baseuri, Listuri, conf.User),
			content: IoRead("t/adlist3.html"),
		},
	}

	for n := 1; n < 10; n++ {
		ads = append(ads, Adsource{
			uri:     fmt.Sprintf("%s/s-anzeige/ad%d/%d", Baseuri, n, n),
			content: IoRead(fmt.Sprintf("t/ad%d.html", n)), // FIXME: add actual ad?.html files
		})
	}

	return ads
}

func SetIntercept(conf *Config) {
	ads := InitSources(conf)

	for _, ad := range ads {
		httpmock.RegisterResponder("GET", ad.uri,
			httpmock.NewStringResponder(200, ad.content))
	}

}

func TestStart(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	conf := &Config{User: 1, Outdir: "t/out"}

	SetIntercept(conf)

	if err := Start(conf); err != nil {
		t.Errorf("failed: %s", err.Error())
	}
}
