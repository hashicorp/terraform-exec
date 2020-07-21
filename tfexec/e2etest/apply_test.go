package e2etest

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestApply(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := tfexec.NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, "basic/main.tf"), td)
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
