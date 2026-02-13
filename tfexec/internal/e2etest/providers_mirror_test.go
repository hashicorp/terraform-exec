// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0
package e2etest

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	providersMirrorMinVersion         = version.Must(version.NewVersion("0.13.0"))
	providersMirrorLockFileMinVersion = version.Must(version.NewVersion("1.10.0"))
)

func TestProvidersMirror(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(providersMirrorMinVersion) {
			t.Skip("terraform providers mirror was added in Terraform 0.13, so test is not valid")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		targetDir := t.TempDir()
		err = tf.ProvidersMirror(context.Background(), targetDir)
		if err != nil {
			t.Fatalf("error running providers mirror: %s", err)
		}

		expectedMirrorPath := filepath.Join(targetDir, "registry.terraform.io", "hashicorp", "null")
		_, err = os.Stat(expectedMirrorPath)
		if err != nil {
			t.Fatalf("providers mirror not found in %s", expectedMirrorPath)
		}
	})
}

func TestProvidersMirror_lockFileFalse(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(providersMirrorLockFileMinVersion) {
			t.Skip("terraform providers mirror -lock-file flag was added in Terraform 1.10, so test is not valid")
		}
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		targetDir := t.TempDir()
		err = tf.ProvidersMirror(context.Background(), targetDir, tfexec.LockFile(false))
		if err != nil {
			t.Fatalf("error running providers mirror: %s", err)
		}
	})
}
