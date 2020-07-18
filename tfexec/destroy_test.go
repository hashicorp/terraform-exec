package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestDestroyCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	destroyCmd := tf.destroyCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(destroyCmd), destroyCmd.Path+" ")

	expected := "destroy -no-color -auto-approve -input=false -lock-timeout=0s -lock=true -parallelism=10 -refresh=true"

	if actual != expected {
		t.Fatalf("expected default arguments of DestroyCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	destroyCmd = tf.destroyCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"))

	actual = strings.TrimPrefix(cmdString(destroyCmd), destroyCmd.Path+" ")

	expected = "destroy -no-color -auto-approve -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -parallelism=99 -refresh=false -target=target1 -target=target2 -var 'var1=foo' -var 'var2=bar'"

	if actual != expected {
		t.Fatalf("expected arguments of DestroyCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
