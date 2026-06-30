// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestOutputCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		var config = outputConfig{}
		outputCmd := tf.outputCmd(context.Background(), config)

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
		}, nil, outputCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		var config = outputConfig{}
		config.state = "teststate"
		outputCmd := tf.outputCmd(context.Background(), config)

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
			"-state=teststate",
		}, nil, outputCmd)
	})

	t.Run("defaults with single output", func(t *testing.T) {
		var config = outputConfig{}
		config.state = "teststate"
		config.name = "testoutput"
		outputCmd := tf.outputCmd(context.Background(), config)

		assertCmd(t, []string{
			"output",
			"-no-color",
			"-json",
			"-state=teststate",
			"testoutput",
		}, nil, outputCmd)
	})
}
