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
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"

	"github.com/corona10/goimagehash"
)

const MaxDistance = 3

type Image struct {
	Filename string
	Hash     *goimagehash.ImageHash
	Data     *bytes.Reader
	URI      string
	Mime     string
}

// used for logging to avoid printing Data
func (img *Image) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("filename", img.Filename),
		slog.String("uri", img.URI),
		slog.String("hash", img.Hash.ToString()),
	)
}

// holds all images of an ad
type Cache []*goimagehash.ImageHash

// filename comes from the scraper, it contains directory/base w/o suffix
func NewImage(buf *bytes.Reader, filename, uri string) (*Image, error) {
	_, imgconfig, err := image.DecodeConfig(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	_, err = buf.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek(0) on image buffer: %w", err)
	}

	if imgconfig == "jpeg" {
		// we're  using the  format as file  extension, but  have used
		// "jpg" in the past, so to be backwards compatible, stay with
		// it.
		imgconfig = "jpg"
	}

	if imgconfig == "" {
		return nil, fmt.Errorf("failed to process image: unknown or unsupported image format (supported: jpg,png,gif,webp)")
	}

	filename += "." + imgconfig

	img := &Image{
		Filename: filename,
		URI:      uri,
		Data:     buf,
		Mime:     imgconfig,
	}

	slog.Debug("image MIME", "mime", img.Mime)

	return img, nil
}

// Calculate diff hash of the image
func (img *Image) CalcHash() error {
	jpgdata, format, err := image.Decode(img.Data)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	if format == "" {
		return fmt.Errorf("failed to decode image: unknown or unsupported image format (supported: jpg,png,gif,webp)")
	}

	hash1, err := goimagehash.DifferenceHash(jpgdata)
	if err != nil {
		return fmt.Errorf("failed to calculate diff hash of image: %w", err)
	}

	img.Hash = hash1

	return nil
}

// checks if 2 images are similar enough to be considered the same
func (img *Image) Similar(hash *goimagehash.ImageHash) bool {
	distance, err := img.Hash.Distance(hash)
	if err != nil {
		slog.Debug("failed to compute diff hash distance", "error", err)

		return false
	}

	if distance < MaxDistance {
		slog.Debug("distance computation", "image-A", img.Hash.ToString(),
			"image-B", hash.ToString(), "distance", distance)

		return true
	}

	return false
}

// check current image against all known hashes.
func (img *Image) SimilarExists(cache Cache) bool {
	for _, otherimg := range cache {
		if img.Similar(otherimg) {
			return true
		}
	}

	return false
}

// read all  JPG images  in a  ad directory,  compute diff  hashes and
// store the results in the slice Images
func ReadImages(addir string, dont bool) (Cache, error) {
	files, err := os.ReadDir(addir)
	if err != nil {
		return nil, fmt.Errorf("failed to read ad directory contents: %w", err)
	}

	cache := Cache{}

	if dont {
		// forced download, -f given
		return cache, nil
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !file.IsDir() && (ext == ".jpg" || ext == ".jpeg" || ext == ".JPG" || ext == ".JPEG") {
			filename := filepath.Join(addir, file.Name())

			data, err := ReadImage(filename)
			if err != nil {
				return nil, err
			}

			reader := bytes.NewReader(data.Bytes())

			img, err := NewImage(reader, filename, "")
			if err != nil {
				return nil, err
			}

			if err := img.CalcHash(); err != nil {
				return nil, err
			}

			if img.Hash != nil {
				slog.Debug("Caching image from file system", "image", img, "hash", img.Hash.ToString())
			}

			cache = append(cache, img.Hash)
		}
	}

	return cache, nil
}
