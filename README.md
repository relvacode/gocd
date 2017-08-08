# gocd
__Change directory to a Go package__

In a multi-package project it's sometimes tedious to have to switch the current directory using `cd $GOPATH/src/github.com/username/pkg`.

gocd is a very simple command line application to automatically change directory based on a go package name.

## Usage

```bash
$ gocd github.com/username/pkg
$ gocd username/pkg
$ gocd pkg
$ gocd pkg <@>
$ gocd
```

If the package name is the full import path then use that, otherwise gocd scans `$GOPATH/src` and finds the first matching occurence of a directory containing .go files. If no arguments are supplied gocd will change directory to `$GOPATH/src`. Directories with `vendor` or `.git` in the path are ignored.

## Installation

```bash
$ go get -v github.com/relvacode/gocd
$ cat `go env GOPATH`/src/github.com/relvacode/gocd/bashrc >> ~/.bashrc
```

  * Run `go get -v github.com/relvacode/gocd` to install package dependencies
  * Add the contents of `bashrc` to your `~/.bashrc`
  * Either `source ~/.bashrc` or re-open your terminal window

## Suggestions

If no direct match can be found, `gocd` will look for the top 10 nearest packages using Levenshtein distance.

```bash
$ gocd txt
  0 golang.org/x/text/internal/format/plural
  1 golang.org/x/text/message
  2 golang.org/x/text/encoding/korean
  3 golang.org/x/text/encoding/japanese
  4 golang.org/x/text/collate/tools/colcmp
  5 golang.org/x/text/encoding/ianaindex
  6 github.com/kr/text/mc
  7 golang.org/x/text/internal/utf8internal
  8 golang.org/x/text/currency
  9 golang.org/x/text/cases
```

Go to a specific package at the correct index you wanted by using the index as the second argument

```bash
$ gocd txt 1 # cd to golang.org/x/text/message
```
