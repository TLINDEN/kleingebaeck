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
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/lmittmann/tint"
	flag "github.com/spf13/pflag"
)

const Usage string = `This is kleingebaeck, the kleinanzeigen.de backup tool.
Usage: kleingebaeck [-dvVhmoc] [<ad-listing-url>,...]
Options:
--user,-u <uid>        Backup ads from user with uid <uid>.
--debug, -d            Enable debug output.
--verbose,-v           Enable verbose output.
--output-dir,-o <dir>  Set output dir (default: current directory)
--manual,-m            Show manual.
--config,-c <file>     Use config file <file> (default: ~/.kleingebaeck).

If one  or more <ad-listing-url>'s  are specified, only  backup those,
otherwise backup all ads of the given user.`

const LevelNotice = slog.Level(2)

func main() {
	os.Exit(Main())
}

func Main() int {
	logLevel := &slog.LevelVar{}
	opts := &tint.Options{
		Level:     logLevel,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	logLevel.Set(LevelNotice)
	var handler slog.Handler = tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	showversion := false
	showhelp := false
	showmanual := false
	enabledebug := false
	enableverbose := false
	uid := 0
	configfile := os.Getenv("HOME") + "/.kleingebaeck.hcl"
	dir := ""

	flag.BoolVarP(&enabledebug, "debug", "d", false, "debug mode")
	flag.BoolVarP(&enableverbose, "verbose", "v", false, "be verbose")
	flag.BoolVarP(&showversion, "version", "V", false, "show version")
	flag.BoolVarP(&showhelp, "help", "h", false, "show usage")
	flag.BoolVarP(&showmanual, "manual", "m", false, "show manual")
	flag.IntVarP(&uid, "user", "u", uid, "user id")
	flag.StringVarP(&dir, "output-dir", "o", dir, "where to store ads")
	flag.StringVarP(&configfile, "config", "c", configfile, "config file")

	flag.Parse()

	if showversion {
		fmt.Printf("This is kleingebaeck version %s\n", VERSION)
		return 0
	}

	if showhelp {
		fmt.Println(Usage)
		return 0
	}

	if showmanual {
		err := man()
		if err != nil {
			return Die(err)
		}
		return 0
	}

	conf, err := ParseConfigfile(configfile)
	if err != nil {
		return Die(err)
	}

	if enableverbose || *conf.Verbose {
		logLevel.Set(slog.LevelInfo)
	}

	if enabledebug {
		// we're using a more verbose logger in debug mode
		buildInfo, _ := debug.ReadBuildInfo()
		opts := &tint.Options{
			Level:     logLevel,
			AddSource: true,
		}

		logLevel.Set(slog.LevelDebug)
		var handler slog.Handler = tint.NewHandler(os.Stdout, opts)
		debuglogger := slog.New(handler).With(
			slog.Group("program_info",
				slog.Int("pid", os.Getpid()),
				slog.String("go_version", buildInfo.GoVersion),
			),
		)
		slog.SetDefault(debuglogger)
	}

	slog.Debug("config", "conf", conf)

	if len(dir) == 0 {
		if len(*conf.Outdir) > 0 {
			dir = *conf.Outdir
		} else {
			dir = Defaultdir
		}
	}

	// prepare output dir
	err = Mkdir(dir)
	if err != nil {
		return Die(err)
	}

	// which template to use
	template := DefaultTemplate
	if runtime.GOOS == "windows" {
		template = DefaultTemplateWin
	}
	if len(*conf.Template) > 0 {
		template = *conf.Template
	}

	// directly backup ad listing[s]
	if len(flag.Args()) >= 1 {
		for _, uri := range flag.Args() {
			err := Scrape(uri, dir, template)
			if err != nil {
				return Die(err)
			}
		}

		return 0
	}

	// backup all ads of the given user (via config or cmdline)
	if uid == 0 && *conf.User > 0 {
		uid = *conf.User
	}

	if uid > 0 {
		err := Start(fmt.Sprintf("%d", uid), dir, template)
		if err != nil {
			return Die(err)
		}
	} else {
		return Die(errors.New("invalid or no user id specified"))
	}

	return 0
}

func Die(err error) int {
	slog.Error("Failure", "error", err.Error())
	return 1
}
