package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// Check that the error is nil; exit otherwise
func Check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}

// GetNdauhome gets the value of the $NDAUHOME environment variable
func GetNdauhome() string {
	nh := os.ExpandEnv("$NDAUHOME")
	if len(nh) > 0 {
		return nh
	}
	return filepath.Join(os.ExpandEnv("$HOME"), ".ndau")
}
