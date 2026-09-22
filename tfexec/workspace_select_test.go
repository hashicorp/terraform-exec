// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestWorkspaceSelectCmd(t *testing.T) {
	tf, err := NewTerraform(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		workspaceSelectCmd, err := tf.workspaceSelectCmd(context.Background(), "workspace-name")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "select",
			"-no-color",
			"workspace-name",
		}, nil, workspaceSelectCmd)
	})

	t.Run("reattach config", func(t *testing.T) {
		workspaceSelectCmd, err := tf.workspaceSelectCmd(context.Background(), "workspace-name", Reattach(map[string]ReattachConfig{
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
			"workspace", "select",
			"-no-color",
			"workspace-name",
		}, map[string]string{
			"TF_REATTACH_PROVIDERS": `{"registry.terraform.io/hashicorp/examplecloud":{"Protocol":"grpc","ProtocolVersion":6,"Pid":1234,"Test":true,"Addr":{"Network":"unix","String":"/fake_folder/T/plugin123"}}}`,
		}, workspaceSelectCmd)
	})
}

func TestWorkspaceSelectOrCreate(t *testing.T) {
	for _, tc := range []struct {
		name      string
		version   string
		enabled   bool
		wantError bool
	}{
		{"enabled", "1.4.0", true, false},
		{"disabled", "1.4.0", false, false},
		{"unsupported", "1.3.9", true, true},
		{"disabled on older version", "1.3.9", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tf, err := NewTerraform(t.TempDir(), "terraform")
			if err != nil {
				t.Fatal(err)
			}
			tf.execVersion = mustVersion(t, tc.version)
			tf.SetEnv(map[string]string{})
			cmd, err := tf.workspaceSelectCmd(context.Background(), "workspace-name", OrCreate(tc.enabled))
			if tc.wantError {
				var mismatch *ErrVersionMismatch
				if !errors.As(err, &mismatch) {
					t.Fatalf("expected version mismatch, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			args := []string{"workspace", "select", "-no-color"}
			if tc.enabled {
				args = append(args, "-or-create")
			}
			args = append(args, "workspace-name")
			assertCmd(t, args, nil, cmd)
		})
	}
}
