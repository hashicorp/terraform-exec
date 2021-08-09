package tfexec

import (
	"bytes"
	"context"
	"os/exec"
)

type stateShowConfig struct {
	state string
}

var defaultStateShowOptions = stateShowConfig{}

// StateShowCmdOption represents options used in the StateShow method.
type StateShowCmdOption interface {
	configureStateShow(*stateShowConfig)
}

func (opt *StateOption) configureStateShow(conf *stateShowConfig) {
	conf.state = opt.path
}

// StateShow represents the terraform state show subcommand.
func (tf *Terraform) StateShow(ctx context.Context, address string, opts ...StateShowCmdOption) (string, error) {
	cmd, err := tf.stateShowCmd(ctx, address, opts...)
	if err != nil {
		return "", err
	}

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	err = tf.runTerraformCmd(ctx, cmd)
	if err != nil {
		return "", err
	}

	return outBuf.String(), nil
}

func (tf *Terraform) stateShowCmd(ctx context.Context, address string, opts ...StateShowCmdOption) (*exec.Cmd, error) {
	c := stateShowConfig{}

	for _, o := range opts {
		o.configureStateShow(&c)
	}

	args := []string{"state", "show", "-no-color"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}
	args = append(args, address)

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}
