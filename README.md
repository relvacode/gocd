# gocd
__Change directory to a Go package__

In a multi-package project it's sometimes tedious to have to switch the current directory using `cd $GOPATH/src/github.com/username/pkg`.

gocd is a very simple command line application to automatically change directory based on a go package name.

## Usage

```bash
$ gocd github.com/username/pkg
$ gocd username/pkg
$ gocd pkg
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
$ gocd golangtext
golang.org/x/text/currency
golang.org/x/text/unicode
golang.org/x/text/message
golang.org/x/text/collate
golang.org/x/text/secure
golang.org/x/text/search
golang.org/x/text/runes
golang.org/x/text/width
golang.org/x/text/cases
golang.org/x/text
```
