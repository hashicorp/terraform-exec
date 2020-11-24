package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateMvCmd(t *testing.T) {
	td := testTempDir(t)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateMvCmd, err := tf.stateMvCmd(context.Background(), "testsource", "testdestination")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"mv",
			"-no-color",
			"-lock-timeout=0s",
			"-lock=true",
			"testsource",
			"testdestination",
		}, nil, stateMvCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateMvCmd, err := tf.stateMvCmd(context.Background(), "testsrc", "testdest", Backup("testbackup"), BackupOut("testbackupout"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"mv",
			"-no-color",
			"-backup=testbackup",
			"-backup-out=testbackupout",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-lock=false",
			"testsrc",
			"testdest",
		}, nil, stateMvCmd)
	})
}
