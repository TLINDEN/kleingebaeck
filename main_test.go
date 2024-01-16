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
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	tpl "text/template"

	"github.com/jarcoal/httpmock"
)

// the ad list, aka:
// https://www.kleinanzeigen.de/s-bestandsliste.html?userId=XXXXXX
// Note, that this HTML code is reduced to the max, so that it only
// contains the stuff required to satisfy goquery
const LISTTPL string = `<!DOCTYPE html>
<html lang="de" >
  <head>
    <title>Ads</title>
  </head>
  <body>
{{ range . }}
     <h2 class="text-module-begin">
        <a class="ellipsis"
           href="/s-anzeige/{{ .Slug }}/{{ .Id }}">{{ .Title }}</a>
     </h2>
{{ end }}
  </body>
</html>
`

// an actual ad listing, aka:
// https://www.kleinanzeigen.de/s-anzeige/ad-text-slug/1010101010
// Note, that this HTML code is reduced to the max, so that it only
// contains the stuff required to satisfy goquery
const ADTPL string = `DOCTYPE html>
<html lang="de">
  <head>
    <title>Ad Listing</title>
  </head>
  <body>

    <div class="l-container-row">
        <div id="vap-brdcrmb" class="breadcrump">
            <a class="breadcrump-link" itemprop="url" href="/" title="Kleinanzeigen ">
                <span itemprop="title">Kleinanzeigen </span>
            </a>
            <a class="breadcrump-link" itemprop="url" href="/egal">
               <span itemprop="title">{{ .Category }}</span></a>
            </div>
    </div>

    {{ range $image := .Images }}
    <div class="galleryimage-element" data-ix="3">
      <img src="{{ $image }}"/>
    </div>
    {{ end }}

    <h1 id="viewad-title" class="boxedarticle--title" itemprop="name" data-soldlabel="Verkauft">
      {{ .Title }}</h1>
    <div class="boxedarticle--flex--container">
      <h2 class="boxedarticle--price" id="viewad-price">
        {{ .Price }}</h2>
    </div>

    <div id="viewad-extra-info" class="boxedarticle--details--full">
      <div><i class="icon icon-small icon-calendar-gray-simple"></i><span>{{ .Created }}</span></div>
    </div>

    <div class="splitlinebox l-container-row" id="viewad-details">
      <ul class="addetailslist">
        <li class="addetailslist--detail">
          Zustand<span class="addetailslist--detail--value" >
          {{ .Condition }}</span>
        </li>
      </ul>
    </div>

    <div class="l-container last-paragraph-no-margin-bottom">
      <p id="viewad-description-text" class="text-force-linebreak " itemprop="description">
        {{ .Text }}
      </p>
    </div>
  </body>
</html>
`

const EMPTYPAGE string = `DOCTYPE html>
<html lang="de">
  <head></head>
  <body></body>
</html>
`

const (
	EMPTYURI       string = `https://www.kleinanzeigen.de/s-anzeige/empty/1`
	INVALID503URI  string = `https://www.kleinanzeigen.de/s-anzeige/503/1`
	INVALIDPATHURI string = `https://www.kleinanzeigen.de/anzeige/name/1`
	INVALID404URI  string = `https://www.kleinanzeigen.de/anzeige/name/1/foo/bar`
	INVALIDURI     string = `https://foo.bar/weird/things`
)

var base = "kleingebaeck -c t/config-empty.conf"

type Tests struct {
	name     string
	args     string
	expect   string
	exitcode int
}

var tests = []Tests{
	{
		name:     "version",
		args:     base + " -V",
		expect:   "This is",
		exitcode: 0,
	},
	{
		name:     "help",
		args:     base + " -h",
		expect:   "Usage:",
		exitcode: 0,
	},
	{
		name:     "debug",
		args:     base + " -d",
		expect:   "program_info",
		exitcode: 1,
	},
	{
		name:     "no-args-no-user",
		args:     base,
		expect:   "invalid or no user id",
		exitcode: 1,
	},
	{
		name:     "download-single-ad",
		args:     base + " -o t/out https://www.kleinanzeigen.de/s-anzeige/first-ad/1",
		expect:   "Successfully downloaded 1 ad with 2 images to t/out",
		exitcode: 0,
	},
	{
		name:     "download-single-ad-verbose",
		args:     base + " -o t/out https://www.kleinanzeigen.de/s-anzeige/first-ad/1 -v",
		expect:   "wrote ad listing",
		exitcode: 0,
	},
	{
		name:     "download-single-ad-debug",
		args:     base + " -o t/out https://www.kleinanzeigen.de/s-anzeige/first-ad/1 -d",
		expect:   "extracted ad listing program_info.pid=",
		exitcode: 0,
	},
	{
		name:     "download-all-ads",
		args:     base + " -o t/out -u 1",
		expect:   "Successfully downloaded 6 ads with 12 images to t/out",
		exitcode: 0,
	},
	{
		name:     "download-all-ads-using-config",
		args:     "kleingebaeck -c t/fullconfig.conf",
		expect:   "Successfully downloaded 6 ads with 12 images to t/out",
		exitcode: 0,
	},
}

