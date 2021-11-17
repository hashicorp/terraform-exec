package tfexec

import (
	"context"
	"errors"
	"testing"
)

func TestWorkspaceShowCmd(t *testing.T) {
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
