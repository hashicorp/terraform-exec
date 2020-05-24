package tfexec

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

const testFixtureDir = "testdata"
const testConfigFileName = "main.tf"
const testStateJsonFileName = "state.json"
const testTerraformStateFileName = "terraform.tfstate"

func TestShow(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, "")
	if err != nil {
		t.Fatal(err)
	}

	// copy state and config files into test dir
	err = copyFile(filepath.Join(testFixtureDir, testTerraformStateFileName), filepath.Join(td, testTerraformStateFileName))
	if err != nil {
		t.Fatalf("error copying state file into test dir: %s", err)
	}
	err = copyFile(filepath.Join(testFixtureDir, testConfigFileName), filepath.Join(td, testConfigFileName))
	if err != nil {
		t.Fatalf("error copying config file into test dir: %s", err)
	}

	expected := tfjson.State{
		FormatVersion:    "0.1",
		TerraformVersion: "0.12.24",
		Values: &tfjson.StateValues{
			RootModule: &tfjson.StateModule{
				Resources: []*tfjson.StateResource{&tfjson.StateResource{
					Address: "null_resource.foo",
					AttributeValues: map[string]interface{}{
						"id":       "5510719323588825107",
						"triggers": nil,
					},
					Mode:         tfjson.ManagedResourceMode,
					Type:         "null_resource",
					Name:         "foo",
					ProviderName: "null",
				}},
			},
		},
	}

	err = tf.Init(context.Background())
	if err != nil {
		t.Fatalf("error running Init in test directory: %s", err)
	}

	actual, err := tf.Show(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, &expected) {
		t.Fatalf("actual: %s\nexpected: %s", spew.Sdump(actual), spew.Sdump(expected))
	}
}

func TestShow_errInitRequired(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, "")
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, testTerraformStateFileName), filepath.Join(td, testTerraformStateFileName))

	// This test will break if the error output of `terraform init`
	// changes significantly. We tolerate this brittleness as a poor
	// man's canary for significant changes to Terraform CLI.
	// FIXME: Parse this in the actual command and return ErrNoInit
	expected := "Error: Could not satisfy plugin requirements"

	_, err = tf.Show(context.Background())
	if err == nil {
		t.Fatal("expected Show to error, but it did not")
	} else {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("expected error %s to contain %s", err, expected)
		}
	}

}

func TestApply(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, "")
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, testConfigFileName), filepath.Join(td, testConfigFileName))
	if err != nil {
		t.Fatalf("error copying config file into test dir: %s", err)
	}

	err = tf.Init(context.Background())
	if err != nil {
		t.Fatalf("error running Init in test directory: %s", err)
	}

	err = tf.Apply(context.Background())
	if err != nil {
		t.Fatalf("error running Apply: %s", err)
	}
}

func TestApplyCmd(t *testing.T) {
	tf, err := NewTerraform("/dev/null", "")
	if err != nil {
		t.Fatal(err)
	}

	applyCmd := tf.ApplyCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"))

	actual := strings.TrimPrefix(applyCmd.String(), applyCmd.Path+" ")

	expected := "apply -auto-approve -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -parallelism=99 -refresh=false -target=target1 -target=target2 -var 'var1=foo' -var 'var2=bar' -no-color"

	if actual != expected {
		t.Fatalf("expected arguments of ApplyCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
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