var invalidtests = []Tests{
	{
		name:     "empty-ad",
		args:     base + " " + EMPTYURI,
		expect:   "could not extract ad data from page, got empty struct",
		exitcode: 1,
	},
	{
		name:     "invalid-ad",
		args:     base + " " + INVALIDURI,
		expect:   "invalid uri",
		exitcode: 1,
	},
	{
		name:     "invalid-path",
		args:     base + " " + INVALIDPATHURI,
		expect:   "could not extract ad data from page, got empty struct",
		exitcode: 1,
	},
	{
		name:     "404",
		args:     base + " " + INVALID404URI,
		expect:   "could not get page via HTTP",
		exitcode: 1,
	},
	{
		name:     "outdir-no-exists",
		args:     base + " -o t/foo/bar/out https://www.kleinanzeigen.de/s-anzeige/first-ad/1 -v",
		expect:   "Failure",
		exitcode: 1,
	},
	{
		name:     "wrong-flag",
		args:     base + " -X",
		expect:   "unknown shorthand flag: 'X' in -X",
		exitcode: 1,
	},
	{
		name:     "no-config",
		args:     "kleingebaeck -c t/invalid.conf",
		expect:   "error loading config file",
		exitcode: 1,
	},
	{
		name:     "503",
		args:     base + " " + INVALID503URI,
		expect:   "could not get page via HTTP",
		exitcode: 1,
	},
}

type AdConfig struct {
	Title     string
	Slug      string
	Id        string
	Price     string
	Category  string
	Condition string
	Created   string
	Text      string
	Images    []string // files in ./t/
}

var adsrc = []AdConfig{
	{
		Title: "First Ad",
		Id:    "1", Price: "5€",
		Category:  "Klimbim",
		Text:      "Thing to sale",
		Slug:      "first-ad",
		Condition: "works",
		Created:   "Yesterday",
		Images:    []string{"t/1.jpg", "t/2.jpg"},
	},
	{
		Title: "Secnd Ad",
		Id:    "2", Price: "5€",
		Category:  "Kram",
		Text:      "Thing to sale",
		Slug:      "second-ad",
		Condition: "works",
		Created:   "Yesterday",
		Images:    []string{"t/1.jpg", "t/2.jpg"},
	},
	{
		Title:     "Third Ad",
		Id:        "3",
		Price:     "5€",
		Category:  "Kuddelmuddel",
		Text:      "Thing to sale",
		Slug:      "third-ad",
		Condition: "works",
		Created:   "Yesterday",
		Images:    []string{"t/1.jpg", "t/2.jpg"},
	},
	{
		Title:     "Forth Ad",
		Id:        "4",
		Price:     "5€",
		Category:  "Krempel",
		Text:      "Thing to sale",
		Slug:      "fourth-ad",
		Condition: "works",
		Created:   "Yesterday",
		Images:    []string{"t/1.jpg", "t/2.jpg"},
	},
	{
		Title:     "Fifth Ad",
		Id:        "5",
		Price:     "5€",
		Category:  "Kladderadatsch",
		Text:      "Thing to sale",
		Slug:      "fifth-ad",
		Condition: "works",
		Created:   "Yesterday",
		Images:    []string{"t/1.jpg", "t/2.jpg"},
	},
	{
		Title:     "Sixth Ad",
		Id:        "6",
		Price:     "5€",
		Category:  "Klunker",
		Text:      "Thing to sale",
		Slug:      "sixth-ad",
		Condition: "works",
		Created:   "Yesterday",
		Images:    []string{"t/1.jpg", "t/2.jpg"},
	},
}

// An  Adsource  is used  to  construct  a  httpmock responder  for  a
// particular    url.    So,     the    code    (scrape.go)    scrapes
// https://kleinanzeigen.de,  but  in  reality httpmock  captures  the
// request and responds with our mock data
type Adsource struct {
	uri     string
	content string
	status  int
}

// Render a HTML template for an adlisting or an ad
func GetTemplate(l []AdConfig, a AdConfig, htmltemplate string) string {
	tmpl, err := tpl.New("template").Parse(htmltemplate)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	if len(a.Id) == 0 {
		err = tmpl.Execute(&out, l)
	} else {
		err = tmpl.Execute(&out, a)
	}

	if err != nil {
		panic(err)
	}

	return out.String()
}

