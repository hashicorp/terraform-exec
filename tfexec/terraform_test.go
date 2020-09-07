package tfexec

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var err error
		installDir, err = ioutil.TempDir("", "tfinstall")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(installDir)

		return m.Run()
	}())
}

func TestSetEnv(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range []struct {
		errManual bool
		name      string
	}{
		{false, "OK_ENV_VAR"},

		{true, "TF_LOG"},
		{true, "TF_VAR_foo"},
	} {
		t.Run(c.name, func(t *testing.T) {
			err = tf.SetEnv(map[string]string{c.name: "foo"})

			if c.errManual {
				var evErr *ErrManualEnvVar
				if !errors.As(err, &evErr) {
					t.Fatalf("expected ErrManualEnvVar, got %T %s", err, err)
				}
			} else {
				if !c.errManual && err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestCheckpointDisablePropagation(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("CHECKPOINT_DISABLE", "1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("CHECKPOINT_DISABLE")

	t.Run("case 1: env var is set in environment and not overridden", func(t *testing.T) {

		err = tf.SetEnv(map[string]string{
			"FOOBAR": "1",
		})
		if err != nil {
			t.Fatal(err)
		}

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
			"-get-plugins=true",
			"-lock=true",
			"-upgrade=false",
			"-verify-plugins=true",
		}, map[string]string{
			"CHECKPOINT_DISABLE": "1",
			"FOOBAR":             "1",
		}, initCmd)
	})

	t.Run("case 2: env var is set in environment and overridden with SetEnv", func(t *testing.T) {
		err = tf.SetEnv(map[string]string{
			"CHECKPOINT_DISABLE": "",
			"FOOBAR":             "2",
		})
		if err != nil {
			t.Fatal(err)
		}

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
			"-get-plugins=true",
			"-lock=true",
			"-upgrade=false",
			"-verify-plugins=true",
		}, map[string]string{
			"CHECKPOINT_DISABLE": "",
			"FOOBAR":             "2",
		}, initCmd)
	})
}

func testTempDir(t *testing.T) string {
	d, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("error creating temporary test directory: %s", err)
	}
	// TODO: add t.Cleanup so we can remove the defers

	return d
}

type installedVersion struct {
	path string
	err  error
}

var (
	installDir           string
	installedVersionLock sync.Mutex
	installedVersions    = map[string]installedVersion{}
)

func tfVersion(t *testing.T, v string) string {
	if installDir == "" {
		t.Fatalf("installDir not yet configured, TestMain must run first")
	}

	installedVersionLock.Lock()
	defer installedVersionLock.Unlock()

	iv, ok := installedVersions[v]
	if !ok {
		dir := filepath.Join(installDir, v)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			t.Fatal(err)
		}
		iv.path, iv.err = tfinstall.Find(context.Background(), tfinstall.ExactVersion(v, dir))
		installedVersions[v] = iv
	}

	if iv.err != nil {
		t.Fatalf("error installing terraform version %q: %s", v, iv.err)
	}

	return iv.path
}
