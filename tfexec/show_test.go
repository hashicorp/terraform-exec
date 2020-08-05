package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestShowCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
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
