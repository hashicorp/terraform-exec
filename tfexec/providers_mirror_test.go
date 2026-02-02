package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestProvidersMirrorCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		mirrorCmd := tf.providersMirrorCmd(context.Background(), "foo")

		assertCmd(t, []string{
			"providers",
			"mirror",
			"foo",
		}, nil, mirrorCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		mirrorCmd := tf.providersMirrorCmd(context.Background(), "foo", Platform("linux_amd64"))

		assertCmd(t, []string{
			"providers",
			"mirror",
			"foo",
			"-platform=linux_amd64",
		}, nil, mirrorCmd)
	})
}
