// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

// comparison functions for tfjson structs used in tests

// diffState returns a human-readable report of the differences between two
// state values. It returns an empty string if the two values are equal.
func diffState(expected *tfjson.State, actual *tfjson.State) string {
	return cmp.Diff(expected, actual, cmpopts.IgnoreFields(tfjson.State{}, "TerraformVersion"), cmpopts.IgnoreFields(tfjson.State{}, "useJSONNumber"))
}

// diffPlan returns a human-readable report of the differences between two
// plan values. It returns an empty string if the two values are equal.
func diffPlan(expected *tfjson.Plan, actual *tfjson.Plan, opts ...cmp.Option) string {
	opts = append(opts, cmpopts.IgnoreFields(tfjson.Plan{}, "TerraformVersion", "useJSONNumber"))

	return cmp.Diff(expected, actual, opts...)
}

// diffSchema returns a human-readable report of the differences between two
// schema values. It returns an empty string if the two values are equal.
func diffSchema(expected *tfjson.ProviderSchemas, actual *tfjson.ProviderSchemas) string {
	return cmp.Diff(expected, actual, ctydebug.CmpOptions)
}
