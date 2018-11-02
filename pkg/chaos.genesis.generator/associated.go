package generator

import "path/filepath"

// DefaultAssociated returns the default path to the associated data
func DefaultAssociated(ndauhome string) string {
	return filepath.Join(ndauhome, "chaos", "associated.toml")
}

// Associated tracks associated data which goes with the mocks.
//
// In particular, it's used for tests. For example, we mock up some
// public/private keypairs for the ReleaseFromEndowment transaction.
// The public halves of those keys are written into the mock file,
// but the private halves are communicated to the test suite by means
// of the Associated struct.
type Associated map[string]interface{}

// AssociatedFile is a file format which tracks associated data over time.
//
// In order that we never clobber old data, we namespace the associated data
// by the BPC public key, base64-encoded.
type AssociatedFile map[string]Associated
