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
	"io"
	"log/slog"
	"os"
	"strings"
	tpl "text/template"
)

func WriteAd(dir string, ad *Ad, template string) error {
	// prepare output dir
	dir = dir + "/" + ad.Slug
	err := Mkdir(dir)
	if err != nil {
		return err
	}

	// write ad file
	listingfile := strings.Join([]string{dir, "Adlisting.txt"}, "/")
	f, err := os.Create(listingfile)
	if err != nil {
		return err
	}

	ad.Text = strings.ReplaceAll(ad.Text, "<br/>", "\n")

	tmpl, err := tpl.New("adlisting").Parse(template)
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, ad)
	if err != nil {
		return err
	}

	slog.Info("wrote ad listing", "listingfile", listingfile)

	return nil
}

func WriteImage(filename string, reader io.ReadCloser) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}
