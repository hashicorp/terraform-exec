package tfexec

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

const testFixtureDir = "testdata"
const testTerraformStateFileName = "terraform.tfstate"

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

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
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

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// case 1: env var is set in environment and not overridden
	err = os.Setenv("CHECKPOINT_DISABLE", "1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("CHECKPOINT_DISABLE")

	tf.SetEnv(map[string]string{
		"FOOBAR": "1",
	})
	initCmd := tf.initCmd(context.Background())
	expected := []string{"CHECKPOINT_DISABLE=1", "FOOBAR=1", "TF_LOG="}
	s := initCmd.Env
	sort.Strings(s)
	actual := s

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected command env to be %s, but it was %s", expected, actual)
	}

	// case 2: env var is set in environment and overridden with SetEnv
	err = tf.SetEnv(map[string]string{
		"CHECKPOINT_DISABLE": "",
		"FOOBAR":             "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	initCmd = tf.initCmd(context.Background())
	expected = []string{"CHECKPOINT_DISABLE=", "FOOBAR=1", "TF_LOG="}
	s = initCmd.Env
	sort.Strings(s)
	actual = s

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected command env to be %s, but it was %s", expected, actual)
	}
}

func testTempDir(t *testing.T) string {
	d, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("error creating temporary test directory: %s", err)
	}

	return d
}

func copyFiles(path string, dstPath string) error {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			// TODO: make recursive with filepath.Walk?
			continue
		}
		err = copyFile(filepath.Join(path, info.Name()), dstPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(path string, dstPath string) error {
	srcF, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcF.Close()

	di, err := os.Stat(dstPath)
	if err != nil {
		return err
	}
	if di.IsDir() {
		_, file := filepath.Split(path)
		dstPath = filepath.Join(dstPath, file)
	}

	dstF, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
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
		iv.path, iv.err = tfinstall.Find(tfinstall.ExactVersion(v, dir))
		installedVersions[v] = iv
	}

	if iv.err != nil {
		t.Fatalf("error installing terraform version %q: %s", v, iv.err)
	}

	return iv.path
}
