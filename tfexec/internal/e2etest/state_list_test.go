package e2etest

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestStateList(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(providerAddressMinVersion) {
			t.Skip("state file provider FQNs not compatible with this Terraform version")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		gotAddresses, err := tf.StateList(context.Background())

		if err != nil {
			t.Fatalf("error running StateList: %s", err)
		}

		expectedAddresses := []string{"null_resource.foo"}

		if !reflect.DeepEqual(gotAddresses, expectedAddresses) {
			t.Errorf("terraform state list = %v, expected %v", gotAddresses, expectedAddresses)
		}
	})
}
