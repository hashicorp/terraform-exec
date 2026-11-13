// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestModulesCmd_v1_10(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_10))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	modulesCmd, err := tf.modulesCmd(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"modules",
		"-json",
	}, nil, modulesCmd)
}

func TestModulesCmd_v1_15(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_15))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	modulesCmd, err := tf.modulesCmd(context.Background(),
		Var("foo=bar"),
		VarFile("path-to.tfvars"))
	if err != nil {
		t.Fatal(err)
	}

	assertCmd(t, []string{
		"modules",
		"-var-file=path-to.tfvars",
		"-var", "foo=bar",
		"-json",
	}, nil, modulesCmd)
}
