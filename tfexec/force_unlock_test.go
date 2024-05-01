// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestForceUnlockCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		forceUnlockCmd, err := tf.forceUnlockCmd(context.Background(), "12345")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"force-unlock",
			"-no-color",
			"-force",
			"12345",
		}, nil, forceUnlockCmd)
	})
}

// The optional final positional [DIR] argument is available
// until v0.15.0.
func TestForceUnlockCmd_pre015(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("Terraform for darwin/arm64 is not available until v1")
	}

	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest014))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("override all defaults", func(t *testing.T) {
		forceUnlockCmd, err := tf.forceUnlockCmd(context.Background(), "12345", Dir("mydir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"force-unlock",
			"-no-color",
			"-force",
			"12345",
			"mydir",
		}, nil, forceUnlockCmd)
	})
}
