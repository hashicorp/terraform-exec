// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestWorkspaceShowCmd_v1(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		cmd, err := tf.workspaceShowCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace",
			"show",
			"-no-color",
		}, nil, cmd)
	})

	t.Run("reattach config", func(t *testing.T) {
		cmd, err := tf.workspaceShowCmd(context.Background(), Reattach(map[string]ReattachConfig{
			"registry.terraform.io/hashicorp/examplecloud": {
				Protocol:        "grpc",
				ProtocolVersion: 6,
				Pid:             1234,
				Test:            true,
				Addr: ReattachConfigAddr{
					Network: "unix",
					String:  "/fake_folder/T/plugin123",
				},
			},
		}))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace",
			"show",
			"-no-color",
		}, map[string]string{
			"TF_REATTACH_PROVIDERS": `{"registry.terraform.io/hashicorp/examplecloud":{"Protocol":"grpc","ProtocolVersion":6,"Pid":1234,"Test":true,"Addr":{"Network":"unix","String":"/fake_folder/T/plugin123"}}}`,
		}, cmd)
	})
}
