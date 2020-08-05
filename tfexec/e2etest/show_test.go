package e2etest

import (
	"context"
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
