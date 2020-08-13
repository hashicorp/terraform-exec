package e2etest

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestImport(t *testing.T) {
	const (
		expectedID      = "asdlfjksdlfkjsdlfk"
		resourceAddress = "random_string.random_string"
	)

	runTest(t, []string{
		testutil.Latest011, // doesn't support show JSON output, but does support import
		testutil.Latest012,
		testutil.Latest013,
	}, "import", func(t *testing.T, tfv string, tf *tfexec.Terraform) {
		ctx := context.Background()

		err := tf.Init(ctx, tfexec.Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		// Config is unnecessary here since its already the working dir, but just testing an additional flag
		err = tf.Import(ctx, resourceAddress, expectedID, tfexec.DisableBackup(), tfexec.Lock(false), tfexec.Config(tf.WorkingDir()))
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
