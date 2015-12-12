# gocd
__Change directory to a Go package__

## Installation
  * Ensure your Go environment is properly configured (you have $GOPATH set in your environment)
  * Run `go install`
  * Paste this into your `~/.bashrc` file  
  
        gocd () {
          if ! dir=$($GOPATH/bin/gocd $1 2>&1); then
            echo "$dir"
          else
            cd "$dir"
          fi
        }
  * For package tab completion add this to your `~/.bashrc` too

        _gopath () {
          local cur
          COMPREPLY=()
          cur=${COMP_WORDS[COMP_CWORD]}
          k=0
          for j in $( compgen -f "$GOPATH/src/$cur" ); do
            if [ -d "$j" ]; then
              COMPREPLY[k++]=${j#$GOPATH/src/}
            fi
          done
        }  
        complete -o nospace -F _gopath gocd          
  * Either `source ~/.bashrc` or re-open your terminal session
  
                 
## Usage
    gocd github.com/me/mypkg
or

    gocd mypkg
or

    gocd

gocd scans `$GOROOT/src` and finds the first occurrence matching 'mypkg' that contains .go files.

Running gocd without arguments will change directory to `$GOROOT/src`

