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
	"runtime/debug"

	"github.com/lmittmann/tint"
)

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
		NoColor: IsNoTty(),
	}

	logLevel.Set(LevelNotice)
	var handler slog.Handler = tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	conf, err := InitConfig()
	if err != nil {
		return Die(err)
	}

	if conf.Showversion {
		fmt.Printf("This is kleingebaeck version %s\n", VERSION)
		return 0
	}

	if conf.Showhelp {
		fmt.Println(Usage)
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
		opts := &tint.Options{
			Level:     logLevel,
			AddSource: true,
			NoColor:   IsNoTty(),
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

	// prepare output dir
	err = Mkdir(conf.Outdir)
	if err != nil {
		return Die(err)
	}

	if len(conf.Adlinks) >= 1 {
		// directly backup ad listing[s]
		for _, uri := range conf.Adlinks {
			err := Scrape(conf, uri)
			if err != nil {
				return Die(err)
			}
		}
	} else if conf.User > 0 {
		// backup all ads of the given user (via config or cmdline)
		err := Start(conf)
		if err != nil {
			return Die(err)
		}
	} else {
		return Die(errors.New("invalid or no user id or no ad link specified"))
	}

	if conf.StatsCountAds > 0 {
		adstr := "ads"
		if conf.StatsCountAds == 1 {
			adstr = "ad"
		}
		fmt.Printf("Successfully downloaded %d %s with %d images to %s.",
			conf.StatsCountAds, adstr, conf.StatsCountImages, conf.Outdir)
		fmt.Println()
	} else {
		fmt.Println("No ads found.")
	}

	return 0
}

func Die(err error) int {
	slog.Error("Failure", "error", err.Error())
	return 1
}
