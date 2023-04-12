// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestUpgrade012(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("Terraform for darwin/arm64 is not available until v1")
	}

	td := t.TempDir()

	t.Run("defaults", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		upgrade012Cmd, err := tf.upgrade012Cmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"0.12upgrade",
			"-no-color",
			"-yes",
		}, nil, upgrade012Cmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
		if err != nil {
			t.Fatal(err)
		}

		// empty env, to avoid environ mismatch in testing
		tf.SetEnv(map[string]string{})

		upgrade012Cmd, err := tf.upgrade012Cmd(context.Background(), Force(true), Dir("upgrade012dir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"0.12upgrade",
			"-no-color",
			"-yes",
			"-force",
			"upgrade012dir",
		}, nil, upgrade012Cmd)
	})

	unsupportedVersions := []string{
		testutil.Latest011,
		testutil.Latest013,
	}
	for _, tfv := range unsupportedVersions {
		t.Run(fmt.Sprintf("unsupported on %s", tfv), func(t *testing.T) {
			tf, err := NewTerraform(td, tfVersion(t, tfv))
			if err != nil {
				t.Fatal(err)
			}

			// empty env, to avoid environ mismatch in testing
			tf.SetEnv(map[string]string{})

			_, err = tf.upgrade012Cmd(context.Background())
			if err == nil {
				t.Fatalf("expected unsupported version %s to fail", tfv)
			}

			var expectedErr *ErrVersionMismatch
			if !errors.As(err, &expectedErr) {
				t.Fatalf("error doesn't match: %#v", err)
			}
		})
	}
}
