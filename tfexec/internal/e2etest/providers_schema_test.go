package e2etest

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func TestProvidersSchema(t *testing.T) {
	for i, c := range []struct {
		fixtureDir string
		expected   *tfjson.ProviderSchemas
	}{
		{
			"basic", &tfjson.ProviderSchemas{
				FormatVersion: "0.1",
				Schemas: map[string]*tfjson.ProviderSchema{
					"null": {
						ConfigSchema: &tfjson.Schema{
							Version: 0,
							Block:   &tfjson.SchemaBlock{},
						},
						ResourceSchemas: map[string]*tfjson.Schema{
							"null_resource": {
								Version: 0,
								Block: &tfjson.SchemaBlock{
									Attributes: map[string]*tfjson.SchemaAttribute{
										"id": {
											AttributeType: cty.String,
											Optional:      true,
											Computed:      true,
										},
										"triggers": {
											AttributeType: cty.Map(cty.String),
											Optional:      true,
										},
									},
								},
							},
						},
						DataSourceSchemas: map[string]*tfjson.Schema{
							"null_data_source": {
								Version: 0,
								Block: &tfjson.SchemaBlock{
									Attributes: map[string]*tfjson.SchemaAttribute{
										"has_computed_default": {
											AttributeType: cty.String,
											Optional:      true,
											Computed:      true,
										},
										"id": {
											AttributeType: cty.String,
											Optional:      true,
											Computed:      true,
										},
										"inputs": {
											AttributeType: cty.Map(cty.String),
											Optional:      true,
										},
										"outputs": {
											AttributeType: cty.Map(cty.String),
											Computed:      true,
										},
										"random": {
											AttributeType: cty.String,
											Computed:      true,
										},
									},
								},
							},
						},
					}},
			},
		},
		{
			"empty_with_tf_file", &tfjson.ProviderSchemas{
				FormatVersion: "0.1",
				Schemas:       nil,
			},
		},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.fixtureDir), func(t *testing.T) {
			runTest(t, []string{
				testutil.Latest012,
			}, c.fixtureDir, func(t *testing.T, tfv string, tf *tfexec.Terraform) {
				err := tf.Init(context.Background())
				if err != nil {
					t.Fatalf("error running Init in test directory: %s", err)
				}

				schemas, err := tf.ProvidersSchema(context.Background())
				if err != nil {
					t.Fatalf("error running ProvidersSchema in test directory: %s", err)
				}

				if !reflect.DeepEqual(schemas, c.expected) {
					t.Fatalf("expected %+v, but got %+v", spew.Sdump(c.expected), spew.Sdump(schemas))
				}
			})
		})
	}

}
