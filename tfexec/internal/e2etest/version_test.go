package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

// // adding 0.13.0 here due to the regression fixed in https://github.com/hashicorp/terraform/pull/25811
// "0.13.0",

func TestVersion(t *testing.T) {
	runTest(t, []string{
		testutil.Latest011,
		testutil.Latest012,
		testutil.Latest013,
	}, "basic", func(t *testing.T, tfv string, tf *tfexec.Terraform) {
		ctx := context.Background()

		err := tf.Init(ctx, tfexec.Lock(false))
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
