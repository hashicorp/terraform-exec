// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateRmCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		stateRmCmd, err := tf.stateRmCmd(context.Background(), "testAddress")
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"rm",
			"-no-color",
			"-lock-timeout=0s",
			"-lock=true",
			"testAddress",
		}, nil, stateRmCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		stateRmCmd, err := tf.stateRmCmd(context.Background(), "testAddress", Backup("testbackup"), BackupOut("testbackupout"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), Lock(false))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"state",
			"rm",
			"-no-color",
			"-backup=testbackup",
			"-backup-out=testbackupout",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-lock=false",
			"testAddress",
		}, nil, stateRmCmd)
	})
}
