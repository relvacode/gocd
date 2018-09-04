package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

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

	// Using '^', try to go to the vendor's parent
	ok, err := TryGoToVendorParent()
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		return
	}

	w := PkgFinder{
		gopath: path,
	}

	matches, err := w.Find(flag.Arg(0))

	if err != nil {
		log.Fatal(err)
	}

	if len(matches) == 1 {
		fmt.Print(matches[0].Target)
		return
	}

	if flag.NArg() > 1 {
		i, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.Fatalf("cannot parse requested index %s: %s", flag.Arg(1), err)
		}

		if i > len(matches) {
			log.Fatalf("%d is an invalid index (max %d)", i, len(matches))
		}

		fmt.Printf(matches[i].Target)
		return
	}

	for i, m := range matches {
		rel, _ := filepath.Rel(path, m.Target)
		log.Printf("  %d %s", i, rel)
	}
	os.Exit(1)

}
