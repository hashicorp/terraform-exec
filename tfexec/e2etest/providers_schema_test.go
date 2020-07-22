package e2etest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-exec/tfexec"
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
					"null": &tfjson.ProviderSchema{
						ConfigSchema: &tfjson.Schema{
							Version: 0,
							Block:   &tfjson.SchemaBlock{},
						},
						ResourceSchemas: map[string]*tfjson.Schema{
							"null_resource": &tfjson.Schema{
								Version: 0,
								Block: &tfjson.SchemaBlock{
									Attributes: map[string]*tfjson.SchemaAttribute{
										"id": &tfjson.SchemaAttribute{
											AttributeType: cty.String,
											Optional:      true,
											Computed:      true,
										},
										"triggers": &tfjson.SchemaAttribute{
											AttributeType: cty.Map(cty.String),
											Optional:      true,
										},
									},
								},
							},
						},
						DataSourceSchemas: map[string]*tfjson.Schema{
							"null_data_source": &tfjson.Schema{
								Version: 0,
								Block: &tfjson.SchemaBlock{
									Attributes: map[string]*tfjson.SchemaAttribute{
										"has_computed_default": &tfjson.SchemaAttribute{
											AttributeType: cty.String,
											Optional:      true,
											Computed:      true,
										},
										"id": &tfjson.SchemaAttribute{
											AttributeType: cty.String,
											Optional:      true,
											Computed:      true,
										},
										"inputs": &tfjson.SchemaAttribute{
											AttributeType: cty.Map(cty.String),
											Optional:      true,
										},
										"outputs": &tfjson.SchemaAttribute{
											AttributeType: cty.Map(cty.String),
											Computed:      true,
										},
										"random": &tfjson.SchemaAttribute{
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
			"empty", &tfjson.ProviderSchemas{
				FormatVersion: "0.1",
				Schemas:       nil,
			},
		},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.fixtureDir), func(t *testing.T) {
			td := testTempDir(t)
			defer os.RemoveAll(td)

			tf, err := tfexec.NewTerraform(td, tfVersion(t, "0.12.28"))
			if err != nil {
				t.Fatal(err)
			}

			err = copyFiles(filepath.Join(testFixtureDir, c.fixtureDir), td)
			if err != nil {
				t.Fatalf("error copying fixtures into test dir: %s", err)
			}

			err = tf.Init(context.Background())
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
	}

}
