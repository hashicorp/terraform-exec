// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestHideWindow(t *testing.T) {
	tf, err := NewTerraform(t.TempDir(), "terraform.exe")
	if err != nil {
		t.Fatal(err)
	}
	build := func() bool {
		cmd := tf.buildTerraformCmd(context.Background(), nil, "version")
		return cmd.SysProcAttr != nil && cmd.SysProcAttr.HideWindow
	}
	if build() {
		t.Fatal("command windows should not be hidden by default")
	}
	tf.SetHideWindow(true)
	if !build() {
		t.Fatal("SetHideWindow(true) did not hide the command window")
	}
	tf.SetHideWindow(false)
	if build() {
		t.Fatal("SetHideWindow(false) did not restore the default")
	}
}

func TestHiddenCommandRuns(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	tf, err := NewTerraform(t.TempDir(), executable)
	if err != nil {
		t.Fatal(err)
	}
	tf.SetHideWindow(true)
	cmd := tf.buildTerraformCmd(context.Background(), map[string]string{"MOCK_SLEEP_DURATION": "1ns"})
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("hidden command failed: %s\n%s", err, output)
	}
}

func TestHideWindowVersion(t *testing.T) {
	tf, err := NewTerraform(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}
	tf.SetHideWindow(true)
	version, _, err := tf.Version(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
	if version.String() != testutil.Latest_v1 {
		t.Fatalf("expected Terraform %s, got %s", testutil.Latest_v1, version)
	}
}
