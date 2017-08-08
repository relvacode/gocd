package main

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

// PkgFinder finds a Go package.
type PkgFinder struct {
	gopath string // GOPATH to use
}

// Implements filepath.Walker.
// Uses GoPackage{} as en error when a package is found.
// If the dirname of the given path fuzzy matches the find key then add it to the slice of PossibleMatches.
func (w *PkgFinder) walker(seen map[string]struct{}, find string) filepath.WalkFunc {
	return func(path string, i os.FileInfo, err error) (e error) {
		// Skip GOPATH/src
		if path == w.gopath {
			return nil
		}
		// Skip if path contains .git or vendor
		if i.IsDir() && (strings.HasPrefix(i.Name(), ".") || strings.HasPrefix(i.Name(), "_") || strings.Contains(path, "vendor")) {
			return filepath.SkipDir
		}
		// Ignore if path is a directory or is not a go file.
		if i.IsDir() || !strings.HasSuffix(i.Name(), "go") {
			return nil
		}

		// Scan every component of the relative path until we find a direct match.
		pkg, _ := filepath.Rel(w.gopath, filepath.Dir(path))

		// Skip already seen packages
		_, ok := seen[pkg]
		if ok {
			return nil
		}
		seen[pkg] = struct{}{}

		components := strings.Split(pkg, string(filepath.Separator))
		for x := len(components) - 1; x >= 0; x-- {
			if find == filepath.Join(components[x:]...) {
				return GoPackage{
					Path: filepath.Dir(path),
					Name: pkg,
				}
			}
		}
		return nil
	}
}

// Find a package by the given key
func (w *PkgFinder) Find(find string) (fuzzy.Ranks, error) {
	// If absolute path then go straight to it
	if path.IsAbs(find) {
		return fuzzy.Ranks{
			{
				Target: find,
			},
		}, nil
	}

	// If path is a path relative to gopath then use it
	abs := filepath.Join(w.gopath, find)
	_, err := os.Stat(abs)
	if err == nil {
		return fuzzy.Ranks{
			{
				Target: abs,
			},
		}, nil
	}

	seen := make(map[string]struct{})

	err = filepath.Walk(w.gopath, w.walker(seen, find))
	if pkg, ok := err.(GoPackage); ok {
		return fuzzy.Ranks{
			{
				Target: pkg.Path,
			},
		}, nil
	}
	// Find possible matches from list of seen packages
	var matches = make(map[string]fuzzy.Rank)
	for k := range seen {

		path := append(strings.Split(k, string(filepath.Separator)), k)
		ranks := fuzzy.RankFindFold(find, path)

		for _, r := range ranks {
			if r.Distance > 10 {
				continue
			}
			m, ok := matches[k]

			if (ok && r.Distance < m.Distance) || !ok {
				r.Target = filepath.Join(w.gopath, k)
				matches[k] = r
			}
		}
	}

	if len(matches) == 0 {
		return nil, ErrNoMatch
	}

	found := make(fuzzy.Ranks, len(matches))
	var i int
	for _, r := range matches {
		found[i] = r
		i++
	}

	sort.Sort(found)

	if len(found) > 10 {
		found = found[:10]
	}

	return found, nil
}
