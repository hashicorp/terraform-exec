// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func TestStateMigrateJSON_basic_migrationBetweenBackends(t *testing.T) {
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
				t.Logf("migration complete: %q", m.Msg)
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
}

func TestStateMigrateJSON_basic_migrationStateStoreConfigChange(t *testing.T) {
	// Developing the test locally requires having a filesystem mirror with 2 versions of the provider available,
	// and using a CLI configuration file that makes the provider be sourced from that mirror.
	//
	// Test fixture specifies v9.0.0
	usedVersion := "9.0.0"

	runTest(t, "state_migrate_state_store_config_change", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(stateMigrateMinVersion) {
			t.Skip("state migrate command not available in this Terraform version")
		}

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
				if m.Severity == tfjson.DiagnosticSeverityError {
					t.Errorf("unexpected error diagnostic: %q", m.Msg)
				}

				// Warnings expected when using fs mirror during development of this test.
				t.Logf("warning diagnostic message: %q", m.Msg)
			case tfjson.MigrationSourceInitializationStartMessage:
				t.Logf("migration source initialization started: %q", m.Msg)
			case tfjson.MigrationSourceInitializationCompleteMessage:
				t.Logf("migration source initialization complete: %q", m.Msg)
			case tfjson.MigrationDestinationInitializationStartMessage:
				t.Logf("migration destination initialization started: %q", m.Msg)
			case tfjson.MigrationDestinationInitializationCompleteMessage:
				t.Logf("migration destination initialization complete: %q", m.Msg)
			case tfjson.ProviderInstallationStartMessage:
				t.Logf("provider installation started: %q", m.Msg)
			case tfjson.ProviderQueryUseLatestMessage:
				t.Logf("provider query use latest: %q", m.Msg)
			case tfjson.ProviderQueryUsePreviousConstraintsMessage:
				// The test fixture doesn't have constraints; latest version should be used instead.
				t.Errorf("unexpected attempt to use version constraints to select provider version: %q", m.Msg)
			case tfjson.ProviderVersionInstallationStartMessage:
				t.Logf("provider version installation started: %q", m.Msg)
			case tfjson.ProviderVersionInstallationCompleteMessage:
				t.Logf("provider version installation complete: %q", m.Msg)
			case tfjson.MigrationStartMessage:
				t.Logf("migration started: %q", m.Msg)
			case tfjson.MigrationCompleteMessage:
				t.Logf("migration complete: %q", m.Msg)
			case tfjson.MigrationFinalizedMessage:
				finalMsg = m.Msg
			}
		}

		expectedMsg := fmt.Sprintf(`Finished migrating state from state store "tfcoremock_fs" (hashicorp/tfcoremock %s) to state store "tfcoremock_fs" (hashicorp/tfcoremock %s).`, usedVersion, usedVersion)
		if finalMsg != expectedMsg {
			t.Fatalf("unexpected final message: %s", finalMsg)
		}

		// The migrated state should be in the new workspace directory specified in the Terraform configuration.
		if _, err = os.Stat(filepath.Join(tf.WorkingDir(), "new-location", "default", "terraform.tfstate")); err != nil {
			if os.IsNotExist(err) {
				t.Fatalf("destination state does not exist: %s", err)
			}
			t.Fatalf("unexpectedly failed to stat destination state: %s", err)
		}

		// We don't compare the contents of the state, because Terraform creates new lineage identifiers for the migrated state;
		// it won't be an exact match of the source state.
	})
}

