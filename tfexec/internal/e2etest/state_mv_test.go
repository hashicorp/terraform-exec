// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestStateMv(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(providerAddressMinVersion) {
			t.Skip("state file provider FQNs not compatible with this Terraform version")
		}

		providerName := "registry.opentofu.org/hashicorp/null"

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.StateMv(context.Background(), "null_resource.foo", "null_resource.bar")
		if err != nil {
			t.Fatalf("error running StateMv: %s", err)
		}

		formatVersion := "0.1"
		var sensitiveValues json.RawMessage
		if tfv.Core().GreaterThanOrEqual(v1_0_1) {
			formatVersion = "0.2"
			sensitiveValues = json.RawMessage([]byte("{}"))
		}
		if tfv.Core().GreaterThanOrEqual(v1_1) {
			formatVersion = "1.0"
		}

		// test that the new state is as expected
		expected := &tfjson.State{
			FormatVersion: formatVersion,
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.bar",
						AttributeValues: map[string]interface{}{
							"id":       "5510719323588825107",
							"triggers": nil,
						},
						SensitiveValues: sensitiveValues,
						Mode:            tfjson.ManagedResourceMode,
						Type:            "null_resource",
						Name:            "bar",
						ProviderName:    providerName,
					}},
				},
			},
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}
