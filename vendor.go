package main

import (
	"flag"
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

const (
	// VendorToken is the argument token used to indicate that gocd
	// should change directory to the vendor's parent.
	VendorToken = "^"
)

func TryGoToVendorParent() (bool, error) {
	arg := flag.Arg(0)
	if arg != VendorToken {
		return false, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	if !strings.Contains(cwd, "vendor") {
		return false, nil
	}

	components := strings.Split(cwd, string(filepath.Separator))
	for i := len(components) - 1; i  >= 0; i -- {
		if components[i] == "vendor" {
			if i == 0 {
				// "vendor" is at the root of the path
				return false, nil
			}

			var abs = append([]string{"/"}, components[:i]...)
			fmt.Print(filepath.Join(abs...))
			return true, nil
		}
	}

	return false, nil
}