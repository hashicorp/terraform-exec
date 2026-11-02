// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

type stateIdentitiesConfig struct {
	reattachInfo ReattachInfo
	state        string
}

var defaultStateIdentitiesOptions = stateIdentitiesConfig{}

// StateIdentitiesOption represents options used in the StateIdentitiesJSON method.
type StateIdentitiesOption interface {
	configureStateIdentities(*stateIdentitiesConfig)
}

func (opt *ReattachOption) configureStateIdentities(conf *stateIdentitiesConfig) {
	conf.reattachInfo = opt.info
}

func (opt *StateOption) configureStateIdentities(conf *stateIdentitiesConfig) {
	conf.state = opt.path
}

// StateIdentities represents the terraform state identities subcommand with the `-json` flag.
func (tf *Terraform) StateIdentities(ctx context.Context, opts ...StateIdentitiesOption) (*tfjson.StateIdentitiesOutput, error) {
	// Terraform supported this since 1.12 but the JSON format has changed
	// in 1.18 to better support long-term maintenance so we opt to treat
	// that as a starting version.
	// See https://github.com/hashicorp/terraform/pull/39443
	err := tf.compatible(ctx, tf1_18_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform state identities -json was added in 1.18.0: %w", err)
	}

	cmd, err := tf.stateIdentitiesCmd(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var output tfjson.StateIdentitiesOutput

	err = tf.runTerraformCmdJSON(ctx, cmd, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (tf *Terraform) stateIdentitiesCmd(ctx context.Context, opts ...StateIdentitiesOption) (*exec.Cmd, error) {
	c := defaultStateIdentitiesOptions

	for _, o := range opts {
		o.configureStateIdentities(&c)
	}

	args := []string{"state", "identities", "-no-color", "-json"}

	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	return tf.buildStateIdentitiesCmd(ctx, c, args)
}

func (tf *Terraform) buildStateIdentitiesCmd(ctx context.Context, c stateIdentitiesConfig, args []string) (*exec.Cmd, error) {
	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}
