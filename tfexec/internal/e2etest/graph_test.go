// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
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

		if diff := cmp.Diff(expectedGraphOutput(tfv), graphOutput); diff != "" {
			t.Fatalf("Graph output does not match: %s", diff)
		}
	})
}

func expectedGraphOutput(tfv *version.Version) string {
	v := tfv.Core()

	if v.LessThan(v0_12_0) {
		// TF <=0.11.15
		return `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] null_resource.foo" [label = "null_resource.foo", shape = "box"]
		"[root] provider.null" [label = "provider.null", shape = "diamond"]
		"[root] meta.count-boundary (count boundary fixup)" -> "[root] null_resource.foo"
		"[root] null_resource.foo" -> "[root] provider.null"
		"[root] provider.null (close)" -> "[root] null_resource.foo"
		"[root] root" -> "[root] meta.count-boundary (count boundary fixup)"
		"[root] root" -> "[root] provider.null (close)"
	}
}

`
	}

	if v.GreaterThanOrEqual(v0_12_0) && v.LessThan(v0_13_0) {
		// TF 0.12.20 - 0.12.31
		return `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] null_resource.foo" [label = "null_resource.foo", shape = "box"]
		"[root] provider.null" [label = "provider.null", shape = "diamond"]
		"[root] meta.count-boundary (EachMode fixup)" -> "[root] null_resource.foo"
		"[root] null_resource.foo" -> "[root] provider.null"
		"[root] provider.null (close)" -> "[root] null_resource.foo"
		"[root] root" -> "[root] meta.count-boundary (EachMode fixup)"
		"[root] root" -> "[root] provider.null (close)"
	}
}

`
	}

	if v.GreaterThanOrEqual(v0_13_0) && v.LessThan(v1_1) {
		// 0.13.0 - 1.0.11
		return `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] null_resource.foo (expand)" [label = "null_resource.foo", shape = "box"]
		"[root] provider[\"registry.opentofu.org/hashicorp/null\"]" [label = "provider[\"registry.opentofu.org/hashicorp/null\"]", shape = "diamond"]
		"[root] meta.count-boundary (EachMode fixup)" -> "[root] null_resource.foo (expand)"
		"[root] null_resource.foo (expand)" -> "[root] provider[\"registry.opentofu.org/hashicorp/null\"]"
		"[root] provider[\"registry.opentofu.org/hashicorp/null\"] (close)" -> "[root] null_resource.foo (expand)"
		"[root] root" -> "[root] meta.count-boundary (EachMode fixup)"
		"[root] root" -> "[root] provider[\"registry.opentofu.org/hashicorp/null\"] (close)"
	}
}

`
	}

	// 1.1.0+
	return `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] null_resource.foo (expand)" [label = "null_resource.foo", shape = "box"]
		"[root] provider[\"registry.opentofu.org/hashicorp/null\"]" [label = "provider[\"registry.opentofu.org/hashicorp/null\"]", shape = "diamond"]
		"[root] null_resource.foo (expand)" -> "[root] provider[\"registry.opentofu.org/hashicorp/null\"]"
		"[root] provider[\"registry.opentofu.org/hashicorp/null\"] (close)" -> "[root] null_resource.foo (expand)"
		"[root] root" -> "[root] provider[\"registry.opentofu.org/hashicorp/null\"] (close)"
	}
}

`
}
