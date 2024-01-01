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
	"os"
	"strings"
	"testing"
)

type Cmdline struct {
	name     string
	args     string
	expect   string
	exitcode int
}

func TestMain(t *testing.T) {
	base := "kleingebaeck -c t/config-empty.conf"
	cmdlines := []Cmdline{
		{
			name:     "version",
			args:     base + " -V",
			expect:   "This is",
			exitcode: 0,
		},
		{
			name:     "help",
			args:     base + " -h",
			expect:   "Usage:",
			exitcode: 0,
		},
		{
			name:     "debug",
			args:     base + " -d",
			expect:   "program_info",
			exitcode: 1,
		},
		{
			name:     "no-args-no-user",
			args:     base,
			expect:   "invalid or no user id",
			exitcode: 1,
		},
		//  FIXME: add  scrape  tests as  well, re-use  scrape_test.go
		// stuff or simply move all of it to here
	}

	oldargs := os.Args
	defer func() { os.Args = oldargs }()

	for _, c := range cmdlines {
		var buf bytes.Buffer
		os.Args = strings.Split(c.args, " ")

		ret := Main(&buf)

		if ret != c.exitcode {
			t.Errorf("%s with cmd <%s> did not exit with %d but %d",
				c.name, c.args, c.exitcode, ret)
		}

		if !strings.Contains(buf.String(), c.expect) {
			t.Errorf("%s with cmd <%s> output did not match.\nExpect: %s\n   Got: %s\n",
				c.name, c.args, c.expect, buf.String())
		}
	}

}
