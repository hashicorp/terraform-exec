package e2etest

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestVersion(t *testing.T) {
	ctx := context.Background()

	for _, tfv := range []string{
		"0.11.14",
		"0.12.28",
		"0.13.0-beta3",
	} {
		t.Run(tfv, func(t *testing.T) {
			td := testTempDir(t)
			defer os.RemoveAll(td)

			err := copyFile(filepath.Join(testFixtureDir, "basic/main.tf"), td)
			if err != nil {
				t.Fatal(err)
			}

			tf, err := tfexec.NewTerraform(td, tfVersion(t, tfv))
			if err != nil {
				t.Fatal(err)
			}

			err = tf.Init(ctx, tfexec.Lock(false))
			if err != nil {
				t.Fatal(err)
			}

			v, _, err := tf.Version(ctx, false)
			if err != nil {
				t.Fatal(err)
			}
			if v.String() != tfv {
				t.Fatalf("expected version %q, got %q", tfv, v)
			}

			// TODO: test/assert provider info

			// force execution / skip cache as well
			v, _, err = tf.Version(ctx, true)
			if err != nil {
				t.Fatal(err)
			}
			if v.String() != tfv {
				t.Fatalf("expected version %q, got %q", tfv, v)
			}
		})
	}
}
