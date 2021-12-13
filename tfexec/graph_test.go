package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestGraphCmd(t *testing.T) {
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

func TestGraphCmd15(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest015))
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
