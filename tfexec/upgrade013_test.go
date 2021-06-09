package tfexec

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestUpgrade013(t *testing.T) {
	td := testTempDir(t)

	t.Run("defaults", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		upgrade013Cmd, err := tf.upgrade013Cmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"0.13upgrade",
			"-no-color",
			"-yes",
		}, nil, upgrade013Cmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		upgrade013Cmd, err := tf.upgrade013Cmd(context.Background(), Dir("upgrade013dir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"0.13upgrade",
			"-no-color",
			"-yes",
			"upgrade013dir",
		}, nil, upgrade013Cmd)
	})

	unsupportedVersions := []string{
		testutil.Latest011,
		testutil.Latest012,
		testutil.Latest014,
		testutil.Latest015,
	}
	for _, tfv := range unsupportedVersions {
		t.Run(fmt.Sprintf("unsupported on %s", tfv), func(t *testing.T) {
			tf, err := NewTerraform(td, tfVersion(t, tfv))
			if err != nil {
				t.Fatal(err)
			}

			// empty env, to avoid environ mismatch in testing
			tf.SetEnv(map[string]string{})

			_, err = tf.upgrade013Cmd(context.Background())
			if err == nil {
				t.Fatalf("expected unsupported version %s to fail", tfv)
			}

			var expectedErr *ErrVersionMismatch
			if !errors.As(err, &expectedErr) {
				t.Fatalf("error doesn't match: %#v", err)
			}
		})
	}
}
