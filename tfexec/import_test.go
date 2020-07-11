package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestImportCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfPath)
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	importCmd := tf.ImportCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(importCmd), importCmd.Path+" ")

	expected := "import -no-color -input=false -lock-timeout=0s -lock=true"

	if actual != expected {
		t.Fatalf("expected default arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	importCmd = tf.ImportCmd(context.Background(),
		Backup("testbackup"),
		LockTimeout("200s"),
		State("teststate"),
		StateOut("teststateout"),
		VarFile("testvarfile"),
		Lock(false),
		Var("var1=foo"),
		Var("var2=bar"),
		AllowMissingConfig(true))

	actual = strings.TrimPrefix(cmdString(importCmd), importCmd.Path+" ")

	expected = "import -no-color -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -allow-missing-config -var 'var1=foo' -var 'var2=bar'"

	if actual != expected {
		t.Fatalf("expected arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
