// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestQueryJSON_TF112(t *testing.T) {
	versions := []string{testutil.Latest_v1_12}

	runTestWithVersions(t, versions, "query", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		re := regexp.MustCompile("terraform query -json was added in 1.14.0")

		_, err = tf.QueryJSON(context.Background())
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

		le, err := tf.QueryJSON(context.Background())
		if err != nil {
			t.Fatalf("error running Query: %s", err)
		}

		results := 0
		var completeData tfjson.ListCompleteData
		for {
			msg, ok, err := le.NextMessage()
			if !ok {
				break
			}
			if err != nil {
				t.Fatalf("error getting next message: %s", err)
			}
			switch m := msg.(type) {
			case tfjson.ListResourceFoundMessage:
				results++
			case tfjson.ListCompleteMessage:
				completeData = m.ListComplete
			}
		}

		if results != 5 {
			t.Fatalf("expected 5 query results, but got %d", results)
		}
		expectedData := tfjson.ListCompleteData{
			Address:      "list.concept_pet.pets",
			ResourceType: "concept_pet",
			Total:        5,
		}
		if diff := cmp.Diff(expectedData, completeData); diff != "" {
			t.Fatalf("unexpected complete message data: %s", diff)
		}
	})
}
