package main

import (
	"errors"
)

// ErrNoMatch means no match at all was found.
var ErrNoMatch = errors.New("No matching package found")
