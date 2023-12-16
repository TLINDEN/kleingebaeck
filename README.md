## Kleingebäck - kleinanzeigen.de Backup

![Kleingebaeck Logo](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/kleingebaecklogo-small.png)

[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://github.com/tlinden/kleingebaeck/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/kleingebaeck)](https://goreportcard.com/report/github.com/tlinden/kleingebaeck) 
![GitHub License](https://img.shields.io/github/license/tlinden/kleingebaeck)
[![GitHub release](https://img.shields.io/github/v/release/tlinden/kleingebaeck?color=%2300a719)](https://github.com/TLINDEN/kleingebaeck/releases/latest)


This tool can be used to backup ads on the german ad page https://kleinanzeigen.de

It downloads all (or  only the specified ones) ads of  one user into a
directory, each ad into its own subdirectory. The backup will contain
a textfile `Adlisting.txt` which contains the ad contents as the
title, body, price etc. All images will be downloaded as well.

The   tool  doesn't   need   authentication  and   doesn't  have   any
dependencies.  Just  download the  binary for  your platform  from the
releases page and you're good to go.

The releases also include a handy tarball which you can use to install
the tool  system-wide including the  manual page. Just extract  it and
type: `make install`.

## Commandline options:

```
Usage: kleingebaeck [-dvVhmoc] [<ad-listing-url>,...]
Options:
--user,-u <uid>        Backup ads from user with uid <uid>.
--debug, -d            Enable debug output.
--verbose,-v           Enable verbose output.
--output-dir,-o <dir>  Set output dir (default: current directory)
--manual,-m            Show manual.
--config,-c <file>     Use config file <file> (default: ~/.kleingebaeck).

If one  or more <ad-listing-url>'s  are specified, only  backup those,
otherwise backup all ads of the given user.
```

## Configfile

You can create a config file to save typing. By default
`~/.kleingebaeck.hcl` is being used but you can specify one with
`-c` as well.

Format is simple:

```
user = 1010101
verbose = true
outdir = "test"
template = ""
```

## Usage

To setup the tool, you need to lookup your userid on
kleinanzeigen.de. Go to your ad overview page while NOT being logged
in:

https://www.kleinanzeigen.de/s-bestandsliste.html?userId=XXXXXX

The `XXXXX` part is your userid.

Put it into the configfile as outlined above. Also specify an output
directory. Then just execute `kleingebaeck`.


## Kleingebäck?

The name is derived from "kleinanzeigen backup": "klein" (german for
small) and "back". In german "bäck" is spelled the same as the english
"back" so "kleinbäck" was short enough, but it's not a valid german
word. "Kleingebäck" however is: it means "Cookies" in english :)

## Getting help

Although I'm happy  to hear from kleingebaeck users  in private email,
that's the best way for me to forget to do something.

In order to report a bug,  unexpected behavior, feature requests or to
submit    a    patch,    please    open   an    issue    on    github:
https://github.com/TLINDEN/kleingebaeck/issues.

Please repeat the failing command with debugging enabled `-d` and
include the output in the issue.

## Copyright und License

Licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Author

T.v.Dein <tom AT vondein DOT org>