func TestStateMigrateJSON_upgradeProvider(t *testing.T) {
	// Developing the test locally requires having a filesystem mirror with 2 versions of the provider available,
	// and using a CLI configuration file that makes the provider be sourced from that mirror.
	//
	// Mirror set up to contain v9.0.0 and v9.1.0
	oldVersion := "9.0.0"
	newVersion := "9.1.0"

	runTest(t, "state_migrate_upgrade_state_store_provider", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(stateMigrateMinVersion) {
			t.Skip("state migrate command not available in this Terraform version")
		}

		upgradeOpt := tfexec.Upgrade(true)
		iter, err := tf.StateMigrateJSON(context.Background(), upgradeOpt)
		if err != nil {
			t.Fatalf("error running StateMigrateJSON: %s", err)
		}

		finalMsg := ""

		installCounter := 0
		sourceProviderInstallMsg := ""
		destinationProviderInstallMsg := ""
		for nextMsg := range iter {
			if nextMsg.Err != nil {
				t.Fatalf("error getting next message: %s", nextMsg.Err)
			}
			switch m := nextMsg.Msg.(type) {
			case tfjson.VersionLogMessage:
				t.Logf("versions: Terraform %q, UI %q", m.Terraform, m.UI)
			case tfjson.DiagnosticLogMessage:
				if m.Severity == tfjson.DiagnosticSeverityError {
					t.Errorf("unexpected error diagnostic: %q", m.Msg)
				}

				// Warnings expected when using fs mirror during development of this test.
				t.Logf("warning diagnostic message: %q", m.Msg)
			case tfjson.MigrationSourceInitializationStartMessage:
				t.Logf("migration source initialization started: %q", m.Msg)
			case tfjson.MigrationSourceInitializationCompleteMessage:
				t.Logf("migration source initialization complete: %q", m.Msg)
			case tfjson.MigrationDestinationInitializationStartMessage:
				t.Logf("migration destination initialization started: %q", m.Msg)
			case tfjson.MigrationDestinationInitializationCompleteMessage:
				t.Logf("migration destination initialization complete: %q", m.Msg)
			case tfjson.ProviderInstallationStartMessage:
				t.Logf("provider installation started: %q", m.Msg)
			case tfjson.ProviderQueryUsePreviousConstraintsMessage:
				// Looking for a newer provider version based on previous constraints
				t.Logf("provider query use previous constraints: %q", m.Msg)
			case tfjson.ProviderQueryUseLatestMessage:
				// Using the latest provider version is not expected in this test; constraints in the
				// test fixture should be used instead.
				t.Errorf("unexpected attempt to use latest provider version: %q", m.Msg)
			case tfjson.ProviderVersionInstallationStartMessage:
				t.Logf("provider version installation started: %q", m.Msg)
			case tfjson.ProviderVersionInstallationCompleteMessage:
				// The old and new provider versions will be installed, in that order.
				t.Logf("provider installation complete: %q", m.Msg)
				if installCounter == 0 {
					sourceProviderInstallMsg = m.Msg
				} else {
					destinationProviderInstallMsg = m.Msg
				}
				installCounter++
			case tfjson.MigrationStartMessage:
				t.Logf("migration started: %q", m.Msg)
			case tfjson.MigrationCompleteMessage:
				t.Logf("migration complete: %q", m.Msg)
			case tfjson.MigrationFinalizedMessage:
				finalMsg = m.Msg
			}
		}

		// Assert we see expected output describing installing the 2 versions of the same provider
		if sourceProviderInstallMsg == "" || !strings.Contains(sourceProviderInstallMsg, oldVersion) {
			t.Fatalf("expected provider installation message to contain old version %q, got: %q", oldVersion, sourceProviderInstallMsg)
		}
		if destinationProviderInstallMsg == "" || !strings.Contains(destinationProviderInstallMsg, newVersion) {
			t.Fatalf("expected provider installation message to contain new version %q, got: %q", newVersion, destinationProviderInstallMsg)
		}

		expectedFinalMsg := fmt.Sprintf(`Finished migrating state from state store "tfcoremock_fs" (hashicorp/tfcoremock %s) to state store "tfcoremock_fs" (hashicorp/tfcoremock %s).`, oldVersion, newVersion)
		if diff := cmp.Diff(expectedFinalMsg, finalMsg); diff != "" {
			t.Fatalf("unexpected final message: %s", diff)
		}

		// The migrated state should be in the new workspace directory specified in the Terraform configuration.
		if _, err = os.Stat(filepath.Join(tf.WorkingDir(), "new-location", "default", "terraform.tfstate")); err != nil {
			if os.IsNotExist(err) {
				t.Fatalf("destination state does not exist: %s", err)
			}
			t.Fatalf("unexpectedly failed to stat destination state: %s", err)
		}

		// We don't compare the contents of the state, because Terraform creates new lineage identifiers for the migrated state;
		// it won't be an exact match of the source state.
	})
}
