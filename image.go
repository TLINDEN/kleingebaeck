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
	"image/jpeg"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
)

const MaxDistance = 3

type Image struct {
	Filename string
	Hash     *goimagehash.ImageHash
	Data     *bytes.Buffer
	Uri      string
}

// used for logging to avoid printing Data
func (img *Image) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("filename", img.Filename),
		slog.String("uri", img.Uri),
		slog.String("hash", img.Hash.ToString()),
	)
}

// holds all images of an ad
type Cache []*goimagehash.ImageHash

func NewImage(buf *bytes.Buffer, filename string, uri string) *Image {
	img := &Image{
		Filename: filename,
		Uri:      uri,
		Data:     buf,
	}

	return img
}

// Calculate diff hash of the image
func (img *Image) CalcHash() error {
	jpgdata, err := jpeg.Decode(img.Data)
	if err != nil {
		return err
	}

	hash1, err := goimagehash.DifferenceHash(jpgdata)
	if err != nil {
		return err
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
	} else {
		return false
	}
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
		return nil, err
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

			img := NewImage(data, filename, "")
			if err = img.CalcHash(); err != nil {
				return nil, err
			}

			slog.Debug("Caching image from file system", "image", img, "hash", img.Hash.ToString())
			cache = append(cache, img.Hash)
		}
	}

	//return nil, errors.New("ende")
	return cache, nil
}
