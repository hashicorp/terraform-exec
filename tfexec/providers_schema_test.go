package tfexec

import (
	"context"
	"os"
	"testing"
)

func TestProvidersSchemaCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	schemaCmd := tf.providersSchemaCmd(context.Background())

	assertCmd(t, []string{
		"providers",
		"schema",
		"-json",
		"-no-color",
	}, nil, schemaCmd)
}
