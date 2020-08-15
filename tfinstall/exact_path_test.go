package tfinstall

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// test that Find returns an appropriate error when given an exact path
// which exists, but is not a terraform executable
func TestExactPath(t *testing.T) {
	// we just want the path to a local executable that definitely exists
	execPath, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Find(ExactPath(execPath))
	if err == nil {
		t.Fatalf("expected Find() to fail when given ExactPath(%s), but it did not", execPath)
	}

	expected := fmt.Sprintf("executable found at path %s is not terraform", execPath)
	if !strings.HasPrefix(err.Error(), expected) {
		t.Fatalf("expected Find() to return %s, but got %s", expected, err)
	}
}
