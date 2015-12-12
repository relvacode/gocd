package main

import (
	"errors"
	"fmt"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Get a package from its full import path
func GetPackage(pkgName string) (string, error) {
	pkg := path.Join(os.Getenv("GOPATH"), "src", pkgName)
	i, e := os.Stat(pkg)
	if e != nil {
		return "", e
	}
	if !i.IsDir() {
		return "", errors.New("not a directory")
	}
	return pkg, nil
}

type GoPackage struct {
	Path string
	Name string
}

func (g GoPackage) Error() string {
	return ""
}

type PkgWalker struct {
	find            string            // package to find
	PossibleMatches map[string]string // Possible matches from fuzzy search
}

// Implements filepath.Walker.
// Uses GoPackage{} as en error when a package is found.
// If the dirname of the given path fuzzy matches the find key then add it to the slice of PossibleMatches.
func (w *PkgWalker) walker(path string, fi os.FileInfo, err error) (e error) {
	if fi.Name() == ".git" || strings.Contains(path, "vendor") {
		return filepath.SkipDir
	}
	if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".go") {
		pkgName := filepath.Base(filepath.Dir(path))
		if pkgName == w.find {
			return GoPackage{Path: filepath.Dir(path), Name: w.find}
		} else if strings.ToLower(pkgName) == strings.ToLower(w.find) {
			if _, ok := w.PossibleMatches[pkgName]; !ok {
				w.PossibleMatches[pkgName] = ""
				return nil
			}
			return nil
		} else if m := fuzzy.FindFold(w.find, []string{pkgName}); len(m) > 0 {
			if _, ok := w.PossibleMatches[pkgName]; !ok {
				w.PossibleMatches[pkgName] = ""
				return nil
			}
		}
	}
	return nil
}

// Find a package by the given key
func (w *PkgWalker) FindPackage(key string) (pkg string, e error) {
	w.find = key
	w.PossibleMatches = make(map[string]string)
	e = filepath.Walk(path.Join(os.Getenv("GOPATH"), "src"), w.walker)
	if _, ok := e.(GoPackage); ok {
		pkg, e = e.(GoPackage).Path, nil
	}
	return
}

// Returns any possible matches as a comma separated string
func (w *PkgWalker) SprintMatches() (s string) {
	for k, _ := range w.PossibleMatches {
		s += fmt.Sprintf("'%s', ", k)
	}
	if len(w.PossibleMatches) > 0 {
		s = s[:len(s)-2]
	}
	return
}

func main() {
	log.SetFlags(0)

	if os.Getenv("GOPATH") == "" {
		log.Fatalln("mising GOPATH from environment")
	}
	if len(os.Args) == 1 {
		fmt.Print(path.Join(os.Getenv("GOPATH"), "src"))
		return
	}

	if pkg, e := GetPackage(os.Args[1]); e == nil {
		fmt.Print(pkg)
		return
	} else {
		w := PkgWalker{}
		pkg, err := w.FindPackage(os.Args[1])
		if err != nil {
			log.Fatalln(e)
		}
		if pkg == "" {
			if len(w.PossibleMatches) > 0 {
				log.Fatalf("cannot locate package %s, maybe you meant %s?", os.Args[1], w.SprintMatches())
			}
			log.Fatalf("cannot locate package %s", os.Args[1])
		}
		fmt.Print(pkg)
		return
	}

}
