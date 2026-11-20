package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateReplaceProviderCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateReplaceProviderCmd, err := tf.stateReplaceProviderCmd(context.Background(), "testfromprovider", "testtoprovider")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"replace-provider",
			"-no-color",
			"-auto-approve",
			"-lock-timeout=0s",
			"-lock=true",
			"testfromprovider",
			"testtoprovider",
		}, nil, stateReplaceProviderCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateReplaceProviderCmd, err := tf.stateReplaceProviderCmd(context.Background(), "testfromprovider", "testtoprovider", Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"replace-provider",
			"-no-color",
			"-auto-approve",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-lock=false",
			"testfromprovider",
			"testtoprovider",
		}, nil, stateReplaceProviderCmd)
	})
}
