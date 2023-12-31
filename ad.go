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
	"log/slog"
)

type Index struct {
	Links []string `goquery:".text-module-begin a,[href]"`
}

type Ad struct {
	Title     string `goquery:"h1"`
	Slug      string
	Id        string
	Condition string
	Category  string
	Price     string   `goquery:"h2#viewad-price"`
	Created   string   `goquery:"#viewad-extra-info,text"`
	Text      string   `goquery:"p#viewad-description-text,html"`
	Images    []string `goquery:".galleryimage-element img,[src]"`
	Meta      []string `goquery:".addetailslist--detail--value,text"`
}

// Used by slog to pretty print an ad
func (ad *Ad) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", ad.Title),
		slog.String("price", ad.Price),
		slog.String("id", ad.Id),
		slog.Int("imagecount", len(ad.Images)),
		slog.Int("bodysize", len(ad.Text)),
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
	if ad.Category == "" || ad.Condition == "" || ad.Created == "" || ad.Text == "" {
		return true
	}

	return false
}
