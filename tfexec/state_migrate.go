// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"
)

const defaultLockFilePath = ".terraform.lock.hcl"

type stateMigrateConfig struct {
	srcProviderLockFile string
	dstProviderLockFile string

	upgrade bool
}

var defaultStateMigrateOptions = stateMigrateConfig{
	srcProviderLockFile: defaultLockFilePath,
	dstProviderLockFile: defaultLockFilePath,
}

// StateMigrateOption represents options used in the Refresh method.
type StateMigrateOption interface {
	configureStateMigrate(*stateMigrateConfig)
}

func (opt *SourceProviderLockFileOption) configureStateMigrate(conf *stateMigrateConfig) {
	conf.srcProviderLockFile = opt.path
}

func (opt *DestinationProviderLockFileOption) configureStateMigrate(conf *stateMigrateConfig) {
	conf.dstProviderLockFile = opt.path
}

func (opt *UpgradeOption) configureStateMigrate(conf *stateMigrateConfig) {
	conf.upgrade = opt.upgrade
}

// StateMigrateJSON executes `terraform state migrate` with `-json` flag
// and any specified options and waits for it to complete.
//
// StateMigrateJSON is likely to be removed in a future major version in favour of
// StateMigrate returning JSON by default.
func (tf *Terraform) StateMigrateJSON(ctx context.Context, opts ...StateMigrateOption) error {
	cmd, err := tf.stateMigrateJSONCmd(ctx, opts...)
	if err != nil {
		return err
	}

	// TODO: JSON parsing
	// TODO: error handling (stderr and parsing?)

	return tf.runTerraformCmd(ctx, cmd)
}

func (tf *Terraform) stateMigrateJSONCmd(ctx context.Context, opts ...StateMigrateOption) (*exec.Cmd, error) {
	c := defaultStateMigrateOptions

	for _, o := range opts {
		o.configureStateMigrate(&c)
	}

	args := tf.buildStateMigrateArgs(c)

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}

func (tf *Terraform) buildStateMigrateArgs(c stateMigrateConfig) []string {
	args := []string{"state", "migrate", "-no-color", "-input=false", "-json"}

	if c.srcProviderLockFile != "" {
		args = append(args, "-source-provider-lock-file="+c.srcProviderLockFile)
	}
	if c.dstProviderLockFile != "" {
		args = append(args, "-destination-provider-lock-file="+c.dstProviderLockFile)
	}
	if c.upgrade {
		args = append(args, "-upgrade")
	}
	return args
}
