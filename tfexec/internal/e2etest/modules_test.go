// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestModules(t *testing.T) {
	runTest(t, "deep_module", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(modulesMinVersion) {
			t.Skip("modules command is not available in this Terraform version")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init: %s", err)
		}

		modules, err := tf.Modules(context.Background())
		if err != nil {
			t.Fatalf("error running Modules: %s", err)
		}

		expected := &tfjson.ModulesOutput{
			FormatVersion: "1.0",
			Modules: []tfjson.ModuleOutput{
				{
					Key:     "foo",
					Source:  "./foo",
					Version: "",
				},
				{
					Key:     "foo.bar",
					Source:  "./bar",
					Version: "",
				},
			},
		}

		if diff := cmp.Diff(expected, modules); diff != "" {
			t.Fatalf("unexpected modules output: %s", diff)
		}
	})
}
