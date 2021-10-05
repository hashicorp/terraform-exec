package e2etest

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	providersSchemaJSONMinVersion = version.Must(version.NewVersion("0.12.0"))
	v0_13_0                       = version.Must(version.NewVersion("0.13.0"))
	v0_15_0                       = version.Must(version.NewVersion("0.15.0"))
	v1_1                          = version.Must(version.NewVersion("1.1.0"))
)

func TestProvidersSchema(t *testing.T) {
	for i, c := range []struct {
		fixtureDir string
		expected   func(*version.Version) *tfjson.ProviderSchemas
	}{
		{
			"basic", func(tfv *version.Version) *tfjson.ProviderSchemas {
				var providerSchema *tfjson.ProviderSchemas

				// TODO: Add handling for v1 format once it lands in core
				// See https://github.com/hashicorp/terraform/pull/29550
				if tfv.Core().GreaterThanOrEqual(v0_15_0) {
					providerSchema = &tfjson.ProviderSchemas{
						FormatVersion: "0.2",
						Schemas: map[string]*tfjson.ProviderSchema{
							"registry.terraform.io/hashicorp/null": {
								ConfigSchema: &tfjson.Schema{
									Version: 0,
									Block: &tfjson.SchemaBlock{
										DescriptionKind: "plain",
									},
								},
								ResourceSchemas: map[string]*tfjson.Schema{
									"null_resource": {
										Version: 0,
										Block: &tfjson.SchemaBlock{
											DescriptionKind: "markdown",
											Description:     "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
											Attributes: map[string]*tfjson.SchemaAttribute{
												"id": {
													AttributeType:   cty.String,
													Optional:        false,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "This is set to a random value at create time.",
												},
												"triggers": {
													AttributeType:   cty.Map(cty.String),
													Optional:        true,
													DescriptionKind: "markdown",
													Description:     "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
												},
											},
										},
									},
								},
								DataSourceSchemas: map[string]*tfjson.Schema{
									"null_data_source": {
										Version: 0,
										Block: &tfjson.SchemaBlock{
											DescriptionKind: "markdown",
											Description: `The ` + "`null_data_source`" + ` data source implements the standard data source lifecycle but does not
interact with any external APIs.

Historically, the ` + "`null_data_source`" + ` was typically used to construct intermediate values to re-use elsewhere in configuration. The
same can now be achieved using [locals](https://www.terraform.io/docs/language/values/locals.html).
`,
											Deprecated: true,
											Attributes: map[string]*tfjson.SchemaAttribute{
												"has_computed_default": {
													AttributeType:   cty.String,
													Optional:        true,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
												},
												"id": {
													AttributeType:   cty.String,
													Optional:        false,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "This attribute is only present for some legacy compatibility issues and should not be used. It will be removed in a future version.",
													Deprecated:      true,
												},
												"inputs": {
													AttributeType:   cty.Map(cty.String),
													Optional:        true,
													DescriptionKind: "markdown",
													Description:     "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
												},
												"outputs": {
													AttributeType:   cty.Map(cty.String),
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "After the data source is \"read\", a copy of the `inputs` map.",
												},
												"random": {
													AttributeType:   cty.String,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "A random value. This is primarily for testing and has little practical use; prefer the [hashicorp/random provider](https://registry.terraform.io/providers/hashicorp/random) for more practical random number use-cases.",
												},
											},
										},
									},
								},
							},
						},
					}
				} else if tfv.Core().GreaterThanOrEqual(v0_13_0) {
					providerSchema = &tfjson.ProviderSchemas{
						FormatVersion: "0.1",
						Schemas: map[string]*tfjson.ProviderSchema{
							"registry.terraform.io/hashicorp/null": {
								ConfigSchema: &tfjson.Schema{
									Version: 0,
									Block: &tfjson.SchemaBlock{
										DescriptionKind: "plain",
									},
								},
								ResourceSchemas: map[string]*tfjson.Schema{
									"null_resource": {
										Version: 0,
										Block: &tfjson.SchemaBlock{
											DescriptionKind: "markdown",
											Description:     "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
											Attributes: map[string]*tfjson.SchemaAttribute{
												"id": {
													AttributeType:   cty.String,
													Optional:        false,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "This is set to a random value at create time.",
												},
												"triggers": {
													AttributeType:   cty.Map(cty.String),
													Optional:        true,
													DescriptionKind: "markdown",
													Description:     "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
												},
											},
										},
									},
								},
								DataSourceSchemas: map[string]*tfjson.Schema{
									"null_data_source": {
										Version: 0,
										Block: &tfjson.SchemaBlock{
											DescriptionKind: "markdown",
											Description: `The ` + "`null_data_source`" + ` data source implements the standard data source lifecycle but does not
interact with any external APIs.

Historically, the ` + "`null_data_source`" + ` was typically used to construct intermediate values to re-use elsewhere in configuration. The
same can now be achieved using [locals](https://www.terraform.io/docs/language/values/locals.html).
`,
											Deprecated: true,
											Attributes: map[string]*tfjson.SchemaAttribute{
												"has_computed_default": {
													AttributeType:   cty.String,
													Optional:        true,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
												},
												"id": {
													AttributeType:   cty.String,
													Optional:        false,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "This attribute is only present for some legacy compatibility issues and should not be used. It will be removed in a future version.",
													Deprecated:      true,
												},
												"inputs": {
													AttributeType:   cty.Map(cty.String),
													Optional:        true,
													DescriptionKind: "markdown",
													Description:     "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
												},
												"outputs": {
													AttributeType:   cty.Map(cty.String),
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "After the data source is \"read\", a copy of the `inputs` map.",
												},
												"random": {
													AttributeType:   cty.String,
													Computed:        true,
													DescriptionKind: "markdown",
													Description:     "A random value. This is primarily for testing and has little practical use; prefer the [hashicorp/random provider](https://registry.terraform.io/providers/hashicorp/random) for more practical random number use-cases.",
												},
											},
										},
									},
								},
							},
						},
					}
				} else {
					providerSchema = &tfjson.ProviderSchemas{
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
											DescriptionKind: "",
											Attributes: map[string]*tfjson.SchemaAttribute{
												"id": {
													AttributeType: cty.String,
													Optional:      false,
													Computed:      true,
													Description:   "This is set to a random value at create time.",
												},
												"triggers": {
													AttributeType:   cty.Map(cty.String),
													Optional:        true,
													DescriptionKind: "",
													Description:     "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
												},
											},
										},
									},
								},
								DataSourceSchemas: map[string]*tfjson.Schema{
									"null_data_source": {
										Version: 0,
										Block: &tfjson.SchemaBlock{
											DescriptionKind: "",
											Description:     "",
											Attributes: map[string]*tfjson.SchemaAttribute{
												"has_computed_default": {
													AttributeType:   cty.String,
													Optional:        true,
													Computed:        true,
													DescriptionKind: "",
													Description:     "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
												},
												"id": {
													AttributeType:   cty.String,
													Optional:        false,
													Computed:        true,
													DescriptionKind: "",
													Description:     "This attribute is only present for some legacy compatibility issues and should not be used. It will be removed in a future version.",
												},
												"inputs": {
													AttributeType:   cty.Map(cty.String),
													Optional:        true,
													DescriptionKind: "",
													Description:     "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
												},
												"outputs": {
													AttributeType:   cty.Map(cty.String),
													Computed:        true,
													DescriptionKind: "",
													Description:     "After the data source is \"read\", a copy of the `inputs` map.",
												},
												"random": {
													AttributeType:   cty.String,
													Computed:        true,
													DescriptionKind: "",
													Description:     "A random value. This is primarily for testing and has little practical use; prefer the [hashicorp/random provider](https://registry.terraform.io/providers/hashicorp/random) for more practical random number use-cases.",
												},
											},
										},
									},
								},
							},
						},
					}
				}

				return providerSchema
			},
		},
		{
			"empty_with_tf_file", func(tfv *version.Version) *tfjson.ProviderSchemas {
				// TODO: Add handling for v1 format once it lands in core
				// See https://github.com/hashicorp/terraform/pull/29550

				if tfv.Core().GreaterThanOrEqual(v0_15_0) {
					return &tfjson.ProviderSchemas{
						FormatVersion: "0.2",
						Schemas:       nil,
					}
				}

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
				if tfv.Core().LessThan(providersSchemaJSONMinVersion) {
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

				if diff := diffSchema(expected, schemas); diff != "" {
					t.Fatalf("mismatch (-want +got):\n%s", diff)
				}
			})
		})
	}

}

func TestProvidersSchema_versionMismatch(t *testing.T) {
	t.Skip("TODO! add version mismatch test for 0.11 as -json was added in 0.12 (I think)")
}
