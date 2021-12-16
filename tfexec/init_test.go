package tfexec

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestInitCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		// defaults
		initCmd, err := tf.initCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-lock-timeout=0s",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
			"-lock=true",
			"-get-plugins=true",
			"-verify-plugins=true",
		}, nil, initCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		initCmd, err := tf.initCmd(context.Background(), Backend(false), BackendConfig("confpath1"), BackendConfig("confpath2"), FromModule("testsource"), Get(false), GetPlugins(false), Lock(false), LockTimeout("999s"), PluginDir("testdir1"), PluginDir("testdir2"), Reconfigure(true), Upgrade(true), VerifyPlugins(false), Dir("initdir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-from-module=testsource",
			"-lock-timeout=999s",
			"-backend=false",
			"-get=false",
			"-upgrade=true",
			"-lock=false",
			"-get-plugins=false",
			"-verify-plugins=false",
			"-reconfigure",
			"-backend-config=confpath1",
			"-backend-config=confpath2",
			"-plugin-dir=testdir1",
			"-plugin-dir=testdir2",
			"initdir",
		}, nil, initCmd)
	})
}

func TestInitCmd_compatible(t *testing.T) {
	// Options -lock, -lock-timeout, -verify-plugins, and -get-plugins were
	// removed in 0.15.
	//
	// The -lock and -lock-timeout options were then reinstated in 1.0.10.
	//
	// We do some grey box testing here to come up with meaningful test
	// combinations of options.

	td := t.TempDir()

	t.Run("options not compatible with 0.15", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest015))
		if err != nil {
			t.Fatal(err)
		}

		var expectedErr *ErrVersionMismatch

		_, err := tf.initCmd(
			context.Background(),
			VerifyPlugins(false),
		)
		if err == nil {
			t.Fatalf("expected -verify-plugins to be unsupported in 0.15")
		}
		if !errors.As(err, &expectedErr) {
			t.Fatalf("expected %#v but got %#v", expectedErr, err)
		}

		// TODO getplugins

		// TODO lock

		// TODO locktimeout
	})

	t.Run("options reinstated in 1.0.10", func(t *testing.T) {

	})

	t.Run("options not reinstated in 1.0.10", func(t *testing.T) {

	})
}
