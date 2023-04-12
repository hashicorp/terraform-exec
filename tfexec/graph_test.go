// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestGraphCmd_v013(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("Terraform for darwin/arm64 is not available until v1")
	}

	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest013))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		graphCmd, _ := tf.graphCmd(context.Background())

		assertCmd(t, []string{
			"graph",
		}, nil, graphCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		graphCmd, _ := tf.graphCmd(context.Background(),
			GraphPlan("teststate"),
			DrawCycles(true),
			GraphType("output"))

		assertCmd(t, []string{
			"graph",
			"teststate",
			"-draw-cycles",
			"-type=output",
		}, nil, graphCmd)
	})
}

func TestGraphCmd_v1(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		graphCmd, _ := tf.graphCmd(context.Background())

		assertCmd(t, []string{
			"graph",
		}, nil, graphCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		graphCmd, _ := tf.graphCmd(context.Background(),
			GraphPlan("teststate"),
			DrawCycles(true),
			GraphType("output"))

		assertCmd(t, []string{
			"graph",
			"-plan=teststate",
			"-draw-cycles",
			"-type=output",
		}, nil, graphCmd)
	})
}
