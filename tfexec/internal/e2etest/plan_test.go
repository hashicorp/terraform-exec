package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestPlan(t *testing.T) {
	tf, cleanup := setupFixture(t, testutil.Latest012, "basic")
	defer cleanup()

	err := tf.Init(context.Background())
	if err != nil {
		t.Fatalf("error running Init in test directory: %s", err)
	}

	err = tf.Plan(context.Background())
	if err != nil {
		t.Fatalf("error running Plan: %s", err)
	}
}
