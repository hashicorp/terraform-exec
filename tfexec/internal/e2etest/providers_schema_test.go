// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	v0_12_0 = version.Must(version.NewVersion("0.12.0"))
	v0_13_0 = version.Must(version.NewVersion("0.13.0"))
	v0_14_0 = version.Must(version.NewVersion("0.14.0"))
	v0_15_0 = version.Must(version.NewVersion("0.15.0"))
	v1_0    = version.Must(version.NewVersion("1.0.0"))
	v1_1    = version.Must(version.NewVersion("1.1.0"))
)

func TestProvidersSchema(t *testing.T) {
	for i, c := range []struct {
		fixtureDir string
		expected   func(*version.Version) *tfjson.ProviderSchemas
	}{
		{
			"basic", func(tfv *version.Version) *tfjson.ProviderSchemas {
				nullSchema := &tfjson.ProviderSchema{
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
										Optional:      false,
										Computed:      true,
										Description:   "This is set to a random value at create time.",
									},
									"triggers": {
										AttributeType: cty.Map(cty.String),
										Optional:      true,
										Description:   "A map of arbitrary strings that, when changed, will force the null resource to be replaced, re-running any associated provisioners.",
									},
								},
							},
						},
					},
					DataSourceSchemas: map[string]*tfjson.Schema{
						"null_data_source": {
							Version: 0,
							Block: &tfjson.SchemaBlock{
								Deprecated: false,
								Attributes: map[string]*tfjson.SchemaAttribute{
									"has_computed_default": {
										AttributeType: cty.String,
										Optional:      true,
										Computed:      true,
										Description:   "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
									},
									"id": {
										AttributeType: cty.String,
										Optional:      false,
										Computed:      true,
										Deprecated:    false,
										Description:   "This attribute is only present for some legacy compatibility issues and should not be used. It will be removed in a future version.",
									},
									"inputs": {
										AttributeType: cty.Map(cty.String),
										Optional:      true,
										Description:   "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
									},
									"outputs": {
										AttributeType: cty.Map(cty.String),
										Computed:      true,
										Description:   "After the data source is \"read\", a copy of the `inputs` map.",
									},
									"random": {
										AttributeType: cty.String,
										Computed:      true,
										Description:   "A random value. This is primarily for testing and has little practical use; prefer the [hashicorp/random provider](https://registry.opentofu.org/providers/hashicorp/random) for more practical random number use-cases.",
									},
								},
							},
						},
					},
				}

				formatVersion := "0.1"
				providerAddr := "null"

				if tfv.Core().GreaterThanOrEqual(v0_13_0) {
					providerAddr = "registry.opentofu.org/hashicorp/null"

					nullSchema = &tfjson.ProviderSchema{
						ConfigSchema: &tfjson.Schema{
							Version: 0,
							Block: &tfjson.SchemaBlock{
								DescriptionKind: tfjson.SchemaDescriptionKindPlain,
							},
						},
						ResourceSchemas: map[string]*tfjson.Schema{
							"null_resource": {
								Version: 0,
								Block: &tfjson.SchemaBlock{
									DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
									Description:     "The `null_resource` resource implements the standard resource lifecycle but takes no further action.\n\nThe `triggers` argument allows specifying an arbitrary set of values that, when changed, will cause the resource to be replaced.",
									Attributes: map[string]*tfjson.SchemaAttribute{
										"id": {
											AttributeType:   cty.String,
											Optional:        false,
											Computed:        true,
											DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
											Description:     "This is set to a random value at create time.",
										},
										"triggers": {
											AttributeType:   cty.Map(cty.String),
											Optional:        true,
											DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
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
									DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
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
											DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
											Description:     "If set, its literal value will be stored and returned. If not, its value defaults to `\"default\"`. This argument exists primarily for testing and has little practical use.",
										},
										"id": {
											AttributeType:   cty.String,
											Optional:        false,
											Computed:        true,
											DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
											Description:     "This attribute is only present for some legacy compatibility issues and should not be used. It will be removed in a future version.",
											Deprecated:      true,
										},
										"inputs": {
											AttributeType:   cty.Map(cty.String),
											Optional:        true,
											DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
											Description:     "A map of arbitrary strings that is copied into the `outputs` attribute, and accessible directly for interpolation.",
										},
										"outputs": {
											AttributeType:   cty.Map(cty.String),
											Computed:        true,
											DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
											Description:     "After the data source is \"read\", a copy of the `inputs` map.",
										},
										"random": {
											AttributeType:   cty.String,
											Computed:        true,
											DescriptionKind: tfjson.SchemaDescriptionKindMarkdown,
											Description:     "A random value. This is primarily for testing and has little practical use; prefer the [hashicorp/random provider](https://registry.opentofu.org/providers/hashicorp/random) for more practical random number use-cases.",
										},
									},
								},
							},
						},
					}
				}
				if tfv.Core().GreaterThanOrEqual(v0_15_0) {
					formatVersion = "0.2"
				}
				if tfv.Core().GreaterThanOrEqual(v1_1) {
					formatVersion = "1.0"
				}

				providerSchema := &tfjson.ProviderSchemas{
					FormatVersion: formatVersion,
					Schemas: map[string]*tfjson.ProviderSchema{
						providerAddr: nullSchema,
					},
				}

				return providerSchema
			},
		},
		{
			"empty_with_tf_file", func(tfv *version.Version) *tfjson.ProviderSchemas {
				if tfv.Core().GreaterThanOrEqual(v1_1) {
					return &tfjson.ProviderSchemas{
						FormatVersion: "1.0",
						Schemas:       nil,
					}
				} else if tfv.Core().GreaterThanOrEqual(v0_15_0) {
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
				if tfv.Core().LessThan(v0_12_0) {
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
