// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateMigrateJSONCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateMigrateJSONCmd, err := tf.stateMigrateJSONCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"migrate",
			"-no-color",
			"-input=false",
			"-force-copy",
			"-json",
		}, nil, stateMigrateJSONCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateMigrateJSONCmd, err := tf.stateMigrateJSONCmd(context.Background(),
			SourceProviderLockFile("src.lock.hcl"),
			DestinationProviderLockFile("dst.lock.hcl"),
			Upgrade(true))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"migrate",
			"-no-color",
			"-input=false",
			"-force-copy",
			"-json",
			"-source-provider-lock-file=src.lock.hcl",
			"-destination-provider-lock-file=dst.lock.hcl",
			"-upgrade",
		}, nil, stateMigrateJSONCmd)
	})
}
