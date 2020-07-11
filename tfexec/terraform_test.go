package tfexec

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

const testFixtureDir = "testdata"
const testConfigFileName = "main.tf"
const testStateJsonFileName = "state.json"
const testTerraformStateFileName = "terraform.tfstate"

var tfPath string

func TestMain(m *testing.M) {
	var err error
	td, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(td)

	tfPath, err = tfinstall.Find(tfinstall.LookPath(), tfinstall.LatestVersion(td, true))
	if err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)

}

func TestCheckpointDisablePropagation(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfPath)
	if err != nil {
		t.Fatal(err)
	}

	// case 1: env var is set in environment and not overridden
	os.Setenv("CHECKPOINT_DISABLE", "1")
	defer os.Unsetenv("CHECKPOINT_DISABLE")
	tf.SetEnv(map[string]string{
		"FOOBAR": "1",
	})
	initCmd := tf.InitCmd(context.Background())
	expected := []string{"FOOBAR=1", "CHECKPOINT_DISABLE=1", "TF_LOG="}
	actual := initCmd.Env

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected command env to be %s, but it was %s", expected, actual)
	}

	// case 2: env var is set in environment and overridden with SetEnv
	tf.SetEnv(map[string]string{
		"CHECKPOINT_DISABLE": "",
		"FOOBAR":             "1",
	})
	initCmd = tf.InitCmd(context.Background())
	expected = []string{"CHECKPOINT_DISABLE=", "FOOBAR=1", "TF_LOG="}
	actual = initCmd.Env

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

func copyFile(path string, dstPath string) error {
	srcF, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcF.Close()

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
