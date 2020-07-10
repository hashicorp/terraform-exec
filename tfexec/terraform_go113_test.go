// This file contains tests that only compile/work in Go 1.13 and forward
// +build go1.13

package tfexec

import (
	"errors"
	"testing"
)

// test that a suitable error is returned if NewTerraform is called without a valid
// executable path
func TestNoTerraformBinary(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	_, err := NewTerraform(td, "")
	if err == nil {
		t.Fatal("expected NewTerraform to error, but it did not")
	}

	var e *ErrNoSuitableBinary
	if !errors.As(err, &e) {
		t.Fatal("expected error to be ErrNoSuitableBinary")
	}
}
