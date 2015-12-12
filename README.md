# gocd
__Change directory to a Go package__

I find myself having to `cd $GOPATH/src/github.com/me/mypkg` a lot whilst moving around a multi-package project.  
So I created gocd as a simple utility to locate and change directory to a Go package based on it's import path, as a bonus I included a search mechanism for quick naviagtion of your source.  

## Installation
  * Ensure your Go environment is properly configured (you have $GOPATH set in your environment)
  * Run `go get -v bitbucket.org/jrelva/gocd` to install package dependencies
  * Run `go install bitbucket.org/jrelva/gocd` to install into your $GOPATH
  * Add the contents of `bashrc` to your `~/.bashrc` with `cat bashrc >> ~/.bashrc`
  * Either `source ~/.bashrc` or re-open your terminal window
                 
## Usage
    gocd github.com/me/mypkg
or

    gocd mypkg
or

    gocd

If the package name is the full import path then use that, otherwise gocd scans `$GOROOT/src` and finds the first occurrence matching 'mypkg' that contains .go files.

Running gocd without arguments will change directory to `$GOROOT/src`.  

Directories with `vendor` or `.git` in the path are ignored.

