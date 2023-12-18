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
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const (
	VERSION         string = "0.0.6"
	Baseuri         string = "https://www.kleinanzeigen.de"
	Listuri         string = "/s-bestandsliste.html"
	Defaultdir      string = "."
	DefaultTemplate string = "Title: {{.Title}}\nPrice: {{.Price}}\nId: {{.Id}}\n" +
		"Category: {{.Category}}\nCondition: {{.Condition}}\nCreated: {{.Created}}\n\n{{.Text}}\n"
	DefaultTemplateWin string = "Title: {{.Title}}\r\nPrice: {{.Price}}\r\nId: {{.Id}}\r\n" +
		"Category: {{.Category}}\r\nCondition: {{.Condition}}\r\nCreated: {{.Created}}\r\n\r\n{{.Text}}\r\n"
	Useragent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

type Config struct {
	Verbose  *bool   `hcl:"verbose"`
	User     *int    `hcl:"user"`
	Outdir   *string `hcl:"outdir"`
	Template *string `hcl:"template"`
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
