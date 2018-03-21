# gocd
__Change directory to a Go package__

A very simple command line application to automatically change directory based on a Go package name

## Install

```bash
$ go get -v github.com/relvacode/gocd
$ cat `go env GOPATH`/src/github.com/relvacode/gocd/bashrc >> ~/.bashrc
```

  * Run `go get -v github.com/relvacode/gocd` to install package dependencies
  * Add the contents of `bashrc` to your `~/.bashrc`
  * Either `source ~/.bashrc` or re-open your terminal window


## Usage

#### Absolute Package Names

You can navigate to a Go package directly

```bash
$ gocd github.com/username/pkg
```

#### Fuzzy Package Names

You can also use a fuzzy match for the package you want

```bash
$ gocd username/pkg
$ gocd pkg
```

gocd will scan your `GOPATH` and look for matches, if one match is found then you are taken to it. 

If more than one match is found supply the requested index as the second argument.

```bash
$ gocd txt
  0 golang.org/x/text
  1 golang.org/x/text/cases
  2 golang.org/x/text/cmd/gotext
  3 golang.org/x/text/cmd/gotext/examples/extract
  4 golang.org/x/text/cmd/gotext/examples/extract_http
  5 golang.org/x/text/cmd/gotext/examples/extract_http/pkg
  
$ gocd txt 0
```

#### Change Directory to Vendor Parent

```bash
$ gocd ^
```

Using `^` will navigate to the parent package of a vendored directory

##### GOPATH

Go to the `GOPATH` by calling gocd without arguments

```bash
$ gocd
```

