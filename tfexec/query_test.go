// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestQueryJSONCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_Alpha_v1_14))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		queryCmd, err := tf.queryJSONCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"query",
			"-no-color",
			"-json",
		}, nil, queryCmd)
	})

	t.Run("override all", func(t *testing.T) {
		queryCmd, err := tf.queryJSONCmd(context.Background(),
			GenerateConfigOut("generated.tf"),
			Var("android=paranoid"),
			Var("brain_size=planet"),
			VarFile("trillian"),
			Dir("earth"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"query",
			"-no-color",
			"-generate-config-out=generated.tf",
			"-var-file=trillian",
			"-var", "android=paranoid",
			"-var", "brain_size=planet",
			"-json",
			"earth",
		}, nil, queryCmd)
	})
}
