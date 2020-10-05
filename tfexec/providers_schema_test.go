package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestProvidersSchemaCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
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
