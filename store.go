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
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	tpl "text/template"
)

func AdDirName(c *Config, ad *Ad) (string, error) {
	tmpl, err := tpl.New("adname").Parse(c.Adnametemplate)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, ad)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func WriteAd(c *Config, ad *Ad) (string, error) {
	// prepare ad dir name
	addir, err := AdDirName(c, ad)
	if err != nil {
		return "", err
	}

	// prepare output dir
	dir := filepath.Join(c.Outdir, addir)
	err = Mkdir(dir)
	if err != nil {
		return "", err
	}

	// write ad file
	listingfile := filepath.Join(dir, "Adlisting.txt")
	f, err := os.Create(listingfile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if runtime.GOOS == "windows" {
		ad.Text = strings.ReplaceAll(ad.Text, "<br/>", "\r\n")
	} else {
		ad.Text = strings.ReplaceAll(ad.Text, "<br/>", "\n")
	}

	tmpl, err := tpl.New("adlisting").Parse(c.Template)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(f, ad)
	if err != nil {
		return "", err
	}

	slog.Info("wrote ad listing", "listingfile", listingfile)

	return addir, nil
}

func WriteImage(filename string, buf []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf)

	if err != nil {
		return err
	}

	return nil
}

func ReadImage(filename string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	if !fileExists(filename) {
		return nil, fmt.Errorf("image %s does not exist", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(data)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
