package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStatePull(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	tf.SetEnv(map[string]string{})

	t.Run("tfstate", func(t *testing.T) {
		statePullCmd := tf.statePullCmd(context.Background())

		assertCmd(t, []string{
			"state",
			"pull",
		}, nil, statePullCmd)
	})
}