// Initialize the valid sources for the httpmock responder
func InitValidSources() []Adsource {
	// valid ad listing page 1
	list1 := []AdConfig{
		adsrc[0], adsrc[1], adsrc[2],
	}

	// valid ad listing page 2
	list2 := []AdConfig{
		adsrc[3], adsrc[4], adsrc[5],
	}

	// valid ad listing page 3, which is empty
	list3 := []AdConfig{}

	// used to signal GetTemplate() to render a listing
	empty := AdConfig{}

	// prepare urls for the listing pages
	ads := []Adsource{
		{
			uri:     fmt.Sprintf("%s%s?userId=1", Baseuri, Listuri),
			content: GetTemplate(list1, empty, LISTTPL),
		},
		{
			uri:     fmt.Sprintf("%s%s?userId=1&pageNum=2", Baseuri, Listuri),
			content: GetTemplate(list2, empty, LISTTPL),
		},
		{
			uri:     fmt.Sprintf("%s%s?userId=1&pageNum=3", Baseuri, Listuri),
			content: GetTemplate(list3, empty, LISTTPL),
		},
	}

	// prepare urls for the ads
	for _, ad := range adsrc {
		ads = append(ads, Adsource{
			uri:     fmt.Sprintf("%s/s-anzeige/%s/%s", Baseuri, ad.Slug, ad.Id),
			content: GetTemplate(nil, ad, ADTPL),
		})
		//panic(GetTemplate(nil, ad, ADTPL))
	}

	return ads
}

func InitInvalidSources() []Adsource {
	empty := AdConfig{}
	ads := []Adsource{
		{
			// valid ad page but without content
			uri:     fmt.Sprintf("%s/s-anzeige/empty/1", Baseuri),
			content: GetTemplate(nil, empty, EMPTYPAGE),
		},
		{
			// some random foreign webpage
			uri:     INVALIDURI,
			content: GetTemplate(nil, empty, "<html>foo</html>"),
		},
		{
			// some invalid page path
			uri:     fmt.Sprintf("%s/anzeige/name/1", Baseuri),
			content: GetTemplate(nil, empty, "<html></html>"),
		},
		{
			// some none-ad page
			uri:     fmt.Sprintf("%s/anzeige/name/1/foo/bar", Baseuri),
			content: GetTemplate(nil, empty, "<html>HTTP 404: /eine-anzeige/ does not exist!</html>"),
			status:  404,
		},
		{
			// valid ad page but 503
			uri:     fmt.Sprintf("%s/s-anzeige/503/1", Baseuri),
			content: GetTemplate(nil, empty, "<html>HTTP 503: service unavailable</html>"),
			status:  503,
		},
	}

	return ads
}

// load a test image from disk
func GetImage(path string) []byte {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return dat
}

// setup httpmock
func SetIntercept(ads []Adsource) {
	for _, ad := range ads {
		if ad.status == 0 {
			ad.status = 200
		}

		httpmock.RegisterResponder("GET", ad.uri,
			httpmock.NewStringResponder(ad.status, ad.content))
	}

	// we just use 2 images, put this here
	for _, image := range []string{"t/1.jpg", "t/2.jpg"} {
		httpmock.RegisterResponder("GET", image,
			httpmock.NewBytesResponder(200, GetImage(image)))
	}

}

func VerifyAd(ad AdConfig) error {
	body := ad.Title + ad.Price + ad.Id + "Kleinanzeigen => " +
		ad.Category + ad.Condition + ad.Created

	// prepare ad dir name using DefaultAdNameTemplate
	c := Config{Adnametemplate: "{{ .Slug }}"}
	adstruct := Ad{Slug: ad.Slug, Id: ad.Id}
	addir, err := AdDirName(&c, &adstruct)
	if err != nil {
		return err
	}

	file := fmt.Sprintf("t/out/%s/Adlisting.txt", addir)
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if body != strings.TrimSpace(string(content)) {
		msg := fmt.Sprintf("ad content doesn't match.\nExpect: %s\n   Got: %s\n", body, content)
		return errors.New(msg)
	}

	return nil
}

func TestMain(t *testing.T) {
	oldargs := os.Args
	defer func() { os.Args = oldargs }()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// prepare httpmock responders
	SetIntercept(InitValidSources())

	// run commandline tests
	for _, tt := range tests {
		var buf bytes.Buffer
		os.Args = strings.Split(tt.args, " ")

		ret := Main(&buf)

		if ret != tt.exitcode {
			t.Errorf("%s with cmd <%s> did not exit with %d but %d",
				tt.name, tt.args, tt.exitcode, ret)
		}

		if !strings.Contains(buf.String(), tt.expect) {
			t.Errorf("%s with cmd <%s> output did not match.\nExpect: %s\n   Got: %s\n",
				tt.name, tt.args, tt.expect, buf.String())
		}
	}

	// verify if downloaded ads match
	for _, ad := range adsrc {
		if err := VerifyAd(ad); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestMainInvalids(t *testing.T) {
	oldargs := os.Args
	defer func() { os.Args = oldargs }()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// prepare httpmock responders
	SetIntercept(InitInvalidSources())

	// run commandline tests
	for _, tt := range invalidtests {
		var buf bytes.Buffer
		os.Args = strings.Split(tt.args, " ")

		ret := Main(&buf)

		if ret != tt.exitcode {
			t.Errorf("%s with cmd <%s> did not exit with %d but %d",
				tt.name, tt.args, tt.exitcode, ret)
		}

		if !strings.Contains(buf.String(), tt.expect) {
			t.Errorf("%s with cmd <%s> output did not match.\nExpect: %s\n   Got: %s\n",
				tt.name, tt.args, tt.expect, buf.String())
		}
	}
}
