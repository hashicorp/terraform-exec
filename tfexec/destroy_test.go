package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestDestroyCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		destroyCmd, err := tf.DestroyCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"destroy",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
		}, nil, destroyCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		destroyCmd, err := tf.DestroyCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), Dir("destroydir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"destroy",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-parallelism=99",
			"-refresh=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"destroydir",
		}, nil, destroyCmd)
	})
}
