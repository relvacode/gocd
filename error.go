package main

import (
	"errors"
)

// ErrNoMatch means no match at all was found.
var ErrNoMatch = errors.New("No matching package found")

// GoPackage is a find result from a PkgFinder.Find().
// It implements Error to be used in a filepath walker.
type GoPackage struct {
	Path string
	Name string
}

func (g GoPackage) Error() string {
	return ""
}
