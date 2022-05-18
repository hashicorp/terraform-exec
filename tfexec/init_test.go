package tfexec

import (
	"context"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestInitCmd_v012(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("Terraform for darwin/arm64 is not available until v1")
	}

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

func TestInitCmd_v1(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
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
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, nil, initCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		initCmd, err := tf.initCmd(context.Background(), Backend(false), BackendConfig("confpath1"), BackendConfig("confpath2"), FromModule("testsource"), Get(false), PluginDir("testdir1"), PluginDir("testdir2"), Reconfigure(true), Upgrade(true), Dir("initdir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-from-module=testsource",
			"-backend=false",
			"-get=false",
			"-upgrade=true",
			"-reconfigure",
			"-backend-config=confpath1",
			"-backend-config=confpath2",
			"-plugin-dir=testdir1",
			"-plugin-dir=testdir2",
			"initdir",
		}, nil, initCmd)
	})
}
