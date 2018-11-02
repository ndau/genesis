package generator

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// ensure that the parents of a given path exist
func ensureDir(path string) error {
	parent := filepath.Dir(path)
	return os.MkdirAll(parent, 0700)
}

// Update updates a toml file, clobbering only the data you specify
//
// path is the string of the file's path
// fileContainer should be a pointer to a value you control.
//   This means that you retain the pointer, and can update it from within the
//   updater func.
// if the updater returns an error, the existing file at the path is unchanged,
// and the error is propagated
func Update(path string, fileContainer interface{}, updater func() error) error {
	err := ensureDir(path)
	if err != nil {
		return err
	}

	in, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		// propagate the error only if it's something other than the file not existing
		// if it's that, we don't care; we can create a file
		return err
	}
	defer in.Close()

	_, err = toml.DecodeReader(in, fileContainer)
	if err != nil {
		// we can't do anything about toml-reading errors
		return err
	}

	err = updater()
	if err != nil {
		// user-supplied error
		return err
	}

	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		return err
	}
	return toml.NewEncoder(out).Encode(fileContainer)
}
