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
	"errors"
	"os"
	"os/exec"
)

func Mkdir(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func man() error {
	man := exec.Command("less", "-")

	var b bytes.Buffer
	b.Write([]byte(manpage))

	man.Stdout = os.Stdout
	man.Stdin = &b
	man.Stderr = os.Stderr

	err := man.Run()

	if err != nil {
		return err
	}

	return nil
}
