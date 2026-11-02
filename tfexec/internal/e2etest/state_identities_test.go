// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateIdentities_tooOldVersion(t *testing.T) {
	versions := []string{testutil.Latest_v1_12}

	runTestWithVersions(t, versions, "state_identities", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		re := regexp.MustCompile("terraform state identities -json was added in 1.18.0")

		_, err := tf.StateIdentities(context.Background())
		if err == nil {
			t.Fatal("expected StateIdentities to fail in old version")
		}
		if !re.MatchString(err.Error()) {
			t.Fatalf("error running StateIdentities: %s", err)
		}
	})
}

func TestStateIdentities_basic(t *testing.T) {
	runTest(t, "state_identities", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(stateIdentitiesMinVersion) {
			t.Skip("state identities -json command not available in this Terraform version")
		}

		output, err := tf.StateIdentities(context.Background())
		if err != nil {
			t.Fatalf("error running StateIdentities: %s", err)
		}

		expected := &tfjson.StateIdentitiesOutput{
			FormatVersion: "1.0",
			Identities: map[string]map[string]any{
				"concept_pet.name": {
					"id":   "choice-monkey",
					"legs": json.Number("4"),
				},
			},
		}

		if diff := cmp.Diff(expected, output); diff != "" {
			t.Fatalf("output mismatch: %s", diff)
		}
	})
}
