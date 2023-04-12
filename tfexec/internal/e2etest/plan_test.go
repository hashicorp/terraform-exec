// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"io"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestPlan(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		hasChanges, err := tf.Plan(context.Background())
		if err != nil {
			t.Fatalf("error running Plan: %s", err)
		}
		if !hasChanges {
			t.Fatalf("expected: true, got: %t", hasChanges)
		}
	})

}

func TestPlanWithState(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(providerAddressMinVersion) {
			t.Skip("state file provider FQNs not compatible with this Terraform version")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		hasChanges, err := tf.Plan(context.Background())
		if err != nil {
			t.Fatalf("error running Plan: %s", err)
		}
		if hasChanges {
			t.Fatalf("expected: false, got: %t", hasChanges)
		}
	})
}

func TestPlanJSON_TF014AndEarlier(t *testing.T) {
	versions := []string{testutil.Latest011, testutil.Latest012, testutil.Latest013, testutil.Latest014}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		re := regexp.MustCompile("terraform plan -json was added in 0.15.3")

		hasChanges, err := tf.PlanJSON(context.Background(), io.Discard)
		if err != nil && !re.MatchString(err.Error()) {
			t.Fatalf("error running Apply: %s", err)
		}
		if hasChanges {
			t.Fatalf("expected: false, got: %t", hasChanges)
		}
	})
}

func TestPlanJSON_TF015AndLater(t *testing.T) {
	versions := []string{testutil.Latest015, testutil.Latest_v1, testutil.Latest_v1_1}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		hasChanges, err := tf.PlanJSON(context.Background(), io.Discard)
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}
		if !hasChanges {
			t.Fatalf("expected: true, got: %t", hasChanges)
		}
	})
}
