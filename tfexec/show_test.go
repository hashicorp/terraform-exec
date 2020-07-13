package tfexec

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestStateShow(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfPath)
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
				Resources: []*tfjson.StateResource{{
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

	actual, err := tf.StateShow(context.Background())
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

	tf, err := NewTerraform(td, tfPath)
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, testTerraformStateFileName), filepath.Join(td, testTerraformStateFileName))

	_, err = tf.StateShow(context.Background())
	if err == nil {
		t.Fatal("expected Show to error, but it did not")
	} else {
		if _, ok := err.(*ErrNoInit); !ok {
			t.Fatalf("expected error %s to be ErrNoInit", err)
		}
	}

}

func TestStateShowCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfPath)
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	showCmd := tf.StateShowCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(showCmd), showCmd.Path+" ")

	expected := "show -json -no-color"

	if actual != expected {
		t.Fatalf("expected default arguments of ShowCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
