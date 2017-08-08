# Bash wrapper to change directory to the output of gocd
gocd () {
  if dir=$($GOPATH/bin/gocd $@); then
    cd "$dir"
  fi
} 