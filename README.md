## Kleingebäck - kleinanzeigen.de Backup

![Kleingebaeck Logo](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/kleingebaecklogo-small.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/kleingebaeck)](https://goreportcard.com/report/github.com/tlinden/kleingebaeck) 
[![Actions](https://github.com/tlinden/kleingebaeck/actions/workflows/ci.yaml/badge.svg)](https://github.com/tlinden/kleingebaeck/actions)
[![Go Coverage](https://github.com/tlinden/kleingebaeck/wiki/coverage.svg)](https://raw.githack.com/wiki/tlinden/kleingebaeck/coverage.html)
![GitHub License](https://img.shields.io/github/license/tlinden/kleingebaeck)
[![GitHub release](https://img.shields.io/github/v/release/tlinden/kleingebaeck?color=%2300a719)](https://github.com/TLINDEN/kleingebaeck/releases/latest)


This tool can be used to backup ads on the german ad page https://kleinanzeigen.de

It downloads all (or  only the specified ones) ads of  one user into a
directory, each ad into its own subdirectory. The backup will contain
a textfile `Adlisting.txt` which contains the ad contents as the
title, body, price etc. All images will be downloaded as well.

## Screenshots

This is the index of my kleinanzeigen.de Account:

![Index](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/kleinanzeigen-index.png)

Here I download my ads on the commandline:

![Download](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/kleinanzeigen-download.png)

And this is the backup directory after download:

![Download](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/kleinanzeigen-backup.png)

Here's a directory for one ad:

![Download](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/kleinanzeigen-ad.png)

**The same thing under windows:**

Downloading ads:

![Download](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/cmd-windows.jpg)

Backup directory after download:

![Download](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/liste-windows.jpg)

And one ad listing directory:

![Download](https://github.com/TLINDEN/kleingebaeck/blob/main/.github/assets/adlisting-windows.jpg)

## Installation

The   tool  doesn't   need   authentication  and   doesn't  have   any
dependencies.  Just  download the  binary for  your platform  from the
releases page and you're good to go.

### Installation using a pre-compiled binary

Go            to             the            [latest            release
page](https://github.com/TLINDEN/kleingebaeck/releases/latest)     and
look for your OS and platform. There are two options to install the binary:

1.    Directly   download    the    binary    for   your    platoform,
   e.g. `kleingebaeck-linux-amd64-0.0.5`, rename  it to `kleingebaeck`
   (or  whatever  you  like  more!)  and put  it  into  your  bin  dir
   (e.g. `$HOME/bin` or as root to `/usr/local/bin`).

Be sure to verify the signature of the binary file. For this also download the matching `kleingebaeck-linux-amd64-0.0.5.sha256` file and:

```shell
cat kleingebaeck-linux-amd64-0.0.5.sha25 && sha256sum kleingebaeck-linux-amd64-0.0.5
```
You should see the same SHA256 hash.

2.  You  may  also  download  a  binary  tarball  for  your  platform,
   e.g.  `kleingebaeck-linux-amd64-0.0.5.tar.gz`,  unpack and  install
   it. GNU Make is required for this:
   
```shell
tar xvfz kleingebaeck-linux-amd64-0.0.5.tar.gz
cd kleingebaeck-linux-amd64-0.0.5
sudo make install
```

### Installation from source

You will need the Golang toolchain  in order to build from source. GNU
Make will also help but is not strictly neccessary.

If you want to compile the tool yourself, use `git clone` to clone the
repository.   Then   execute   `go    mod   tidy`   to   install   all
dependencies. Then  just enter `go  build` or -  if you have  GNU Make
installed - `make`.

To install after building either copy the binary or execute `sudo make install`.

## Commandline options:

```
Usage: kleingebaeck [-dvVhmoc] [<ad-listing-url>,...]
Options:
--user    -u <uid>      Backup ads from user with uid <uid>.
--debug   -d            Enable debug output.
--verbose -v            Enable verbose output.
--outdir  -o <dir>      Set output dir (default: current directory)
--limit   -l <num>      Limit the ads to download to <num>, default: load all.
--config  -c <file>     Use config file <file> (default: ~/.kleingebaeck).
--manual  -m            Show manual.
--help    -h            Show usage.

If one  or more <ad-listing-url>'s  are specified, only  backup those,
otherwise backup all ads of the given user.
```

## Configfile

You can create a config file to save typing. By default
`~/.kleingebaeck` is being used but you can specify one with
`-c` as well.

Format is simple:

```
user = 1010101
loglevel = verbose
outdir = "test"
```

## Usage

To setup the tool, you need to lookup your userid on
kleinanzeigen.de. Go to your ad overview page while NOT being logged
in:

https://www.kleinanzeigen.de/s-bestandsliste.html?userId=XXXXXX

The `XXXXX` part is your userid.

Put it into the configfile as outlined above. Also specify an output
directory. Then just execute `kleingebaeck`.

Inside the output directory you'll find a new subdirectory for each
ad. Every directory contains a file `Adlisting.txt`, which will look
somewhat like this:

```default
Title: A book I sell
Price: 99 € VB
Id: 1919191919
Category: Sachbücher
Condition: Sehr Gut
Created: 10.12.2023

This is the description text.

Pay with paypal.
```

You can change the formatting using the `template` config
variable. The supplied sample config contains the default template.

All images will be stored in the same directory.

## Documentation

You can read the documentation [online](https://github.com/TLINDEN/kleingebaeck/blob/main/kleingebaeck.pod) or locally once you have installed kleingebaeck with: `kleingebaeck --manual`.

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

## Related projects

I could not find any projects specifically designed to backup
kleinanzeigen.de ads, however there's a bot project which is also able
to download ads:
[kleinanzeigen-bot](https://github.com/Second-Hand-Friends/kleinanzeigen-bot/). However,
be aware that kleinanzeigen.de is actively fighting bots! Look at this
[issue](https://github.com/Second-Hand-Friends/kleinanzeigen-bot/issues/219). The
problem with these kind of bots is, that they login into your account
using your credentials. If the company is able to detect bot activity
they can associate it easily with your account and **lock you
out**. So be careful.

**kleingebäck** doesn't need to login, it just accesses public
available web pages. Kleinanzeigen.de could hardly do anything against
it, once because it is legal. There's no difference between a browser
and a commandline client. Both run on the clientside and it is not
kleinanzeigen.de's decision which software one uses to access their
pages. And second: because you can use it to download any ads, not
just yours. So it is not really clear if the activity is associated in
any way with the ad owner. In addition to that comes the fact that
kleingebäck is just a backup tool. It is not intendet to be used on a
daily basis. You cannot use it to view regular ads or maintain your
own ads. You'll need to use the mobile app or the browser page with a
login. So, in my point of view, the risk is very minimal.

There is another Tool available named [kleinanzeigen-enhanded](https://kleinanzeigen-enhanced.de/). It is a complete Ad management system targeting primarily commercial users. You have to pay a monthly fee, perhaps there's also a free version available, but I haven't checked. The tool is implemented as a Chrome browser extension, which explains why it was possible to implement it without an API. It seems to be a nice solution for power users by the looks of it. And it includes backups.

## Copyright und License

Licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Author

T.v.Dein <tom AT vondein DOT org>

