package e2etest

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestShow(t *testing.T) {
	tf, cleanup := setupFixture(t, testutil.Latest012, "basic_with_state")
	defer cleanup()

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
					ProviderName: "null",
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
}

func TestShow_errInitRequired(t *testing.T) {
	tf, cleanup := setupFixture(t, testutil.Latest012, "basic")
	defer cleanup()

	_, err := tf.Show(context.Background())
	if err == nil {
		t.Fatal("expected Show to error, but it did not")
	}
	if _, ok := err.(*tfexec.ErrNoInit); !ok {
		t.Fatalf("expected error %s to be ErrNoInit", err)
	}
}

func TestShow_compatible(t *testing.T) {
	tf, cleanup := setupFixture(t, testutil.Latest011, "basic")
	defer cleanup()

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
}
