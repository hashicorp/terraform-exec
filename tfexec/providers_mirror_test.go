// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestProvidersMirrorCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		mirrorCmd, err := tf.providersMirrorCmd(context.Background(), "path")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"providers",
			"mirror",
			"path",
		}, nil, mirrorCmd)
	})

	t.Run("multiple platforms", func(t *testing.T) {
		mirrorCmd, err := tf.providersMirrorCmd(context.Background(), "path", Platform("IBM-Z"), Platform("Solaris"), Platform("Commodore64"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"providers",
			"mirror",
			"-platform=IBM-Z",
			"-platform=Solaris",
			"-platform=Commodore64",
			"path",
		}, nil, mirrorCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		mirrorCmd, err := tf.providersMirrorCmd(context.Background(), "path", LockFile(false), Platform("IBM-Z"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"providers",
			"mirror",
			"-platform=IBM-Z",
			"-lock-file=false",
			"path",
		}, nil, mirrorCmd)
	})
}
