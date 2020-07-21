package tfexec

import (
	"context"
	"os"
	"testing"
)

func TestShowCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	// defaults
	showCmd := tf.showCmd(context.Background())

	assertCmd(t, []string{
		"show",
		"-json",
		"-no-color",
	}, nil, showCmd)
}
