package util

// ----- ---- --- -- -
// Copyright 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

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
