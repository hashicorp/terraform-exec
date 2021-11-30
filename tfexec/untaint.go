package tfexec

import (
	"context"
	"fmt"
	"os/exec"
)

type untaintConfig struct {
	state string
}

var defaultUntaintOptions = untaintConfig{}

// OutputOption represents options used in the Output method.
type UntaintOption interface {
	configureUntaint(*untaintConfig)
}

func (opt *StateOption) configureUntaint(conf *untaintConfig) {
	conf.state = opt.path
}

// Untaint represents the terraform untaint subcommand.
func (tf *Terraform) Untaint(ctx context.Context, address string, opts ...UntaintOption) error {
	err := tf.compatible(ctx, tf0_6_13, nil)
	if err != nil {
		return fmt.Errorf("untaint was first introduced in Terraform 0.6.13: %w", err)
	}
	untaintCmd := tf.untaintCmd(ctx, address, opts...)
	return tf.runTerraformCmd(ctx, untaintCmd)
}

func (tf *Terraform) untaintCmd(ctx context.Context, address string, opts ...UntaintOption) *exec.Cmd {
	c := defaultUntaintOptions

	for _, o := range opts {
		o.configureUntaint(&c)
	}

	args := []string{"untaint", "-no-color"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	args = append(args, address)

	return tf.buildTerraformCmd(ctx, nil, args...)
}
