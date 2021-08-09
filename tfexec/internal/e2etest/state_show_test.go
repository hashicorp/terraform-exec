package e2etest

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestStateShow(t *testing.T) {
	runTest(t, "basic_with_state", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(providerAddressMinVersion) {
			t.Skip("state file provider FQNs not compatible with this Terraform version")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		resources := []struct {
			address               string
			expectedOutputPattern string
			expectedErr           bool
		}{
			{
				"null_resource.foo",
				"^# null_resource\\.foo:\nresource \"null_resource\" \"foo\" {",
				false,
			},
			{
				"null_resource.bar",
				// failure results in empty output
				"^$",
				true,
			},
		}

		for _, resource := range resources {
			gotOutput, err := tf.StateShow(context.Background(), resource.address)
			gotErr := err != nil

			if gotErr && !resource.expectedErr {
				t.Fatalf("unexpected error running StateShow: %s", err)
			}

			if gotErr != resource.expectedErr {
				t.Errorf("terraform state show %s error: %s", resource.address, err)
			}

			expectedRegexp, err := regexp.Compile(resource.expectedOutputPattern)
			if err != nil {
				t.Fatalf("unable to compile regexp: %s", err)
			}

			if !expectedRegexp.Match([]byte(gotOutput)) {
				t.Errorf("terraform state show %s = %s", resource.address, gotOutput)
			}
		}
	})
}
