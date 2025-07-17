// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"bytes"
	"context"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestQueryJSON_TF112(t *testing.T) {
	versions := []string{testutil.Latest_v1_12}

	runTestWithVersions(t, versions, "query", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		re := regexp.MustCompile("terraform query -json was added in 1.14.0")

		err = tf.QueryJSON(context.Background(), io.Discard)
		if err != nil && !re.MatchString(err.Error()) {
			t.Fatalf("error running Query: %s", err)
		}
	})
}

func TestQueryJSON_TF114(t *testing.T) {
	versions := []string{testutil.Latest_Alpha_v1_14}

	runTestWithVersions(t, versions, "query", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		var output bytes.Buffer
		err = tf.QueryJSON(context.Background(), &output)
		if err != nil {
			t.Fatalf("error running Query: %s", err)
		}

		results := strings.Count(output.String(), "list.concept_pet.pets: Result found")
		if results != 5 {
			t.Fatalf("expected 5 query results, but got %d", results)
		}
	})
}
