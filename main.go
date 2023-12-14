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
	"errors"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

const VERSION string = "0.0.1"
const Useragent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
const Baseuri string = "https://www.kleinanzeigen.de"
const Listuri string = "/s-bestandsliste.html"
const Defaultdir string = "."

func main() {
	os.Exit(Main())
}

func Main() int {
	showversion := false
	showhelp := false
	showmanual := false
	enabledebug := false
	configfile := ""
	dir := Defaultdir

	flag.BoolVarP(&enabledebug, "debug", "d", false, "debug mode")
	flag.BoolVarP(&showversion, "version", "v", false, "show version")
	flag.BoolVarP(&showhelp, "help", "h", false, "show usage")
	flag.BoolVarP(&showmanual, "manual", "m", false, "show manual")
	flag.StringVarP(&dir, "output-dir", "o", dir, "where to store ads")
	flag.StringVarP(&configfile, "config", "c",
		os.Getenv("HOME")+"/.kleingebaeck", "config file")

	flag.Parse()

	if showversion {
		fmt.Printf("This is kleingebaeck version %s\n", VERSION)
		return 0
	}

	/*

		 if showhelp {
			fmt.Println(Usage)
			return 0
		}

		if enabledebug {
			calc.ToggleDebug()
		}

		if showmanual {
			man()
			return 0
		}

	*/

	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return Die(err)
		}
	}

	if len(flag.Args()) == 1 {
		Start(flag.Args()[0], dir)
	}

	return 0
}

func Die(err error) int {
	fmt.Println(err)
	return 1
}
