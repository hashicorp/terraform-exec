// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

var (
	v1_0_1 = version.Must(version.NewVersion("1.0.1"))
)

func TestShow(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {

		if tfv.LessThan(providerAddressMinVersion) {
			t.Skip("state file provider FQNs not compatible with this Terraform version")
		}

		providerName := "registry.opentofu.org/hashicorp/null"
		if tfv.LessThan(providerAddressMinVersion) {
			providerName = "null"
		}

		formatVersion := "0.1"
		var sensitiveValues json.RawMessage
		if tfv.Core().GreaterThanOrEqual(v1_0_1) {
			formatVersion = "0.2"
			sensitiveValues = json.RawMessage([]byte("{}"))
		}
		if tfv.Core().GreaterThanOrEqual(v1_1) {
			formatVersion = "1.0"
		}

		expected := &tfjson.State{
			FormatVersion: formatVersion,
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "5510719323588825107",
							"triggers": nil,
						},
						SensitiveValues: sensitiveValues,
						Mode:            tfjson.ManagedResourceMode,
						Type:            "null_resource",
						Name:            "foo",
						ProviderName:    providerName,
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_emptyDir(t *testing.T) {
	runTest(t, "empty", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		formatVersion := "0.1"
		if tfv.Core().GreaterThanOrEqual(v1_0_1) {
			formatVersion = "0.2"
		}
		if tfv.Core().GreaterThanOrEqual(v1_1) {
			formatVersion = "1.0"
		}

		expected := &tfjson.State{
			FormatVersion: formatVersion,
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_noInitBasic(t *testing.T) {
	// Prior to v1.2.0, running show before init always results in an error.
	// In the basic case, in which the local backend is implicit and there are
	// no providers to download, this is unintended behaviour, as
	// init is not actually necessary. This is considered a known issue in
	// pre-1.2.0 versions.
	runTestWithVersions(t, []string{testutil.Latest012, testutil.Latest013, testutil.Latest014, testutil.Latest015, testutil.Latest_v1, testutil.Latest_v1_1}, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})

	// From v1.2.0 onwards, running show before init in the basic case returns
	// an empty state with no error.
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		// HACK KEM: Really I mean tfv.LessThan(version.Must(version.NewVersion("1.2.0"))),
		// but I want this test to run for refs/heads/main prior to the release of v1.2.0.
		if tfv.LessThan(version.Must(version.NewVersion("1.2.0"))) {

			t.Skip("test applies only to v1.2.0 and greater")
		}
		expected := &tfjson.State{
			FormatVersion: "1.0",
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_noInitModule(t *testing.T) {
	// Prior to v1.2.0, running show before init always results in an error.
	// In the basic case, in which the local backend is implicit and there are
	// no providers to download, this is unintended behaviour, as
	// init is not actually necessary. This is considered a known issue in
	// pre-1.2.0 versions.
	runTestWithVersions(t, []string{testutil.Latest012, testutil.Latest013, testutil.Latest014, testutil.Latest015, testutil.Latest_v1, testutil.Latest_v1_1}, "registry_module", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})

	// From v1.2.0 onwards, running show before init in the basic case returns
	// an empty state with no error.
	runTest(t, "registry_module", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		// HACK KEM: Really I mean tfv.LessThan(version.Must(version.NewVersion("1.2.0"))),
		// but I want this test to run for refs/heads/main prior to the release of v1.2.0.
		if tfv.LessThanOrEqual(version.Must(version.NewVersion(testutil.Latest_v1_1))) {
			t.Skip("test applies only to v1.2.0 and greater")
		}
		expected := &tfjson.State{
			FormatVersion: "1.0",
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_noInitInmemBackend(t *testing.T) {
	runTest(t, "inmem_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitLocalBackendNonDefaultState(t *testing.T) {
	runTest(t, "local_backend_non_default_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitCloudBackend(t *testing.T) {
	runTest(t, "cloud_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(version.Must(version.NewVersion("1.1.0"))) {
			t.Skip("cloud backend was added in Terraform 1.1, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitEtcdBackend(t *testing.T) {
	runTest(t, "etcd_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		if tfv.GreaterThanOrEqual(version.Must(version.NewVersion("1.3.0"))) || tfv.Prerelease() != "" {
			t.Skip("etcd backend was removed in Terraform 1.3, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_noInitRemoteBackend(t *testing.T) {
	runTest(t, "remote_backend", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_statefileDoesNotExist(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		_, err = tf.ShowStateFile(context.Background(), "statefilefoo")
		if err == nil {
			t.Fatalf("expected error, but did not get one")
		}
	})
}

func TestShow_versionMismatch(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		// only testing versions without show
		if tfv.GreaterThanOrEqual(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		var mismatch *tfexec.ErrVersionMismatch
		_, err := tf.Show(context.Background())
		if !errors.As(err, &mismatch) {
			t.Fatalf("expected version mismatch error, got %T %s", err, err)
		}
		if mismatch.Actual != "0.11.15" {
			t.Fatalf("expected version 0.11.15, got %q", mismatch.Actual)
		}
		if mismatch.MinInclusive != "0.12.0" {
			t.Fatalf("expected min 0.12.0, got %q", mismatch.MinInclusive)
		}
		if mismatch.MaxExclusive != "-" {
			t.Fatalf("expected max -, got %q", mismatch.MaxExclusive)
		}
	})
}

// Non-default state files cannot be migrated from 0.12 to 0.13,
// so we maintain one fixture per supported version.
// See github.com/hashicorp/terraform/25920
func TestShowStateFile012(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest012}, "non_default_statefile_012", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		expected := &tfjson.State{
			FormatVersion: "0.1",
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "2363759301357831073",
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: "null",
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowStateFile(context.Background(), "statefilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowStateFile013(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest013, testutil.Latest014}, "non_default_statefile_013", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		expected := &tfjson.State{
			FormatVersion: "0.1",
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "6724959521006014491",
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: "registry.opentofu.org/hashicorp/null",
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowStateFile(context.Background(), "statefilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowStateFile014(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest014}, "non_default_statefile_014", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		expected := &tfjson.State{
			FormatVersion: "0.1",
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "3544690470898862261",
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: "registry.opentofu.org/hashicorp/null",
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowStateFile(context.Background(), "statefilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

// Plan files cannot be transferred between different Terraform versions,
// so we maintain one fixture per supported version
func TestShowPlanFile012_linux(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest012}, "non_default_planfile_012", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if runtime.GOOS != "linux" {
			t.Skip("plan file created in 0.12 on Linux is not compatible with other systems")
		}

		providerName := "null"

		expected := &tfjson.Plan{
			FormatVersion: "0.1",
			// TerraformVersion is ignored to facilitate latest version testing
			PlannedValues: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: providerName,
					}},
				},
			},
			ResourceChanges: []*tfjson.ResourceChange{{
				Address:      "null_resource.foo",
				Mode:         tfjson.ManagedResourceMode,
				Type:         "null_resource",
				Name:         "foo",
				ProviderName: providerName,
				Change: &tfjson.Change{
					Actions:      tfjson.Actions{tfjson.ActionCreate},
					After:        map[string]interface{}{"triggers": nil},
					AfterUnknown: map[string]interface{}{"id": (true)},
				},
			}},
			Config: &tfjson.Config{
				RootModule: &tfjson.ConfigModule{
					Resources: []*tfjson.ConfigResource{{
						Address:           "null_resource.foo",
						Mode:              tfjson.ManagedResourceMode,
						Type:              "null_resource",
						Name:              "foo",
						ProviderConfigKey: "null",
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowPlanFile(context.Background(), "planfilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffPlan(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowPlanFile013(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest013}, "non_default_planfile_013", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		providerName := "registry.opentofu.org/hashicorp/null"

		expected := &tfjson.Plan{
			// TerraformVersion is ignored to facilitate latsest version testing
			FormatVersion: "0.1",
			PlannedValues: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: providerName,
					}},
				},
			},
			ResourceChanges: []*tfjson.ResourceChange{{
				Address:      "null_resource.foo",
				Mode:         tfjson.ManagedResourceMode,
				Type:         "null_resource",
				Name:         "foo",
				ProviderName: providerName,
				Change: &tfjson.Change{
					Actions:      tfjson.Actions{tfjson.ActionCreate},
					After:        map[string]interface{}{"triggers": nil},
					AfterUnknown: map[string]interface{}{"id": true},
				},
			}},
			Config: &tfjson.Config{
				RootModule: &tfjson.ConfigModule{
					Resources: []*tfjson.ConfigResource{{
						Address:           "null_resource.foo",
						Mode:              tfjson.ManagedResourceMode,
						Type:              "null_resource",
						Name:              "foo",
						ProviderConfigKey: "null",
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowPlanFile(context.Background(), "planfilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffPlan(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowPlanFile014(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest014}, "non_default_planfile_014", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		providerName := "registry.opentofu.org/hashicorp/null"

		expected := &tfjson.Plan{
			// TerraformVersion is ignored to facilitate latsest version testing
			FormatVersion: "0.1",
			PlannedValues: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: providerName,
					}},
				},
			},
			ResourceChanges: []*tfjson.ResourceChange{{
				Address:      "null_resource.foo",
				Mode:         tfjson.ManagedResourceMode,
				Type:         "null_resource",
				Name:         "foo",
				ProviderName: providerName,
				Change: &tfjson.Change{
					Actions:      tfjson.Actions{tfjson.ActionCreate},
					After:        map[string]interface{}{"triggers": nil},
					AfterUnknown: map[string]interface{}{"id": true},
				},
			}},
			Config: &tfjson.Config{
				ProviderConfigs: map[string]*tfjson.ProviderConfig{
					"null": {Name: "null", VersionConstraint: "3.0.0"},
				},
				RootModule: &tfjson.ConfigModule{
					Resources: []*tfjson.ConfigResource{{
						Address:           "null_resource.foo",
						Mode:              tfjson.ManagedResourceMode,
						Type:              "null_resource",
						Name:              "foo",
						ProviderConfigKey: "null",
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowPlanFile(context.Background(), "planfilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffPlan(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowPlanFileRaw012_linux(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest012}, "non_default_planfile_012", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if runtime.GOOS != "linux" {
			t.Skip("plan file created in 0.12 on Linux is not compatible with other systems")
		}

		f, err := os.Open("testdata/non_default_planfile_012/human_readable_output.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		expected, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		err = tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowPlanFileRaw(context.Background(), "planfilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(normalizePlanOutput(actual), normalizePlanOutput(string(expected))); diff != "" {
			t.Fatalf("unexpected difference: %s", diff)
		}
	})
}

func TestShowPlanFileRaw013(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest013}, "non_default_planfile_013", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		f, err := os.Open("testdata/non_default_planfile_013/human_readable_output.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		expected, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		err = tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowPlanFileRaw(context.Background(), "planfilefoo")
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(normalizePlanOutput(actual), normalizePlanOutput(string(expected))); diff != "" {
			t.Fatalf("unexpected difference: %s", diff)
		}
	})
}

func TestShowPlanFileRaw014(t *testing.T) {
	runTestWithVersions(t, []string{testutil.Latest014}, "non_default_planfile_014", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		f, err := os.Open("testdata/non_default_planfile_013/human_readable_output.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		expected, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		err = tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		actual, err := tf.ShowPlanFileRaw(context.Background(), "planfilefoo")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(normalizePlanOutput(actual), normalizePlanOutput(string(expected))); diff != "" {
			t.Fatalf("unexpected difference: %s", diff)
		}
	})
}

func TestShowBigInt(t *testing.T) {
	runTest(t, "bigint", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		providerName := "registry.opentofu.org/hashicorp/random"
		if tfv.LessThan(providerAddressMinVersion) {
			providerName = "random"
		}

		formatVersion := "0.1"
		var sensitiveValues json.RawMessage

		if tfv.Core().GreaterThanOrEqual(v1_0_1) {
			formatVersion = "0.2"
			sensitiveValues = json.RawMessage([]byte("{}"))
		}
		if tfv.Core().GreaterThanOrEqual(v1_1) {
			formatVersion = "1.0"
		}

		expected := &tfjson.State{
			FormatVersion: formatVersion,
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "random_integer.bigint",
						AttributeValues: map[string]interface{}{
							"id":      "7227701560655103598",
							"max":     json.Number("7227701560655103598"),
							"min":     json.Number("7227701560655103597"),
							"result":  json.Number("7227701560655103598"),
							"seed":    "12345",
							"keepers": nil,
						},
						SensitiveValues: sensitiveValues,
						Mode:            tfjson.ManagedResourceMode,
						Type:            "random_integer",
						Name:            "bigint",
						ProviderName:    providerName,
					}},
				},
			},
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.Apply(context.Background())
		if err != nil {
			t.Fatalf("error running Apply in test directory: %s", err)
		}

		actual, err := tf.Show(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := diffState(expected, actual); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

// Since our plan strings are not large, prefer simple cross-platform
// normalization handling over pulling in a dependency.
func normalizePlanOutput(str string) string {
	// Ignore any extra newlines at the beginning or end of output
	str = strings.TrimSpace(str)
	// Normalize CRLF to LF for cross-platform testing
	str = strings.Replace(str, "\r\n", "\n", -1)

	return str
}
