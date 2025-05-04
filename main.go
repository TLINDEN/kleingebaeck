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
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/inconshreveable/mousetrap"
	"github.com/lmittmann/tint"
	"github.com/tlinden/yadu"
)

const LevelNotice = slog.Level(2)

func main() {
	os.Exit(Main(os.Stdout))
}

func init() {
	// if we're running on Windows  AND if the user double clicked the
	// exe  file from explorer, we  tell them and then  wait until any
	// key has been hit, which  will make the cmd window disappear and
	// thus give the user time to read it.
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("Do no double click kleingebaeck.exe!")
			fmt.Println("Please open a command shell and run it from there.")
			fmt.Println()
			fmt.Print("Press any key to quit: ")
			_, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				panic(err)
			}
		}
	}
}

func Main(output io.Writer) int {
	logLevel := &slog.LevelVar{}
	opts := &tint.Options{
		Level:     logLevel,
		AddSource: false,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			// Remove time from the output
			if attr.Key == slog.TimeKey {
				return slog.Attr{}
			}

			return attr
		},
		NoColor: IsNoTty(),
	}

	logLevel.Set(LevelNotice)

	handler := tint.NewHandler(output, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	conf, err := InitConfig(output)
	if err != nil {
		return Die(err)
	}

	if conf.Showversion {
		_, err := fmt.Fprintf(output, "This is kleingebaeck version %s\n", VERSION)
		if err != nil {
			panic(err)
		}

		return 0
	}

	if conf.Showhelp {
		_, err := fmt.Fprintln(output, Usage)
		if err != nil {
			panic(err)
		}

		return 0
	}

	if conf.Showmanual {
		err := man()
		if err != nil {
			return Die(err)
		}

		return 0
	}

	if conf.Verbose {
		logLevel.Set(slog.LevelInfo)
	}

	if conf.Debug {
		// we're using a more verbose logger in debug mode
		buildInfo, _ := debug.ReadBuildInfo()
		opts := &yadu.Options{
			Level:     logLevel,
			AddSource: true,
			//NoColor:   IsNoTty(),
		}

		logLevel.Set(slog.LevelDebug)

		handler := yadu.NewHandler(output, opts)
		debuglogger := slog.New(handler).With(
			slog.Group("program_info",
				slog.Int("pid", os.Getpid()),
				slog.String("go_version", buildInfo.GoVersion),
			),
		)
		slog.SetDefault(debuglogger)
	}

	slog.Debug("config", "conf", conf)

	// prepare output dir
	outdir, err := OutDirName(conf)
	if err != nil {
		return Die(err)
	}
	conf.Outdir = outdir

	// used for all HTTP requests
	fetch, err := NewFetcher(conf)
	if err != nil {
		return Die(err)
	}

	// setup ad dir registry, needed to check for duplicates
	DirsVisited = make(map[string]int)

	switch {
	case len(conf.Adlinks) >= 1:
		// directly backup ad listing[s]
		for _, uri := range conf.Adlinks {
			err := ScrapeAd(fetch, uri)
			if err != nil {
				return Die(err)
			}
		}
	case conf.User > 0:
		// backup all ads of the given user (via config or cmdline)
		err := ScrapeUser(fetch)
		if err != nil {
			return Die(err)
		}
	default:
		return Die(errors.New("invalid or no user id or no ad link specified"))
	}

	if conf.StatsCountAds > 0 {
		adstr := "ads"

		if conf.StatsCountAds == 1 {
			adstr = "ad"
		}

		_, err := fmt.Fprintf(output, "Successfully downloaded %d %s with %d images to %s.\n",
			conf.StatsCountAds, adstr, conf.StatsCountImages, conf.Outdir)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := fmt.Fprintf(output, "No ads found.")
		if err != nil {
			panic(err)
		}
	}

	return 0
}

func Die(err error) int {
	slog.Error("Failure", "error", err.Error())

	return 1
}
