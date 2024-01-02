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
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

const (
	VERSION         string = "0.1.0"
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

const Usage string = `This is kleingebaeck, the kleinanzeigen.de backup tool.

Usage: kleingebaeck [-dvVhmoclu] [<ad-listing-url>,...]

Options:
--user    -u <uid>      Backup ads from user with uid <uid>.
--debug   -d            Enable debug output.
--verbose -v            Enable verbose output.
--outdir  -o <dir>      Set output dir (default: current directory)
--limit   -l <num>      Limit the ads to download to <num>, default: load all.
--config  -c <file>     Use config file <file> (default: ~/.kleingebaeck).
--manual  -m            Show manual.
--help    -h            Show usage.
--version -V            Show program version.

If one  or more ad listing url's  are specified, only  backup those,
otherwise backup all ads of the given user.`

type Config struct {
	Verbose          bool   `koanf:"verbose"` // loglevel=info
	Debug            bool   `koanf:"debug"`   // loglevel=debug
	Showversion      bool   `koanf:"version"` // -v
	Showhelp         bool   `koanf:"help"`    // -h
	Showmanual       bool   `koanf:"manual"`  // -m
	User             int    `koanf:"user"`
	Outdir           string `koanf:"outdir"`
	Template         string `koanf:"template"`
	Loglevel         string `koanf:"loglevel"`
	Limit            int    `koanf:"limit"`
	Adlinks          []string
	StatsCountAds    int
	StatsCountImages int
}

func (c *Config) IncrAds() {
	c.StatsCountAds++
}

func (c *Config) IncrImgs(num int) {
	c.StatsCountImages += num
}

// load commandline flags and config file
func InitConfig(w io.Writer) (*Config, error) {
	var k = koanf.New(".")

	// determine template based on os
	template := DefaultTemplate
	if runtime.GOOS == "windows" {
		template = DefaultTemplateWin
	}

	// Load default values using the confmap provider.
	if err := k.Load(confmap.Provider(map[string]interface{}{
		"template": template,
		"outdir":   ".",
		"loglevel": "notice",
		"userid":   0,
	}, "."), nil); err != nil {
		return nil, err
	}

	// setup custom usage
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Fprintln(w, Usage)
		os.Exit(0)
	}

	// parse commandline flags
	f.StringP("config", "c", "", "config file")
	f.StringP("outdir", "o", "", "directory where to store ads")
	f.IntP("user", "u", 0, "user id")
	f.IntP("limit", "l", 0, "limit ads to be downloaded (default 0, unlimited)")
	f.BoolP("verbose", "v", false, "be verbose")
	f.BoolP("debug", "d", false, "enable debug log")
	f.BoolP("version", "V", false, "show program version")
	f.BoolP("help", "h", false, "show usage")
	f.BoolP("manual", "m", false, "show manual")

	if err := f.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	// generate a  list of config files to try  to load, including the
	// one provided via -c, if any
	var configfiles []string
	configfile, _ := f.GetString("config")
	home, _ := os.UserHomeDir()
	if configfile != "" {
		configfiles = []string{configfile}
	} else {
		configfiles = []string{
			"/etc/kleingebaeck.conf", "/usr/local/etc/kleingebaeck.conf", // unix variants
			filepath.Join(home, ".config", "kleingebaeck", "config"),
			filepath.Join(home, ".kleingebaeck"),
			"kleingebaeck.conf",
		}
	}

	// Load the config file[s]
	for _, cfgfile := range configfiles {
		if path, err := os.Stat(cfgfile); !os.IsNotExist(err) {
			if !path.IsDir() {
				if err := k.Load(file.Provider(cfgfile), toml.Parser()); err != nil {
					return nil, errors.New("error loading config file: " + err.Error())
				}
			}
		}
		// else: we ignore the file if it doesn't exists
	}

	// command line overrides config file
	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		return nil, errors.New("error loading flags: " + err.Error())
	}

	// fetch values
	conf := &Config{}
	if err := k.Unmarshal("", &conf); err != nil {
		return nil, errors.New("error unmarshalling: " + err.Error())
	}

	// adjust loglevel
	switch conf.Loglevel {
	case "verbose":
		conf.Verbose = true
	case "debug":
		conf.Debug = true
	}

	// are there any args left on commandline? if so threat them as adlinks
	conf.Adlinks = f.Args()

	return conf, nil
}
