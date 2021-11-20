package tfexec

import (
	"context"
	"os/exec"
)

type taintConfig struct {
	state string
}

var defaultTaintOptions = taintConfig{}

// TaintOption represents options used in the Taint method.
type TaintOption interface {
	configureTaint(*taintConfig)
}

func (opt *StateOption) configureTaint(conf *taintConfig) {
	conf.state = opt.path
}

// Taint represents the terraform taint subcommand.
func (tf *Terraform) Taint(ctx context.Context, address string, opts ...TaintOption) error {
	taintCmd := tf.taintCmd(ctx, address, opts...)
	return tf.runTerraformCmd(ctx, taintCmd)
}

func (tf *Terraform) taintCmd(ctx context.Context, address string, opts ...TaintOption) *exec.Cmd {
	c := defaultTaintOptions

	for _, o := range opts {
		o.configureTaint(&c)
	}

	args := []string{"taint", "-allow-missing", "-lock-timeout=0s", "-lock=true"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	args = append(args, address)

	return tf.buildTerraformCmd(ctx, nil, args...)
}
