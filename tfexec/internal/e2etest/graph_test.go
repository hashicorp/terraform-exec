package e2etest

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestGraph(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply: %s", err)
		}

		graphOutput, err := tf.Graph(context.Background())
		if err != nil {
			t.Fatalf("error running Graph: %s", err)
		}

		// Graph output differs slightly between versions, but resource subgraph remains consistent
		if !strings.Contains(graphOutput, `"[root] null_resource.foo" [label = "null_resource.foo", shape = "box"]`) {
			t.Fatalf("error running Graph. Graph output does not contain expected strings. Returned: %s", graphOutput)
		}
	})
}
