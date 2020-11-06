package e2etest

import (
	"context"
	"errors"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/andybalholm/crlf"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

var (
	showMinVersion = version.Must(version.NewVersion("0.12.0"))

	providerAddressMinVersion = version.Must(version.NewVersion("0.13.0"))
)

func TestShow(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		providerName := "registry.terraform.io/-/null"
		if tfv.LessThan(providerAddressMinVersion) {
			providerName = "null"
		}

		expected := &tfjson.State{
			FormatVersion: "0.1",
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "5510719323588825107",
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: providerName,
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

		if diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(tfjson.State{}, "TerraformVersion")); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShow_errInitRequired(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		var noInit *tfexec.ErrNoInit
		_, err := tf.Show(context.Background())
		if !errors.As(err, &noInit) {
			t.Fatalf("expected error ErrNoInit, got %T: %s", err, err)
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
		if mismatch.Actual != "0.11.14" {
			t.Fatalf("expected version 0.11.14, got %q", mismatch.Actual)
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
	runTestVersions(t, []string{testutil.Latest012}, "non_default_statefile_012", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		expected := &tfjson.State{
			FormatVersion: "0.1",
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "3610244792381545397",
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

		if diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(tfjson.State{}, "TerraformVersion")); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowStateFile013(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest013}, "non_default_statefile_013", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		expected := &tfjson.State{
			FormatVersion: "0.1",
			// TerraformVersion is ignored to facilitate latest version testing
			Values: &tfjson.StateValues{
				RootModule: &tfjson.StateModule{
					Resources: []*tfjson.StateResource{{
						Address: "null_resource.foo",
						AttributeValues: map[string]interface{}{
							"id":       "3610244792381545397",
							"triggers": nil,
						},
						Mode:         tfjson.ManagedResourceMode,
						Type:         "null_resource",
						Name:         "foo",
						ProviderName: "registry.terraform.io/hashicorp/null",
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

		if diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(tfjson.State{}, "TerraformVersion")); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

// Plan files cannot be transferred between different Terraform versions,
// so we maintain one fixture per supported version
func TestShowPlanFile012_linux(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest012}, "non_default_planfile_012", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
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

		if diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(tfjson.Plan{}, "TerraformVersion")); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowPlanFile013(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest013}, "non_default_planfile_013", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		providerName := "registry.terraform.io/hashicorp/null"

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

		if diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(tfjson.Plan{}, "TerraformVersion")); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestShowPlanFileRaw012_linux(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest012}, "non_default_planfile_012", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if runtime.GOOS != "linux" {
			t.Skip("plan file created in 0.12 on Linux is not compatible with other systems")
		}

		// crlf will standardize our line endings for us
		f, err := crlf.Open("testdata/non_default_planfile_012/human_readable_output.txt")
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

		if strings.TrimSpace(actual) != strings.TrimSpace(string(expected)) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(strings.TrimSpace(actual), strings.TrimSpace(string(expected)), false)
			t.Fatalf("actual:\n\n%s\n\nexpected:\n\n%s\n\ndiff:\n\n%s", actual, string(expected), dmp.DiffPrettyText(diffs))
		}
	})
}

func TestShowPlanFileRaw013(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest013}, "non_default_planfile_013", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		// crlf will standardize our line endings for us
		f, err := crlf.Open("testdata/non_default_planfile_013/human_readable_output.txt")
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

		if strings.TrimSpace(actual) != strings.TrimSpace(string(expected)) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(strings.TrimSpace(actual), strings.TrimSpace(string(expected)), false)
			t.Fatalf("actual:\n\n%s\n\nexpected:\n\n%s\n\ndiff:\n\n%s", actual, string(expected), dmp.DiffPrettyText(diffs))
		}
	})
}
