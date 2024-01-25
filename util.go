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
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/mattn/go-isatty"
)

func Mkdir(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
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
		return fmt.Errorf("failed to execute 'less': %w", err)
	}

	return nil
}

// returns TRUE if stdout is NOT a tty or windows
func IsNoTty() bool {
	if runtime.GOOS == "windows" || !isatty.IsTerminal(os.Stdout.Fd()) {
		return true
	}

	// it is a tty
	return false
}

func GetThrottleTime() time.Duration {
	return time.Duration(rand.Intn(MaxThrottle-MinThrottle+1)+MinThrottle) * time.Millisecond
}
