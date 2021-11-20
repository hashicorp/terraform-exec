package tfexec

import (
	"context"
	"os/exec"
)

type untaintConfig struct {
	state string
}

var defaultUnTaintOptions = untaintConfig{}

// OutputOption represents options used in the Output method.
type UnTaintOption interface {
	configureUnTaint(*untaintConfig)
}

func (opt *StateOption) configureUnTaint(conf *untaintConfig) {
	conf.state = opt.path
}

// Untaint represents the terraform untaint subcommand.
func (tf *Terraform) UnTaint(ctx context.Context, address string, opts ...UnTaintOption) error {
	unTaintCmd := tf.unTaintCmd(ctx, address, opts...)
	return tf.runTerraformCmd(ctx, unTaintCmd)
}

func (tf *Terraform) unTaintCmd(ctx context.Context, address string, opts ...UnTaintOption) *exec.Cmd {
	c := defaultUnTaintOptions

	for _, o := range opts {
		o.configureUnTaint(&c)
	}

	args := []string{"untaint", "-allow-missing", "-lock-timeout=0s", "-lock=true", "-no-color"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	args = append(args, address)

	return tf.buildTerraformCmd(ctx, nil, args...)
}
