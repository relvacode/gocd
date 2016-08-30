# gocd
__Change directory to a Go package__

In a multi-package project it's sometimes tedious to have to switch the current directory using `cd $GOPATH/src/github.com/me/pkg`.

gocd is a very simple command line application to automatically change directory based on a go package name.

## Usage
    gocd github.com/me/mypkg
    
or

    gocd mypkg
    
or

    gocd

If the package name is the full import path then use that, otherwise gocd scans `$GOPATH/src` and finds the first occurrence matching `mypkg` that contains .go files.

Running gocd without arguments will change directory to `$GOPATH/src`.  

Directories with `vendor` or `.git` in the path are ignored.

## Installation

  * Ensure your Go environment is properly configured (you have `$GOPATH` set in your environment)
  * Run `go get -v github.com/relvacode/gocd` to install package dependencies
  * Add the contents of `bashrc` to your `~/.bashrc` with `cat $GOPATH/src/github.com/relvacode/gocd/bashrc >> ~/.bashrc`
  * Either `source ~/.bashrc` or re-open your terminal window
