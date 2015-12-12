# gocd
__Change directory to a Go package__

## Installation
  * Ensure your Go environment is properly configured (you have $GOPATH set in your environment)
  * Run `go get -v` to install package dependencies
  * Run `go install` to install into your $GOPATH
  * Add the contents of `bashrc` to your `~/.bashrc` with `cat bashrc >> ~/.bashrc`
  * Either `source ~/.bashrc` or re-open your terminal window
                 
## Usage
    gocd github.com/me/mypkg
or

    gocd mypkg
or

    gocd

gocd scans `$GOROOT/src` and finds the first occurrence matching 'mypkg' that contains .go files.
gocd uses fuzzy matching for packages that cannot be found to suggest possible matching packages:


Running gocd without arguments will change directory to `$GOROOT/src`

