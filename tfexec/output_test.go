package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestOutputCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	outputCmd := tf.outputCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(outputCmd), outputCmd.Path+" ")

	expected := "output -no-color -json"

	if actual != expected {
		t.Fatalf("expected default arguments of OutputCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	outputCmd = tf.outputCmd(context.Background(),
		State("teststate"))

	actual = strings.TrimPrefix(cmdString(outputCmd), outputCmd.Path+" ")

	expected = "output -no-color -json -state=teststate"

	if actual != expected {
		t.Fatalf("expected arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
