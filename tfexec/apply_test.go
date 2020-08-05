package tfexec

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestApplyCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("basic", func(t *testing.T) {
		applyCmd := tf.applyCmd(context.Background(),
			Backup("testbackup"),
			LockTimeout("200s"),
			State("teststate"),
			StateOut("teststateout"),
			VarFile("testvarfile"),
			Lock(false),
			Parallelism(99),
			Refresh(false),
			Target("target1"),
			Target("target2"),
			Var("var1=foo"),
			Var("var2=bar"),
			DirOrPlan("testfile"),
		)

		assertCmd(t, []string{
			"apply",
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
			"testfile",
		}, nil, applyCmd)
	})
}
