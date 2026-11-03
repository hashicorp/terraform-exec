// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateIdentitiesCmd(t *testing.T) {
	tf, err := NewTerraform(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateIdentitiesCmd, err := tf.stateIdentitiesCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state", "identities",
			"-no-color",
			"-json",
		}, nil, stateIdentitiesCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateIdentitiesCmd, err := tf.stateIdentitiesCmd(context.Background(), State("teststate"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state", "identities",
			"-no-color",
			"-json",
			"-state=teststate",
		}, nil, stateIdentitiesCmd)
	})

	t.Run("reattach config", func(t *testing.T) {
		stateIdentitiesCmd, err := tf.stateIdentitiesCmd(context.Background(), Reattach(map[string]ReattachConfig{
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
			"state", "identities",
			"-no-color",
			"-json",
		}, map[string]string{
			"TF_REATTACH_PROVIDERS": `{"registry.terraform.io/hashicorp/examplecloud":{"Protocol":"grpc","ProtocolVersion":6,"Pid":1234,"Test":true,"Addr":{"Network":"unix","String":"/fake_folder/T/plugin123"}}}`,
		}, stateIdentitiesCmd)
	})
}
