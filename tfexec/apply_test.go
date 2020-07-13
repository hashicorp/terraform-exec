package tfexec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, "basic/main.tf"), td)
	if err != nil {
		t.Fatalf("error copying config file into test dir: %s", err)
	}

	err = tf.Init(context.Background())
	if err != nil {
		t.Fatalf("error running Init in test directory: %s", err)
	}

	err = tf.Apply(context.Background())
	if err != nil {
		t.Fatalf("error running Apply: %s", err)
	}
}

func TestApplyCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	applyCmd := tf.ApplyCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), DirOrPlan("testfile"))

	actual := strings.TrimPrefix(cmdString(applyCmd), applyCmd.Path+" ")

	expected := "apply -no-color -auto-approve -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -parallelism=99 -refresh=false -target=target1 -target=target2 -var 'var1=foo' -var 'var2=bar' testfile"

	if actual != expected {
		t.Fatalf("expected arguments of ApplyCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
