package tfexec

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestUpgrade012(t *testing.T) {
	td := testTempDir(t)

	t.Run("defaults", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		upgrade012Cmd, err := tf.upgrade012Cmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"0.12upgrade",
		}, nil, upgrade012Cmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		upgrade012Cmd, err := tf.upgrade012Cmd(context.Background(), Yes(true), Force(true))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"0.12upgrade",
			"-yes",
			"-force",
		}, nil, upgrade012Cmd)
	})

	t.Run("unsupported on 0.13", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		_, err = tf.upgrade012Cmd(context.Background())
		if err == nil {
			t.Fatal("expected old version to fail")
		}

		var expectedErr *ErrVersionMismatch
		if !errors.As(err, &expectedErr) {
			t.Fatalf("error doesn't match: %#v", err)
		}
	})
}
