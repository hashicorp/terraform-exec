package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestTaintCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		taintCmd := tf.taintCmd(context.Background(), "aws_instance.foo")

		assertCmd(t, []string{
			"taint",
			"-no-color",
			"aws_instance.foo",
		}, nil, taintCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		taintCmd := tf.taintCmd(context.Background(), "aws_instance.foo", State("teststate"))

		assertCmd(t, []string{
			"taint",
			"-no-color",
			"-state=teststate",
			"aws_instance.foo",
		}, nil, taintCmd)
	})
}
