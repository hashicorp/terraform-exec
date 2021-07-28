package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestProvidersLockCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		lockCmd := tf.providersLockCmd(context.Background())

		assertCmd(t, []string{
			"providers",
			"lock",
		}, nil, lockCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		lockCmd := tf.providersLockCmd(context.Background(), FSMirror("test"), NetMirror("test"), Platform("linux_amd64"), Provider("workingdir"))

		assertCmd(t, []string{
			"providers",
			"lock",
			"-fs-mirror=test",
			"-net-mirror=test",
			"-platform=linux_amd64",
			"workingdir",
		}, nil, lockCmd)
	})
}
