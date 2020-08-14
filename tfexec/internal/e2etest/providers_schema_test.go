package e2etest

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	providersSchemaJSONMinVersion = version.Must(version.NewVersion("0.12.0"))
)

func TestProvidersSchema(t *testing.T) {
	for i, c := range []struct {
		fixtureDir string
		expected   func(*version.Version) *tfjson.ProviderSchemas
	}{
		{
			"basic", func(tfv *version.Version) *tfjson.ProviderSchemas {
				dk := tfjson.SchemaDescriptionKindPlain
				// HACK: this is a bug? the contstant value is "plaintext" but "plain" is output,
				// not sure if its in tfjson or what
				dk = "plain"

				providerName := "registry.terraform.io/hashicorp/null"

				if tfv.LessThan(version.Must(version.NewVersion("0.13.0"))) {
					providerName = "null"
					dk = ""
				}

				return &tfjson.ProviderSchemas{
					FormatVersion: "0.1",
					Schemas: map[string]*tfjson.ProviderSchema{
						providerName: {
							ConfigSchema: &tfjson.Schema{
								Version: 0,
								Block: &tfjson.SchemaBlock{
									DescriptionKind: dk,
								},
							},
							ResourceSchemas: map[string]*tfjson.Schema{
								"null_resource": {
									Version: 0,
									Block: &tfjson.SchemaBlock{
										DescriptionKind: dk,
										Attributes: map[string]*tfjson.SchemaAttribute{
											"id": {
												AttributeType:   cty.String,
												Optional:        true,
												Computed:        true,
												DescriptionKind: dk,
											},
											"triggers": {
												AttributeType:   cty.Map(cty.String),
												Optional:        true,
												DescriptionKind: dk,
											},
										},
									},
								},
							},
							DataSourceSchemas: map[string]*tfjson.Schema{
								"null_data_source": {
									Version: 0,
									Block: &tfjson.SchemaBlock{
										DescriptionKind: dk,
										Attributes: map[string]*tfjson.SchemaAttribute{
											"has_computed_default": {
												AttributeType:   cty.String,
												Optional:        true,
												Computed:        true,
												DescriptionKind: dk,
											},
											"id": {
												AttributeType:   cty.String,
												Optional:        true,
												Computed:        true,
												DescriptionKind: dk,
											},
											"inputs": {
												AttributeType:   cty.Map(cty.String),
												Optional:        true,
												DescriptionKind: dk,
											},
											"outputs": {
												AttributeType:   cty.Map(cty.String),
												Computed:        true,
												DescriptionKind: dk,
											},
											"random": {
												AttributeType:   cty.String,
												Computed:        true,
												DescriptionKind: dk,
											},
										},
									},
								},
							},
						}},
				}
			},
		},
		{
			"empty_with_tf_file", func(*version.Version) *tfjson.ProviderSchemas {
				return &tfjson.ProviderSchemas{
					FormatVersion: "0.1",
					Schemas:       nil,
				}
			},
		},
	} {
		c := c
		t.Run(fmt.Sprintf("%d %s", i, c.fixtureDir), func(t *testing.T) {
			runTest(t, c.fixtureDir, func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
				if tfv.LessThan(providersSchemaJSONMinVersion) {
					t.Skip("providers schema -json was added in 0.12")
				}

				err := tf.Init(context.Background())
				if err != nil {
					t.Fatalf("error running Init in test directory: %s", err)
				}

				schemas, err := tf.ProvidersSchema(context.Background())
				if err != nil {
					t.Fatalf("error running ProvidersSchema in test directory: %s", err)
				}

				expected := c.expected(tfv)
				if !reflect.DeepEqual(schemas, expected) {
					t.Fatalf("expected %+v, but got %+v", spew.Sdump(expected), spew.Sdump(schemas))
				}
			})
		})
	}

}

func TestProvidersSchema_versionMismatch(t *testing.T) {
	t.Skip("TODO! add version mismatch test for 0.11 as -json was added in 0.12 (I think)")
}
