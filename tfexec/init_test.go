package tfexec

import (
	"context"
	"os"
	"testing"
)

func TestInitCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		// defaults
		initCmd := tf.initCmd(context.Background())
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
			"-get-plugins=true",
			"-lock=true",
			"-upgrade=false",
			"-verify-plugins=true",
		}, nil, initCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		initCmd := tf.initCmd(context.Background(), Backend(false), BackendConfig("confpath1"), BackendConfig("confpath2"), FromModule("testsource"), Get(false), GetPlugins(false), Lock(false), LockTimeout("999s"), PluginDir("testdir1"), PluginDir("testdir2"), Reconfigure(true), Upgrade(true), VerifyPlugins(false))

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-from-module=testsource",
			"-lock-timeout=999s",
			"-backend=false",
			"-get=false",
			"-get-plugins=false",
			"-lock=false",
			"-upgrade=true",
			"-verify-plugins=false",
			"-reconfigure",
			"-backend-config=confpath1",
			"-backend-config=confpath2",
			"-plugin-dir=testdir1",
			"-plugin-dir=testdir2",
		}, nil, initCmd)
	})
}
