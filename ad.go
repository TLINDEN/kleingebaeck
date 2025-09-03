/*
Copyright © 2023-2025 Thomas von Dein

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
	"bufio"
	"log/slog"
	"strings"
	"time"
)

type Index struct {
	Links []string `goquery:".text-module-begin a,[href]"`
}

type Ad struct {
	Title        string `goquery:"h1"`
	Slug         string
	ID           string
	Details      string            `goquery:".addetailslist--detail,text"`
	Attributes   map[string]string // processed afterwards
	Condition    string            // post processed from details for backward compatibility
	Type         string            // post processed from details for backward compatibility
	Color        string            // post processed from details for backward compatibility
	Material     string            // post processed from details for backward compatibility
	Category     string
	CategoryTree []string `goquery:".breadcrump-link,text"`
	Price        string   `goquery:"h2#viewad-price"`
	Created      string   `goquery:"#viewad-extra-info,text"`
	Text         string   `goquery:"p#viewad-description-text,html"`
	Images       []string `goquery:".galleryimage-element img,[src]"`
	Shipping     string   `goquery:".boxedarticle--details--shipping,text"` // not always filled
	Expire       string

	// runtime computed
	Year, Day, Month string
}

// Used by slog to pretty print an ad
func (ad *Ad) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", ad.Title),
		slog.String("price", ad.Price),
		slog.String("id", ad.ID),
		slog.Int("imagecount", len(ad.Images)),
		slog.Int("bodysize", len(ad.Text)),
		slog.String("categorytree", strings.Join(ad.CategoryTree, "+")),
		slog.String("created", ad.Created),
		slog.String("expire", ad.Expire),
		slog.String("shipping", ad.Shipping),
		slog.String("details", ad.Details),
	)
}

// check for  completeness.  I  erected these  fields to  be mandatory
// (though I really don't know if  they really are). I consider images
// and meta  optional. So,  if either  of the  checked fields  here is
// empty we  return an  error.  All the  checked fields  are extracted
// using goquery. However,  I think price is optional  since there are
// ads for gifts as well.
//
// Note: we return true for "ad is incomplete" and false for "ad is complete"!
func (ad *Ad) Incomplete() bool {
	if ad.Category == "" || ad.Created == "" || ad.Text == "" {
		return true
	}

	return false
}

func (ad *Ad) CalculateExpire() {
	if ad.Created != "" {
		ts, err := time.Parse("02.01.2006", ad.Created)
		if err == nil {
			ad.Expire = ts.AddDate(0, 0, ExpireDays).Format("02.01.2006")
		}
	}
}

/*
Decode attributes like color or condition. See
https://github.com/TLINDEN/kleingebaeck/issues/117
for more details. In short: the HTML delivered by
kleinanzeigen.de has no css attribute for the keys
so we cannot extract key=>value mappings of the
ad details but have to parse them manually.

The ad.Details member contains this after goq run:

Art

	Weitere Kinderzimmermöbel

	Farbe
	Holz

	Zustand
	In Ordnung

We parse this into ad.Attributes and fill in some
static members for backward compatibility reasons.
*/
func (ad *Ad) DecodeAttributes() {
	rd := strings.NewReader(ad.Details)
	scanner := bufio.NewScanner(rd)

	isattr := true
	attr := ""
	attrmap := map[string]string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if isattr {
			attr = line
		} else {
			attrmap[attr] = line
		}

		isattr = !isattr
	}

	ad.Attributes = attrmap

	if Exists(ad.Attributes, "Zustand") {
		ad.Condition = ad.Attributes["Zustand"]
	}

	if Exists(ad.Attributes, "Farbe") {
		ad.Color = ad.Attributes["Farbe"]
	}

	if Exists(ad.Attributes, "Art") {
		ad.Type = ad.Attributes["Art"]
	}

	if Exists(ad.Attributes, "Material") {
		ad.Material = ad.Attributes["Material"]
	}

	slog.Debug("parsed attributes", "attributes", ad.Attributes)

	ad.Shipping = strings.Replace(ad.Shipping, "+ Versand ab ", "", 1)
}
