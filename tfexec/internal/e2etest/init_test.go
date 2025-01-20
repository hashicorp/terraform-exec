// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"io"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestInit(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}
	})
}

func TestInitJSON_TF18AndEarlier(t *testing.T) {
	versions := []string{
		testutil.Latest011,
		testutil.Latest012,
		testutil.Latest013,
		testutil.Latest_v1_6,
		testutil.Latest_v1_7,
		testutil.Latest_v1_8,
	}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		re := regexp.MustCompile("terraform init -json was added in 1.9.0")

		err = tf.InitJSON(context.Background(), io.Discard)
		if err != nil && !re.MatchString(err.Error()) {
			t.Fatalf("error running Init: %s", err)
		}
	})
}

func TestInitJSON_TF19AndLater(t *testing.T) {
	versions := []string{
		testutil.Latest_v1_9,
		testutil.Latest_Alpha_v1_9,
		testutil.Latest_Alpha_v1_10,
	}

	runTestWithVersions(t, versions, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.InitJSON(context.Background(), io.Discard)
		if err != nil {
			t.Fatalf("error running Init: %s", err)
		}
	})
}
