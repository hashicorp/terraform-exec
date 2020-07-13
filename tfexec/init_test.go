package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestInitCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfPath)
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	initCmd := tf.InitCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(initCmd), initCmd.Path+" ")

	expected := "init -no-color -force-copy -input=false -lock-timeout=0s -backend=true -get=true -get-plugins=true -lock=true -upgrade=false -verify-plugins=true"

	if actual != expected {
		t.Fatalf("expected default arguments of InitCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	initCmd = tf.InitCmd(context.Background(), Backend(false), BackendConfig("confpath1"), BackendConfig("confpath2"), FromModule("testsource"), Get(false), GetPlugins(false), Lock(false), LockTimeout("999s"), PluginDir("testdir1"), PluginDir("testdir2"), Reconfigure(true), Upgrade(true), VerifyPlugins(false))

	actual = strings.TrimPrefix(cmdString(initCmd), initCmd.Path+" ")

	expected = "init -no-color -force-copy -input=false -from-module=testsource -lock-timeout=999s -backend=false -get=false -get-plugins=false -lock=false -upgrade=true -verify-plugins=false -reconfigure -backend-config=confpath1 -backend-config=confpath2 -plugin-dir=testdir1 -plugin-dir=testdir2"

	if actual != expected {
		t.Fatalf("expected arguments of InitCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
