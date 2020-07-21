package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestShowCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	showCmd := tf.showCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(showCmd), showCmd.Path+" ")

	expected := "show -json -no-color"

	if actual != expected {
		t.Fatalf("expected default arguments of ShowCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
