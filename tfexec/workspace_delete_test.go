package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestWorkspaceDeleteCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		workspaceDeleteCmd, err := tf.workspaceDeleteCmd(context.Background(), "workspace-name")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "delete",
			"-no-color",
			"workspace-name",
		}, nil, workspaceDeleteCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		workspaceDeleteCmd, err := tf.workspaceDeleteCmd(context.Background(), "workspace-name",
			LockTimeout("200s"),
			Force(true),
			Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "delete",
			"-no-color",
			"-force",
			"-lock-timeout=200s",
			"-lock=false",
			"workspace-name",
		}, nil, workspaceDeleteCmd)
	})
}
