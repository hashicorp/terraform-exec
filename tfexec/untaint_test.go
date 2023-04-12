// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestUntaintCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		untaintCmd := tf.untaintCmd(context.Background(), "aws_instance.foo")

		assertCmd(t, []string{
			"untaint",
			"-no-color",
			"-lock=true",
			"aws_instance.foo",
		}, nil, untaintCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		untaintCmd := tf.untaintCmd(context.Background(), "aws_instance.foo",
			State("teststate"),
			AllowMissing(true),
			LockTimeout("200s"),
			Lock(false))

		assertCmd(t, []string{
			"untaint",
			"-no-color",
			"-lock-timeout=200s",
			"-state=teststate",
			"-lock=false",
			"-allow-missing",
			"aws_instance.foo",
		}, nil, untaintCmd)
	})
}
