package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateListCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateListCmd, err := tf.stateListCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"list",
			"-no-color",
		}, nil, stateListCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateListCmd, err := tf.stateListCmd(context.Background(), State("teststate"), Id("testid"), Address("testaddress1"), Address("testaddress2"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"list",
			"-no-color",
			"-state=teststate",
			"-id=testid",
			"testaddress1",
			"testaddress2",
		}, nil, stateListCmd)
	})
}
