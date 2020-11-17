package tfexec

import (
	"context"
	"fmt"
	"os/exec"
)

type upgrade012Config struct {
	yes   bool
	force bool

	reattachInfo ReattachInfo
}

var defaultUpgrade012Options = upgrade012Config{
	yes:   false,
	force: false,
}

// Upgrade012Option represents options used in the Destroy method.
type Upgrade012Option interface {
	configureUpgrade012(*upgrade012Config)
}

func (opt *ForceOption) configureUpgrade012(conf *upgrade012Config) {
	conf.force = opt.force
}

func (opt *YesOption) configureUpgrade012(conf *upgrade012Config) {
	conf.yes = opt.yes
}

func (opt *ReattachOption) configureUpgrade012(conf *upgrade012Config) {
	conf.reattachInfo = opt.info
}

// Upgrade012 represents the terraform 0.12upgrade subcommand.
func (tf *Terraform) Upgrade012(ctx context.Context, opts ...Upgrade012Option) error {
	cmd, err := tf.upgrade012Cmd(ctx, opts...)
	if err != nil {
		return err
	}
	return tf.runTerraformCmd(cmd)
}

func (tf *Terraform) upgrade012Cmd(ctx context.Context, opts ...Upgrade012Option) (*exec.Cmd, error) {
	err := tf.compatible(ctx, tf0_12_0, tf0_13_0)
	if err != nil {
		return nil, fmt.Errorf("terraform 0.12upgrade is only supported in 0.12 releases: %w", err)
	}

	c := defaultUpgrade012Options

	for _, o := range opts {
		o.configureUpgrade012(&c)
	}

	args := []string{"0.12upgrade"}

	// boolean opts: only pass if set
	if c.yes {
		args = append(args, "-yes")
	}
	if c.force {
		args = append(args, "-force")
	}

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
