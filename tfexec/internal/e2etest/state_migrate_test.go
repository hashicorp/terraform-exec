// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestStateMigrateJSON_tooOldVersion(t *testing.T) {
	versions := []string{testutil.Latest_v1_12}

	runTestWithVersions(t, versions, "state_migrate_local_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		re := regexp.MustCompile("terraform state migrate -json was added in 1.18.0")

		_, err := tf.StateMigrateJSON(context.Background())
		if err == nil {
			t.Fatal("expected StateMigrateJSON to fail in old version")
		}
		if !re.MatchString(err.Error()) {
			t.Fatalf("error running StateMigrateJSON: %s", err)
		}
	})
}

func TestStateMigrateJSON_basic(t *testing.T) {
	runTest(t, "state_migrate_local_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(stateMigrateMinVersion) {
			t.Skip("state migrate command not available in this Terraform version")
		}

		srcBytes, err := os.ReadFile(filepath.Join(tf.WorkingDir(), "src.tfstate"))
		if err != nil {
			t.Fatalf("failed to read source state: %s", err)
		}
		sourceState := string(srcBytes)

		iter, err := tf.StateMigrateJSON(context.Background())
		if err != nil {
			t.Fatalf("error running StateMigrateJSON: %s", err)
		}

		finalMsg := ""
		for nextMsg := range iter {
			if nextMsg.Err != nil {
				t.Fatalf("error getting next message: %s", nextMsg.Err)
			}
			switch m := nextMsg.Msg.(type) {
			case tfjson.VersionLogMessage:
				t.Logf("versions: Terraform %q, UI %q", m.Terraform, m.UI)
			case tfjson.DiagnosticLogMessage:
				t.Errorf("unexpected diagnostic: %q", m.Msg)
			case tfjson.MigrationCompleteMessage:
				t.Logf("MSG complete: %q", m.Msg)
			case tfjson.MigrationFinalizedMessage:
				finalMsg = m.Msg
			}
		}

		expectedMsg := `Finished migrating state from backend "local" to backend "local".`
		if diff := cmp.Diff(expectedMsg, finalMsg); diff != "" {
			t.Fatalf("unexpected final message: %s", diff)
		}

		dstBytes, err := os.ReadFile(filepath.Join(tf.WorkingDir(), "dst.tfstate"))
		if err != nil {
			t.Fatalf("failed to read destination state: %s", err)
		}
		destinationState := string(dstBytes)

		if diff := cmp.Diff(sourceState, destinationState); diff != "" {
			t.Fatalf("unexpected difference after migration: %s", diff)
		}
	})

	// TODO: basic state store test (requires an implemented provider)

	// TODO: provider upgrade test (requires an implemented provider)
}
