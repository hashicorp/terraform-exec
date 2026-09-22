// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"os/exec"
)

type workspaceSelectConfig struct {
	orCreate     bool
	reattachInfo ReattachInfo
}

var defaultWorkspaceSelectOptions = workspaceSelectConfig{}

type WorkspaceSelectOption interface {
	configureWorkspaceSelect(*workspaceSelectConfig)
}

func (opt *ReattachOption) configureWorkspaceSelect(conf *workspaceSelectConfig) {
	conf.reattachInfo = opt.info
}

func (opt *OrCreateOption) configureWorkspaceSelect(conf *workspaceSelectConfig) {
	conf.orCreate = opt.orCreate
}

// WorkspaceSelect represents the workspace select subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceSelect(ctx context.Context, workspace string, opts ...WorkspaceSelectOption) error {
	cmd, err := tf.workspaceSelectCmd(ctx, workspace, opts...)
	if err != nil {
		return err
	}

	return tf.runTerraformCmd(ctx, cmd)
}

func (tf *Terraform) workspaceSelectCmd(ctx context.Context, workspace string, opts ...WorkspaceSelectOption) (*exec.Cmd, error) {
	c := defaultWorkspaceSelectOptions

	for _, o := range opts {
		o.configureWorkspaceSelect(&c)
	}

	mergeEnv := map[string]string{}
	args := []string{"workspace", "select", "-no-color"}
	if c.orCreate {
		if err := tf.compatible(ctx, tf1_4_0, nil); err != nil {
			return nil, fmt.Errorf("-or-create was added to workspace select in Terraform 1.4: %w", err)
		}
		args = append(args, "-or-create")
	}
	args = append(args, workspace)
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}
