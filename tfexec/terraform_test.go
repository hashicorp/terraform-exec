package tfexec

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

var tfCache *testutil.TFCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var err error
		installDir, err := ioutil.TempDir("", "tfinstall")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(installDir)

		tfCache = testutil.NewTFCache(installDir)

		return m.Run()
	}())
}

func TestSetEnv(t *testing.T) {
	td := t.TempDir()

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

func TestSetLog(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))

	if err != nil {
		t.Fatalf("unexpected NewTerraform error: %s", err)
	}

	// Required so all testing environment variables are not copied.
	err = tf.SetEnv(map[string]string{
		"CLEARENV": "1",
	})

	if err != nil {
		t.Fatalf("unexpected SetEnv error: %s", err)
	}

	t.Run("case 1: SetLog <= 0.15 error", func(t *testing.T) {
		if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
			t.Skip("Terraform for darwin/arm64 is not available until v1")
		}

		td012 := t.TempDir()

		tf012, err := NewTerraform(td012, tfVersion(t, testutil.Latest012))

		if err != nil {
			t.Fatalf("unexpected NewTerraform error: %s", err)
		}

		err = tf012.SetLog("TRACE")

		if err == nil {
			t.Fatal("expected SetLog error, got none")
		}
	})

	t.Run("case 2: SetLog TRACE no SetLogPath", func(t *testing.T) {
		err := tf.SetLog("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV": "1",
			"TF_LOG":   "",
		}, initCmd)
	})

	t.Run("case 3: SetLog TRACE and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":    "1",
			"TF_LOG":      "TRACE",
			"TF_LOG_PATH": tfLogPath,
		}, initCmd)
	})

	t.Run("case 4: SetLog DEBUG and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("DEBUG")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":    "1",
			"TF_LOG":      "DEBUG",
			"TF_LOG_PATH": tfLogPath,
		}, initCmd)
	})
}

func TestSetLogPath(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))

	if err != nil {
		t.Fatalf("unexpected NewTerraform error: %s", err)
	}

	// Required so all testing environment variables are not copied.
	err = tf.SetEnv(map[string]string{
		"CLEARENV": "1",
	})

	if err != nil {
		t.Fatalf("unexpected SetEnv error: %s", err)
	}

	t.Run("case 1: No SetLogPath", func(t *testing.T) {
		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":    "1",
			"TF_LOG":      "",
			"TF_LOG_PATH": "",
		}, initCmd)
	})

	t.Run("case 2: SetLogPath sets TF_LOG and TF_LOG_PATH", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err = tf.SetLogPath(tfLogPath)

		if err != nil {
			t.Fatalf("unexpected SetLogPath error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-force-copy",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":    "1",
			"TF_LOG":      "TRACE",
			"TF_LOG_PATH": tfLogPath,
		}, initCmd)
	})
}

func TestCheckpointDisablePropagation(t *testing.T) {
	td := t.TempDir()

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
			"-upgrade=false",
			"-lock=true",
			"-get-plugins=true",
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
			"-upgrade=false",
			"-lock=true",
			"-get-plugins=true",
			"-verify-plugins=true",
		}, map[string]string{
			"CHECKPOINT_DISABLE": "",
			"FOOBAR":             "2",
		}, initCmd)
	})
}

// test that a suitable error is returned if NewTerraform is called without a valid
// executable path
func TestNoTerraformBinary(t *testing.T) {
	td := t.TempDir()

	_, err := NewTerraform(td, "")
	if err == nil {
		t.Fatal("expected NewTerraform to error, but it did not")
	}

	var e *ErrNoSuitableBinary
	if !errors.As(err, &e) {
		t.Fatal("expected error to be ErrNoSuitableBinary")
	}
}

func tfVersion(t *testing.T, v string) string {
	if tfCache == nil {
		t.Fatalf("tfCache not yet configured, TestMain must run first")
	}

	return tfCache.Version(t, v)
}
