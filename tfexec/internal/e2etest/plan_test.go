package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestPlan(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		ok, err := tf.Plan(context.Background())
		if err != nil {
			t.Fatalf("error running Plan: %s", err)
		}
		if !ok {
			t.Fatalf("expected: true, got: %t", ok)
		}
	})

}

func TestPlanWithState(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(version.Must(version.NewVersion("0.12.0"))) {
			t.Skip("state file is not compatiable with Terraform 0.11")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		ok, err := tf.Plan(context.Background())
		if err != nil {
			t.Fatalf("error running Plan: %s", err)
		}
		if ok {
			t.Fatalf("expected: false, got: %t", ok)
		}
	})

}
