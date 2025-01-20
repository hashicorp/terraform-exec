// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestApplyCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("basic", func(t *testing.T) {
		applyCmd, err := tf.applyCmd(context.Background(),
			Backup("testbackup"),
			LockTimeout("200s"),
			State("teststate"),
			StateOut("teststateout"),
			VarFile("foo.tfvars"),
			VarFile("bar.tfvars"),
			Lock(false),
			Parallelism(99),
			Refresh(false),
			Replace("aws_instance.test"),
			Replace("google_pubsub_topic.test"),
			Target("target1"),
			Target("target2"),
			Var("var1=foo"),
			Var("var2=bar"),
			Destroy(true),
			DirOrPlan("testfile"),
		)
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=foo.tfvars",
			"-var-file=bar.tfvars",
			"-lock=false",
			"-parallelism=99",
			"-refresh=false",
			"-replace=aws_instance.test",
			"-replace=google_pubsub_topic.test",
			"-destroy",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"testfile",
		}, nil, applyCmd)
	})

	t.Run("refresh-only operation", func(t *testing.T) {
		applyCmd, err := tf.applyCmd(context.Background(),
			RefreshOnly(true),
		)
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
			"-refresh-only",
		}, nil, applyCmd)
	})
}

func TestApplyJSONCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("basic", func(t *testing.T) {
		applyCmd, err := tf.applyJSONCmd(context.Background(),
			Backup("testbackup"),
			LockTimeout("200s"),
			State("teststate"),
			StateOut("teststateout"),
			VarFile("foo.tfvars"),
			VarFile("bar.tfvars"),
			Lock(false),
			Parallelism(99),
			Refresh(false),
			Replace("aws_instance.test"),
			Replace("google_pubsub_topic.test"),
			Target("target1"),
			Target("target2"),
			Var("var1=foo"),
			Var("var2=bar"),
			DirOrPlan("testfile"),
		)
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=foo.tfvars",
			"-var-file=bar.tfvars",
			"-lock=false",
			"-parallelism=99",
			"-refresh=false",
			"-replace=aws_instance.test",
			"-replace=google_pubsub_topic.test",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"-json",
			"testfile",
		}, nil, applyCmd)
	})
}

func TestApplyCmd_AllowDeferral(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_Alpha_v1_9))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("allow deferrals during apply", func(t *testing.T) {
		applyCmd, err := tf.applyCmd(context.Background(),
			AllowDeferral(true),
		)
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
			"-allow-deferral",
		}, nil, applyCmd)
	})
}
