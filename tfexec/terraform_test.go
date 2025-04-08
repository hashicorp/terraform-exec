// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

var tfCache *testutil.TFCache

func TestMain(m *testing.M) {
	if rawDuration := os.Getenv("MOCK_SLEEP_DURATION"); rawDuration != "" {
		sleepMock(rawDuration)
		return
	}

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

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "TRACE",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "DEBUG",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})
}

func TestSetLogCore(t *testing.T) {
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

	t.Run("case 1: SetLogCore <= 0.15 error", func(t *testing.T) {
		if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
			t.Skip("Terraform for darwin/arm64 is not available until v1")
		}

		td012 := t.TempDir()

		tf012, err := NewTerraform(td012, tfVersion(t, testutil.Latest012))

		if err != nil {
			t.Fatalf("unexpected NewTerraform error: %s", err)
		}

		err = tf012.SetLogCore("TRACE")

		if err == nil {
			t.Fatal("expected SetLogCore error, got none")
		}
	})

	t.Run("case 2: SetLogCore TRACE no SetLogPath", func(t *testing.T) {
		err := tf.SetLogCore("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 3: SetLogCore TRACE and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogCore("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "TRACE",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 4: SetLogCore DEBUG and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogCore("DEBUG")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "DEBUG",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 2: SetLogPath sets TF_LOG (if no TF_LOG_CORE or TF_LOG_PROVIDER) and TF_LOG_PATH", func(t *testing.T) {
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "TRACE",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 3: SetLogPath does not set TF_LOG if TF_LOG_CORE", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogCore("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		err = tf.SetLogProvider("")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "TRACE",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 4: SetLogPath does not set TF_LOG if TF_LOG_PROVIDER", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLog("")

		if err != nil {
			t.Fatalf("unexpected SetLog error: %s", err)
		}

		err = tf.SetLogCore("")

		if err != nil {
			t.Fatalf("unexpected SetLogCore error: %s", err)
		}

		err = tf.SetLogProvider("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "TRACE",
		}, initCmd)
	})
}

func TestSetLogProvider(t *testing.T) {
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

	t.Run("case 1: SetLogProvider <= 0.15 error", func(t *testing.T) {
		if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
			t.Skip("Terraform for darwin/arm64 is not available until v1")
		}

		td012 := t.TempDir()

		tf012, err := NewTerraform(td012, tfVersion(t, testutil.Latest012))

		if err != nil {
			t.Fatalf("unexpected NewTerraform error: %s", err)
		}

		err = tf012.SetLogProvider("TRACE")

		if err == nil {
			t.Fatal("expected SetLogProvider error, got none")
		}
	})

	t.Run("case 2: SetLogProvider TRACE no SetLogPath", func(t *testing.T) {
		err := tf.SetLogProvider("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
		}

		initCmd, err := tf.initCmd(context.Background())

		if err != nil {
			t.Fatalf("unexpected command error: %s", err)
		}

		assertCmd(t, []string{
			"init",
			"-no-color",
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     "",
			"TF_LOG_PROVIDER": "",
		}, initCmd)
	})

	t.Run("case 3: SetLogProvider TRACE and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogProvider("TRACE")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "TRACE",
		}, initCmd)
	})

	t.Run("case 4: SetLogProvider DEBUG and SetLogPath", func(t *testing.T) {
		tfLogPath := filepath.Join(td, "test.log")

		err := tf.SetLogProvider("DEBUG")

		if err != nil {
			t.Fatalf("unexpected SetLogProvider error: %s", err)
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CLEARENV":        "1",
			"TF_LOG":          "",
			"TF_LOG_CORE":     "",
			"TF_LOG_PATH":     tfLogPath,
			"TF_LOG_PROVIDER": "DEBUG",
		}, initCmd)
	})
}

func TestCheckpointDisablePropagation_v012(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("Terraform for darwin/arm64 is not available until v1")
	}

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

func TestCheckpointDisablePropagation_v1(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
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
			"-input=false",
			"-backend=true",
			"-get=true",
			"-upgrade=false",
		}, map[string]string{
			"CHECKPOINT_DISABLE": "",
			"FOOBAR":             "2",
		}, initCmd)
	})
}

func TestGracefulCancellation_interruption(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("graceful cancellation not supported on windows")
	}
	mockExecPath, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	td := t.TempDir()

	tf, err := NewTerraform(td, mockExecPath)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 100*time.Millisecond)
	t.Cleanup(cancelFunc)

	_, _, err = tf.version(ctx)
	if err != nil {
		var exitErr *exec.ExitError
		isExitErr := errors.As(err, &exitErr)
		if isExitErr && exitErr.ProcessState.String() == "signal: interrupt" {
			return
		}
		if isExitErr {
			t.Fatalf("expected interrupt signal, received %q", exitErr)
		}
		t.Fatalf("unexpected command error: %s", err)
	}
}

func TestGracefulCancellation_withDelay(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("graceful cancellation not supported on windows")
	}
	mockExecPath, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	td := t.TempDir()
	tf, err := NewTerraform(td, mockExecPath)
	if err != nil {
		t.Fatal(err)
	}
	tf.SetEnv(map[string]string{
		"MOCK_SLEEP_DURATION": "5s",
	})
	tf.SetLogger(testutil.TestLogger())
	tf.SetWaitDelay(100 * time.Millisecond)

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 100*time.Millisecond)
	t.Cleanup(cancelFunc)

	_, _, err = tf.version(ctx)
	if err != nil {
		var exitErr *exec.ExitError
		isExitErr := errors.As(err, &exitErr)
		if isExitErr && exitErr.ProcessState.String() == "signal: killed" {
			return
		}
		if isExitErr {
			t.Fatalf("expected kill signal, received %q", exitErr)
		}
		t.Fatalf("unexpected command error: %s", err)
	}
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
