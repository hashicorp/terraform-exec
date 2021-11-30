package tfexec

import (
	"context"
	"fmt"
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
	err := tf.compatible(ctx, tf0_4_1, nil)
	if err != nil {
		return fmt.Errorf("taint was first introduced in Terraform 0.4.1: %w", err)
	}
	taintCmd := tf.taintCmd(ctx, address, opts...)
	return tf.runTerraformCmd(ctx, taintCmd)
}

func (tf *Terraform) taintCmd(ctx context.Context, address string, opts ...TaintOption) *exec.Cmd {
	c := defaultTaintOptions

	for _, o := range opts {
		o.configureTaint(&c)
	}

	args := []string{"taint", "-no-color"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	args = append(args, address)

	return tf.buildTerraformCmd(ctx, nil, args...)
}
