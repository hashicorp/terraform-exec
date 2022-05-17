package tfexec

import (
	"context"
	"errors"
	"runtime"
	"testing"
)

func TestWorkspaceShowCmd(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("Terraform for darwin/arm64 is not available until v1")
	}

	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, "0.9.11"))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("too old version", func(t *testing.T) {
		_, err := tf.workspaceShowCmd(context.Background())
		if err == nil {
			t.Fatal("expected old version to fail")
		}

		var expectedErr *ErrVersionMismatch
		if !errors.As(err, &expectedErr) {
			t.Fatalf("error doesn't match: %#v", err)
		}
	})
}
