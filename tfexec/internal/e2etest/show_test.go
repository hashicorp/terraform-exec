package e2etest

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-exec/tfexec"
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

		expected := tfjson.State{
			FormatVersion: "0.1",
			// this is the version that wrote state, not the version that is running
			TerraformVersion: "0.12.24",
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

		if !reflect.DeepEqual(actual, &expected) {
			t.Fatalf("actual: %s\nexpected: %s", spew.Sdump(actual), spew.Sdump(expected))
		}
	})
}

func TestShow_errInitRequired(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(showMinVersion) {
			t.Skip("terraform show was added in Terraform 0.12, so test is not valid")
		}

		_, err := tf.Show(context.Background())
		if err == nil {
			t.Fatal("expected Show to error, but it did not")
		}
		if _, ok := err.(*tfexec.ErrNoInit); !ok {
			t.Fatalf("expected error %s to be ErrNoInit", err)
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
			t.Fatal("expected version mismatch error")
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
