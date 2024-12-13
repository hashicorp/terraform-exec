// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

type providersSchemaConfig struct {
	reattachInfo ReattachInfo
}

var defaultProvidersSchemaOptions = providersSchemaConfig{}

type ProvidersSchemaOption interface {
	configureProvidersSchema(*providersSchemaConfig)
}

func (opt *ReattachOption) configureProvidersSchema(conf *providersSchemaConfig) {
	conf.reattachInfo = opt.info
}

// ProvidersSchema represents the terraform providers schema -json subcommand.
func (tf *Terraform) ProvidersSchema(ctx context.Context, opts ...ProvidersSchemaOption) (*tfjson.ProviderSchemas, error) {
	c := defaultProvidersSchemaOptions
	for _, o := range opts {
		o.configureProvidersSchema(&c)
	}

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	schemaCmd := tf.providersSchemaCmd(ctx, mergeEnv)

	var ret tfjson.ProviderSchemas
	err := tf.runTerraformCmdJSON(ctx, schemaCmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *Terraform) providersSchemaCmd(ctx context.Context, mergeEnv map[string]string, args ...string) *exec.Cmd {
	allArgs := []string{"providers", "schema", "-json", "-no-color"}
	allArgs = append(allArgs, args...)

	return tf.buildTerraformCmd(ctx, mergeEnv, allArgs...)
}
