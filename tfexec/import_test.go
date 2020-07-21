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

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	importCmd := tf.importCmd(context.Background(), "my-addr", "my-id")

	actual := strings.TrimPrefix(cmdString(importCmd), importCmd.Path+" ")

	expected := "import -no-color -input=false -lock-timeout=0s -lock=true my-addr my-id"

	if actual != expected {
		t.Fatalf("expected default arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	importCmd = tf.importCmd(context.Background(), "my-addr2", "my-id2",
		Backup("testbackup"),
		LockTimeout("200s"),
		State("teststate"),
		StateOut("teststateout"),
		VarFile("testvarfile"),
		Lock(false),
		Var("var1=foo"),
		Var("var2=bar"),
		AllowMissingConfig(true),
	)

	actual = strings.TrimPrefix(cmdString(importCmd), importCmd.Path+" ")

	expected = "import -no-color -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -allow-missing-config -var 'var1=foo' -var 'var2=bar' my-addr2 my-id2"

	if actual != expected {
		t.Fatalf("expected arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
