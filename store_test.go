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
	"testing"
)

// this is  a weird thing.  WriteImage() is being called  in scrape.go
// which is  being tested by  TestMain() in main_test.go.  However, it
// doesn't  show up  in the  coverage report  for unknown  reasons, so
// here's a single test for it
func TestWriteImage(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	file := "t/out/t.jpg"

	err := WriteImage(file, buf)
	if err != nil {
		t.Errorf("Could not write mock image to %s: %s", file, err.Error())
	}

}
