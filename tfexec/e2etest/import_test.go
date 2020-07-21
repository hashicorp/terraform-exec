package e2etest

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestImport(t *testing.T) {
	const (
		expectedID      = "asdlfjksdlfkjsdlfk"
		resourceAddress = "random_string.random_string"
	)
	ctx := context.Background()

	for _, tfv := range []string{
		"0.11.14", // doesn't support show JSON output, but does support import
		"0.12.28",
		"0.13.0-beta3",
	} {
		t.Run(tfv, func(t *testing.T) {
			td := testTempDir(t)
			defer os.RemoveAll(td)

			err := copyFiles(filepath.Join(testFixtureDir, "import"), td)
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

			// Config is unnecessary here since its already the working dir, but just testing an additional flag
			err = tf.Import(ctx, resourceAddress, expectedID, tfexec.DisableBackup(), tfexec.Lock(false), tfexec.Config(td))
			if err != nil {
				t.Fatal(err)
			}

			if strings.HasPrefix(tfv, "0.11.") {
				t.Logf("skipping state assertion for 0.11")
				return
			}

			state, err := tf.Show(ctx)
			if err != nil {
				t.Fatal(err)
			}

			for _, r := range state.Values.RootModule.Resources {
				if r.Address != resourceAddress {
					continue
				}

				raw, ok := r.AttributeValues["id"]
				if !ok {
					t.Fatal("value not found for \"id\" attribute")
				}
				actual, ok := raw.(string)
				if !ok {
					t.Fatalf("unable to cast %T to string: %#v", raw, raw)
				}

				if actual != expectedID {
					t.Fatalf("expected %q, got %q", expectedID, actual)
				}

				// success
				return
			}

			t.Fatalf("imported resource %q not found", resourceAddress)
		})
	}
}
