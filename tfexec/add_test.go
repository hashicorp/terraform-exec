package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestAddCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("default", func(t *testing.T) {
		addCmd, err := tf.addCmd(context.Background(), "my-addr")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"add",
			"-from-state=false",
			"-optional=false",
			"my-addr",
		}, nil, addCmd)
	})

	t.Run("override-default", func(t *testing.T) {
		addCmd, err := tf.addCmd(context.Background(),
			"my-addr",
			FromState(true),
			Out("out"),
			IncludeOptional(true),
			Provider("my-provider"),
		)
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"add",
			"-from-state=true",
			"-out=out",
			"-optional=true",
			"-provider=my-provider",
			"my-addr",
		}, nil, addCmd)
	})
}
