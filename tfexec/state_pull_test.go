// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStatePull(t *testing.T) {
	tf, err := NewTerraform(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	tf.SetEnv(map[string]string{})

	t.Run("tfstate", func(t *testing.T) {
		statePullCmd := tf.statePullCmd(context.Background(), nil)

		assertCmd(t, []string{
			"state",
			"pull",
		}, nil, statePullCmd)
	})
}
