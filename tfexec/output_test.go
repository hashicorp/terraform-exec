package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestOutputCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		outputCmd := tf.outputCmd(context.Background())

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
		}, nil, outputCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		outputCmd := tf.outputCmd(context.Background(),
			State("teststate"))

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
			"-state=teststate",
		}, nil, outputCmd)
	})
}
