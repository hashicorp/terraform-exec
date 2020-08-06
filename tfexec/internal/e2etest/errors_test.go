// This file contains tests that only compile/work in Go 1.13 and forward
// +build go1.13

package e2etest

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestUnparsedError(t *testing.T) {
	// This simulates an unparsed error from the Cmd.Run method (in this case file not found). This
	// is to ensure we don't miss raising unexpected errors in addition to parsed / well known ones.
	for _, tfv := range []string{
		testutil.Latest011,
		testutil.Latest012,
		testutil.Latest013,
	} {
		t.Run(tfv, func(t *testing.T) {
			tf, cleanup := setupFixture(t, tfv, "")
			defer cleanup()

			// force delete the working dir to cause an os.PathError
			err := os.RemoveAll(tf.WorkingDir())
			if err != nil {
				t.Fatal(err)
			}

			err = tf.Init(context.Background())
			if err == nil {
				t.Fatalf("expected error running Init, none returned")
			}
			var e *os.PathError
			if !errors.As(err, &e) {
				t.Fatalf("expected os.PathError, got %T, %s", err, err)
			}
		})
	}
}
