package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestOutput(t *testing.T) {
	tf, cleanup := setupFixture(t, testutil.Latest012, "basic")
	defer cleanup()

	err := tf.Init(context.Background())
	if err != nil {
		t.Fatalf("error running Init in test directory: %s", err)
	}

	_, err = tf.Output(context.Background())
	if err != nil {
		t.Fatalf("error running Output: %s", err)
	}
}
