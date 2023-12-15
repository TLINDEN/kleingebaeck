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
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
)

type Config struct {
	Verbose bool   `hcl:"verbose"`
	User    int    `hcl:"user"`
	Outdir  string `hcl:"outdir"`
}

func ParseConfigfile(file string) (*Config, error) {
	c := Config{}
	if path, err := os.Stat(file); !os.IsNotExist(err) {
		if !path.IsDir() {
			configstring, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}

			err = hclsimple.Decode(
				path.Name(), configstring,
				nil, &c,
			)

			if err != nil {
				return nil, err
			}
		}
	}

	return &c, nil
}
