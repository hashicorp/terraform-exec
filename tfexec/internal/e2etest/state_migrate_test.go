// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestStateMigrateJSON(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(stateMigrateMinVersion) {
			t.Skip("state migrate command not available in this Terraform version")
		}

		// TODO: backend local(path foo) -> backend local(path boo)
		err := tf.StateMigrateJSON(context.Background())
		if err != nil {
			t.Fatalf("error running StateMigrateJSON: %s", err)
		}
	})

	// TODO: upgrade test? Is there a provider to test with?
	// TODO: state store test? Is there a provider to test with?
}
