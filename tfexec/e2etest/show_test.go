package e2etest

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestShow(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := tfexec.NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// copy state and config files into test dir
	err = copyFiles(filepath.Join(testFixtureDir, "basic"), td)
	if err != nil {
		t.Fatalf("error copying files into test dir: %s", err)
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

	tf, err := tfexec.NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, "basic", testTerraformStateFileName), td)

	_, err = tf.Show(context.Background())
	if err == nil {
		t.Fatal("expected Show to error, but it did not")
	} else {
		if _, ok := err.(*tfexec.ErrNoInit); !ok {
			t.Fatalf("expected error %s to be ErrNoInit", err)
		}
	}

}

func TestShow_compatible(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tfPath, err := tfinstall.Find(tfinstall.ExactVersion("0.11.14", td))
	if err != nil {
		t.Fatal(err)
	}

	tf, err := tfexec.NewTerraform(td, tfPath)
	if err != nil {
		t.Fatal(err)
	}

	var mismatch *tfexec.ErrVersionMismatch
	_, err = tf.Show(context.Background())
	if !errors.As(err, &mismatch) {
		t.Fatal("expected version mismatch error")
	}
	if mismatch.Actual != "0.11.14" {
		t.Fatalf("expected version 0.11.14, got %q", mismatch.Actual)
	}
	if mismatch.MinInclusive != "0.12.0" {
		t.Fatalf("expected min 0.12.0, got %q", mismatch.MinInclusive)
	}
	if mismatch.MaxExclusive != "-" {
		t.Fatalf("expected max -, got %q", mismatch.MaxExclusive)
	}
}
