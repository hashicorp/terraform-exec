package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestProvidersSchemaCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	schemaCmd := tf.providersSchemaCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(schemaCmd), schemaCmd.Path+" ")

	expected := "providers schema -json -no-color"

	if actual != expected {
		t.Fatalf("expected default arguments of ProvidersSchemaCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
