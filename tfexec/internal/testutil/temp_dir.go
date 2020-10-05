// +build go1.14

package testutil

import (
	"io/ioutil"
	"os"
	"testing"
)

// TODO: Remove once we drop support for Go <1.15
// in favour of native t.TempDir()
func TempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("error creating temporary test directory: %s", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	return dir
}
