package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestOutputCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		outputCmd := tf.outputCmd(context.Background())

		assertCmd(t, []string{
			"output",
			"-json",
			"-no-color",
		}, nil, outputCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		outputCmd := tf.outputCmd(context.Background(),
			State("teststate"))

		assertCmd(t, []string{
			"output",
			"-json",
			"-state=teststate",
			"-no-color",
		}, nil, outputCmd)
	})
}
