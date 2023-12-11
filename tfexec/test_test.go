// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestTestCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_6))

	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		testCmd := tf.testCmd(context.Background())

		assertCmd(t, []string{
			"test",
			"-json",
		}, nil, testCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		testCmd := tf.testCmd(context.Background(), TestsDirectory("test"))

		assertCmd(t, []string{
			"test",
			"-json",
			"-tests-directory=test",
		}, nil, testCmd)
	})
}
