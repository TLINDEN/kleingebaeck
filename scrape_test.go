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
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	tpl "text/template"

	"github.com/jarcoal/httpmock"
)

// used to fill an ad template and the ad listing page template
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
          Art<span class="addetailslist--detail--value" >
          {{ .Category }}</span>
        </li>
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

// An  Adsource  is used  to  construct  a  httpmock responder  for  a
// particular    url.    So,     the    code    (scrape.go)    scrapes
// https://kleinanzeigen.de,  but  in  reality httpmock  captures  the
// request and responds with our mock data
type Adsource struct {
	uri     string
	content string
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

func InitAds() []AdConfig {
	return []AdConfig{
		{Title: "First Ad", Id: "1", Price: "5€", Category: "Klimbim", Text: "Thing to sale", Slug: "first-ad",
			Condition: "works", Created: "Yesterday", Images: []string{"t/1.jpg", "t/2.jpg"}},
		{Title: "Secnd Ad", Id: "2", Price: "5€", Category: "Kram", Text: "Thing to sale", Slug: "second-ad",
			Condition: "works", Created: "Yesterday", Images: []string{"t/1.jpg", "t/2.jpg"}},
		{Title: "Third Ad", Id: "3", Price: "5€", Category: "Kuddelmuddel", Text: "Thing to sale", Slug: "third-ad",
			Condition: "works", Created: "Yesterday", Images: []string{"t/1.jpg", "t/2.jpg"}},
		{Title: "Forth Ad", Id: "4", Price: "5€", Category: "Krempel", Text: "Thing to sale", Slug: "fourth-ad",
			Condition: "works", Created: "Yesterday", Images: []string{"t/1.jpg", "t/2.jpg"}},
		{Title: "Fifth Ad", Id: "5", Price: "5€", Category: "Kladderadatsch", Text: "Thing to sale", Slug: "fifth-ad",
			Condition: "works", Created: "Yesterday", Images: []string{"t/1.jpg", "t/2.jpg"}},
		{Title: "Sixth Ad", Id: "6", Price: "5€", Category: "Klunker", Text: "Thing to sale", Slug: "sixth-ad",
			Condition: "works", Created: "Yesterday", Images: []string{"t/1.jpg", "t/2.jpg"}},
	}
}

// Initialize the valid sources for the httpmock responder
func InitValidSources(conf *Config) []Adsource {
	// all our valid ads
	adsrc := InitAds()

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
			uri:     fmt.Sprintf("%s%s?userId=%d", Baseuri, Listuri, conf.User),
			content: GetTemplate(list1, empty, LISTTPL),
		},
		{
			uri:     fmt.Sprintf("%s%s?userId=%d&pageNum=2", Baseuri, Listuri, conf.User),
			content: GetTemplate(list2, empty, LISTTPL),
		},
		{
			uri:     fmt.Sprintf("%s%s?userId=%d&pageNum=3", Baseuri, Listuri, conf.User),
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

// load a test image from disk
func GetImage(path string) []byte {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return dat
}

// setup httpmock
func SetIntercept(conf *Config) {
	ads := InitValidSources(conf)

	for _, ad := range ads {
		httpmock.RegisterResponder("GET", ad.uri,
			httpmock.NewStringResponder(200, ad.content))
	}

	// we just use 2 images, put this here
	for _, image := range []string{"t/1.jpg", "t/2.jpg"} {
		httpmock.RegisterResponder("GET", image, httpmock.NewBytesResponder(200, GetImage(image)))
	}

}

// the  actual  test, calls  Start()  from  scrape, which  recursively
// scrapes ads from a user
func TestStart(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// fake config
	conf := &Config{User: 1, Outdir: "t/out", Template: DefaultTemplate}

	// prepare httpmock responders
	SetIntercept(conf)

	// run
	if err := Start(conf); err != nil {
		t.Errorf("failed to scrape: %s", err.Error())
	}

	// verify
	for _, ad := range InitAds() {
		file := fmt.Sprintf("t/out/%s/Adlisting.txt", ad.Slug)
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("failed to read adlisting: %s", err.Error())
		}

		if !strings.Contains(string(content), ad.Category) && !strings.Contains(string(content), ad.Title) {
			t.Errorf("failed to verify: %s content doesn't contain expected data", file)
		}
	}

	// uncomment to see slogs
	//t.Errorf("debug")
}
