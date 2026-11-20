// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestDestroyCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		destroyCmd, err := tf.destroyCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"destroy",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
		}, nil, destroyCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		destroyCmd, err := tf.destroyCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), Dir("destroydir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"destroy",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-parallelism=99",
			"-refresh=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"--",
			"destroydir",
		}, nil, destroyCmd)
	})
}

func TestDestroyJSONCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		destroyCmd, err := tf.destroyJSONCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"destroy",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
			"-json",
		}, nil, destroyCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		destroyCmd, err := tf.destroyJSONCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), Dir("destroydir"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"destroy",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-parallelism=99",
			"-refresh=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"-json",
			"--",
			"destroydir",
		}, nil, destroyCmd)
	})
}

func TestDestroyCmd_MinimalRefresh(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_17))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("minimal-refresh destroy", func(t *testing.T) {
		destroyCmd, err := tf.destroyCmd(context.Background(), MinimalRefresh(true))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"destroy",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock-timeout=0s",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
			"-minimal-refresh",
		}, nil, destroyCmd)
	})

	t.Run("minimal-refresh unsupported version", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_16))
		if err != nil {
			t.Fatal(err)
		}

		_, err = tf.destroyCmd(context.Background(), MinimalRefresh(true))
		if err == nil {
			t.Fatal("expected an error but received none")
		}

		expectedErr := "minimal-refresh option was introduced in Terraform 1.17.0"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("expected error to contain %q, got: %q", expectedErr, err)
		}
	})

	t.Run("minimal-refresh and refresh=false returns an error", func(t *testing.T) {
		_, err := tf.destroyCmd(context.Background(), Refresh(false), MinimalRefresh(true))
		if err == nil {
			t.Fatal("expected an error but received none")
		}

		expectedErr := "cannot use refresh=false with minimal-refresh"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("expected error to contain %q, got: %q", expectedErr, err)
		}
	})
}
