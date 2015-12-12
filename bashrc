# Bash wrapper to change directory to the output of gocd
gocd () {
  if dir=$($GOPATH/bin/gocd $1); then
    cd "$dir"
  fi
}

# Optional tab completion wrapper for $GOPATH/src
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