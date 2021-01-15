package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestVersion(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		ctx := context.Background()

		err := tf.Init(ctx)
		if err != nil {
			t.Fatal(err)
		}

		v, _, err := tf.Version(ctx, false)
		if err != nil {
			t.Fatal(err)
		}
		if !v.Equal(tfv) {
			t.Fatalf("expected version %q, got %q", tfv, v)
		}

		// TODO: test/assert provider info

		// force execution / skip cache as well
		v, _, err = tf.Version(ctx, true)
		if err != nil {
			t.Fatal(err)
		}
		if !v.Equal(tfv) {
			t.Fatalf("expected version %q, got %q", tfv, v)
		}
	})
}
