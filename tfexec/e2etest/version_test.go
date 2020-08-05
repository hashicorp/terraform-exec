package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestVersion(t *testing.T) {
	ctx := context.Background()

	for _, tfv := range []string{
		testutil.Latest011,
		testutil.Latest012,
		testutil.Latest013,
	} {
		t.Run(tfv, func(t *testing.T) {
			tf, cleanup := setupFixture(t, tfv, "basic")
			defer cleanup()

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
}
