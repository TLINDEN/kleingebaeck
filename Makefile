# Copyright © 2023 Thomas von Dein

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.


#
# no need to modify anything below
tool      = kleingebaeck
VERSION   = $(shell grep VERSION config.go | head -1 | cut -d '"' -f2)
archs     = darwin freebsd linux windows
PREFIX    = /usr/local
UID       = root
GID       = 0
HAVE_POD := $(shell pod2text -h 2>/dev/null)

all: $(tool).1 $(tool).go buildlocal

%.1: %.pod
ifdef HAVE_POD
	  pod2man -c "User Commands" -r 1 -s 1 $*.pod > $*.1
endif

%.go: %.pod
ifdef HAVE_POD
	  echo "package main" > $*.go
	  echo >> $*.go
	  echo "var manpage = \`" >> $*.go
	  pod2text $*.pod >> $*.go
	  echo "\`" >> $*.go
endif

buildlocal:
	CGO_LDFLAGS='-static' go build -tags osusergo,netgo -ldflags "-extldflags=-static" -o $(tool)

install: buildlocal
	install -d -o $(UID) -g $(GID) $(PREFIX)/bin
	install -d -o $(UID) -g $(GID) $(PREFIX)/man/man1
	install -o $(UID) -g $(GID) -m 555 $(tool)  $(PREFIX)/sbin/
	install -o $(UID) -g $(GID) -m 444 $(tool).1 $(PREFIX)/man/man1/

clean:
	rm -rf $(tool) coverage.out testdata t/out

test: clean
	mkdir -p t/out
	go test ./... $(ARGS)

testlint: test lint

lint:
	golangci-lint run

lint-full:
	golangci-lint run --enable-all --exclude-use-default --disable exhaustivestruct,exhaustruct,depguard,interfacer,deadcode,golint,structcheck,scopelint,varcheck,ifshort,maligned,nosnakecase,godot,funlen,gofumpt,cyclop

testfuzzy: clean
	go test -fuzz ./... $(ARGS)

singletest:
	@echo "Call like this: make singletest TEST=TestPrepareColumns ARGS=-v"
	go test -run $(TEST) $(ARGS)

cover-report:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

goupdate:
	go get -t -u=patch ./...

buildall:
	./mkrel.sh $(tool) $(VERSION)

release: buildall
	gh release create v$(VERSION) --generate-notes releases/*

show-versions: buildlocal
	@echo "### kleingebaeck version:"
	@./kleingebaeck -V

	@echo
	@echo "### go module versions:"
	@go list -m all

	@echo
	@echo "### go version used for building:"
	@grep -m 1 go go.mod

# lint:
# 	golangci-lint run -p bugs -p unused
