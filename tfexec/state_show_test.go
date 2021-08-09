package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateShowCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateShowCmd, err := tf.stateShowCmd(context.Background(), "testaddress")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"show",
			"-no-color",
			"testaddress",
		}, nil, stateShowCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateShowCmd, err := tf.stateShowCmd(context.Background(), "testaddress", State("teststate"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"show",
			"-no-color",
			"-state=teststate",
			"testaddress",
		}, nil, stateShowCmd)
	})
}
