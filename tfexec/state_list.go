package tfexec

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

type stateListConfig struct {
	state     string
	id        string
	addresses []string
}

// StateListCmdOption represents options used in the StateList method.
type StateListCmdOption interface {
	configureStateList(*stateListConfig)
}

func (opt *StateOption) configureStateList(conf *stateListConfig) {
	conf.state = opt.path
}

func (opt *IdOption) configureStateList(conf *stateListConfig) {
	conf.id = opt.id
}

func (opt *AddressOption) configureStateList(conf *stateListConfig) {
	if conf.addresses == nil {
		conf.addresses = []string{}
	}

	conf.addresses = append(conf.addresses, opt.address)
}

// StateList represents the terraform state list subcommand.
func (tf *Terraform) StateList(ctx context.Context, opts ...StateListCmdOption) ([]string, error) {
	cmd, err := tf.stateListCmd(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	err = tf.runTerraformCmd(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return parseStateList(outBuf.String()), nil
}

func (tf *Terraform) stateListCmd(ctx context.Context, opts ...StateListCmdOption) (*exec.Cmd, error) {
	c := stateListConfig{}

	for _, o := range opts {
		o.configureStateList(&c)
	}

	args := []string{"state", "list", "-no-color"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}
	if c.id != "" {
		args = append(args, "-id="+c.id)
	}

	for _, address := range c.addresses {
		args = append(args, address)
	}

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}

func parseStateList(stdout string) []string {
	lines := strings.Split(stdout, "\n")

	addresses := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		addresses = append(addresses, line)
	}

	return addresses
}
