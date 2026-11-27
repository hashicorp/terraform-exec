// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

type modulesConfig struct {
	vars     []string
	varFiles []string
}

var defaultModulesOptions = modulesConfig{}

// ModulesOption represents options used in the Modules method.
type ModulesOption interface {
	configureModules(*modulesConfig)
}

func (opt *VarFileOption) configureModules(conf *modulesConfig) {
	conf.varFiles = append(conf.varFiles, opt.path)
}

func (opt *VarOption) configureModules(conf *modulesConfig) {
	conf.vars = append(conf.vars, opt.assignment)
}

// Modules represents the terraform modules -json subcommand.
func (tf *Terraform) Modules(ctx context.Context, opts ...ModulesOption) (*tfjson.ModulesOutput, error) {
	err := tf.compatible(ctx, tf1_10_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform modules -json was added in 1.10.0: %w", err)
	}

	modulesCmd, err := tf.modulesCmd(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var ret tfjson.ModulesOutput
	err = tf.runTerraformCmdJSON(ctx, modulesCmd, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *Terraform) modulesCmd(ctx context.Context, opts ...ModulesOption) (*exec.Cmd, error) {
	c := defaultModulesOptions

	for _, o := range opts {
		switch o.(type) {
		case *VarOption:
			// https://github.com/hashicorp/terraform/pull/38276
			err := tf.compatible(ctx, tf1_15_0, nil)
			if err != nil {
				return nil, fmt.Errorf("-var was added to modules in Terraform 1.15: %w", err)
			}
		case *VarFileOption:
			// https://github.com/hashicorp/terraform/pull/38276
			err := tf.compatible(ctx, tf1_15_0, nil)
			if err != nil {
				return nil, fmt.Errorf("-var-file was added to modules in Terraform 1.15: %w", err)
			}
		}

		o.configureModules(&c)
	}

	args, err := tf.buildModulesArgs(c)
	if err != nil {
		return nil, err
	}

	args = append(args, "-json")
	return tf.buildTerraformCmd(ctx, nil, args...), nil
}

func (tf *Terraform) buildModulesArgs(c modulesConfig) ([]string, error) {
	args := []string{"modules"}

	for _, vf := range c.varFiles {
		args = append(args, "-var-file="+vf)
	}

	for _, v := range c.vars {
		args = append(args, "-var", v)
	}

	return args, nil
}
