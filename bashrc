# Bash wrapper to change directory to the output of gocd
gocd () {
  if dir=$($GOPATH/bin/gocd $1); then
    cd "$dir"
  fi
} 