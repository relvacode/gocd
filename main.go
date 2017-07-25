package main

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

// GoPackage is a find result from a PkgFinder.Find().
// It implements Error to be used in a filepath walker.
type GoPackage struct {
	Path string
	Name string
}

func (g GoPackage) Error() string {
	return ""
}

// ErrNoMatch means no match at all was found.
var ErrNoMatch = errors.New("no match found")

// ErrNoDirectMatch indicactes that no direct match was found.
type ErrNoDirectMatch struct {
	matches fuzzy.Ranks
}

// Error returns any possible matches as a comma separated string
func (err ErrNoDirectMatch) Error() string {
	sort.Sort(sort.Reverse(err.matches))

	// Take only the top 10 distances
	if len(err.matches) > 10 {
		err.matches = err.matches[len(err.matches)-10:]
	}

	var names = make([]string, len(err.matches))
	var i int
	for _, n := range err.matches {
		if n.Distance > 50 {
			continue
		}
		names[i] = n.Target
		i++
	}
	names = names[:i]

	return strings.Join(names, "\n")
}

// PkgFinder finds a Go package.
type PkgFinder struct {
	find   string // package to find
	gopath string // GOPATH to use
	seen   map[string]struct{}
}

// Implements filepath.Walker.
// Uses GoPackage{} as en error when a package is found.
// If the dirname of the given path fuzzy matches the find key then add it to the slice of PossibleMatches.
func (w *PkgFinder) walker(path string, i os.FileInfo, err error) (e error) {
	// Skip GOPATH/src
	if path == w.gopath {
		return nil
	}
	// Skip if path contains .git or vendor
	if i.IsDir() && (strings.HasPrefix(i.Name(), ".") || strings.Contains(path, "vendor")) {
		return filepath.SkipDir
	}
	// Ignore if path is a directory or is not a go file.
	if i.IsDir() || !strings.HasSuffix(i.Name(), "go") {
		return nil
	}

	// Scan every component of the relative path until we find a direct match.
	pkg, _ := filepath.Rel(w.gopath, filepath.Dir(path))

	// Skip already seen packages
	_, ok := w.seen[pkg]
	if ok {
		return nil
	}
	w.seen[pkg] = struct{}{}

	components := strings.Split(pkg, string(filepath.Separator))
	for x := len(components) - 1; x >= 0; x-- {
		if w.find == filepath.Join(components[x:]...) {
			return GoPackage{
				Path: filepath.Dir(path),
				Name: pkg,
			}
		}
	}
	return nil
}

// Find a package by the given key
func (w *PkgFinder) Find() (string, error) {
	// If absolute path then go straight to it
	if path.IsAbs(w.find) {
		return w.find, nil
	}

	// If path is a path relative to gopath then use it
	abs := filepath.Join(w.gopath, w.find)
	_, err := os.Stat(abs)
	if err == nil {
		return abs, nil
	}

	err = filepath.Walk(w.gopath, w.walker)
	if pkg, ok := err.(GoPackage); ok {
		return pkg.Path, nil
	}
	// Find possible matches from list of seen packages
	var pkgs = make([]string, len(w.seen))
	var idx int
	for k := range w.seen {
		pkgs[idx] = k
		idx++
	}
	matches := fuzzy.RankFindFold(w.find, pkgs)
	if len(matches) == 0 {
		return "", ErrNoMatch
	}
	return "", ErrNoDirectMatch{
		matches: matches,
	}
}

// Gopath attempts to get the currently used GOPATH/src.
func gopath() (string, error) {
	// Try to use GOPATH by default
	if path := os.Getenv("GOPATH"); path != "" {
		return filepath.Join(path, "src"), nil
	}
	// Otherwise use the system default.
	path := filepath.Join(build.Default.GOPATH, "src")
	_, err := os.Stat(path)
	return path, err
}

func main() {
	log.SetFlags(0)

	path, err := gopath()
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	// If no path supplied then change directory to GOPATH.
	if flag.NArg() == 0 {
		fmt.Print(path)
		return
	}

	w := PkgFinder{
		gopath: path,
		find:   flag.Arg(0),
		seen:   make(map[string]struct{}),
	}
	pkg, err := w.Find()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(pkg)
}
