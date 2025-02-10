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
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	tpl "text/template"
	"time"
)

type OutdirData struct {
	Year, Day, Month string
}

func OutDirName(conf *Config) (string, error) {
	tmpl, err := tpl.New("outdir").Parse(conf.Outdir)
	if err != nil {
		return "", fmt.Errorf("failed to parse outdir template: %w", err)
	}

	buf := bytes.Buffer{}

	now := time.Now()
	data := OutdirData{
		Year:  now.Format("2006"),
		Month: now.Format("01"),
		Day:   now.Format("02"),
	}

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute outdir template: %w", err)
	}

	return buf.String(), nil
}

func AdDirName(conf *Config, advertisement *Ad) (string, error) {
	tmpl, err := tpl.New("adname").Parse(conf.Adnametemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse adname template: %w", err)
	}

	buf := bytes.Buffer{}

	err = tmpl.Execute(&buf, advertisement)
	if err != nil {
		return "", fmt.Errorf("failed to execute adname template: %w", err)
	}

	return buf.String(), nil
}

func WriteAd(conf *Config, advertisement *Ad, addir string) error {
	// prepare output dir
	dir := filepath.Join(conf.Outdir, addir)

	err := Mkdir(dir)
	if err != nil {
		return err
	}

	// write ad file
	listingfile := filepath.Join(dir, "Adlisting.txt")

	listingfd, err := os.Create(listingfile)
	if err != nil {
		return fmt.Errorf("failed to create Adlisting.txt: %w", err)
	}
	defer listingfd.Close()

	if runtime.GOOS == WIN {
		advertisement.Text = strings.ReplaceAll(advertisement.Text, "<br/>", "\r\n")
	} else {
		advertisement.Text = strings.ReplaceAll(advertisement.Text, "<br/>", "\n")
	}

	tmpl, err := tpl.New("adlisting").Parse(conf.Template)
	if err != nil {
		return fmt.Errorf("failed to parse adlisting template: %w", err)
	}

	err = tmpl.Execute(listingfd, advertisement)
	if err != nil {
		return fmt.Errorf("failed to execute adlisting template: %w", err)
	}

	slog.Info("wrote ad listing", "listingfile", listingfile)

	return nil
}

func WriteImage(filename string, reader *bytes.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	_, err = reader.WriteTo(file)

	if err != nil {
		return fmt.Errorf("failed to write to image file: %w", err)
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
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	_, err = buf.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write image into buffer: %w", err)
	}

	return &buf, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil {
		// return false on any error
		return false
	}

	return !info.IsDir()
}

// check if  an addir has  already been  processed by current  run and
// decide what to do
func CheckAdVisited(conf *Config, adname string) bool {
	if Exists(DirsVisited, adname) {
		if conf.ForceDownload {
			slog.Warn("an ad with the same name has already been downloaded, overwriting", "addir", adname)
			return true
		}

		// don't overwrite
		slog.Warn("an ad with the same name has already been downloaded, skipping (use -f to overwrite)", "addir", adname)
		return false
	}

	// overwrite
	return true
}
