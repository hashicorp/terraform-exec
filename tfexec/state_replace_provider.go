package tfexec

import (
	"context"
	"os/exec"
	"strconv"
)

type stateReplaceProviderConfig struct {
	backup      string
	lock        bool
	lockTimeout string
	state       string
	stateOut    string
}

var defaultStateReplaceProviderOptions = stateReplaceProviderConfig{
	lock:        true,
	lockTimeout: "0s",
}

// StateReplaceProviderCmdOption represents options used in the Refresh method.
type StateReplaceProviderCmdOption interface {
	configureStateReplaceProvider(*stateReplaceProviderConfig)
}

func (opt *BackupOption) configureStateReplaceProvider(conf *stateReplaceProviderConfig) {
	conf.backup = opt.path
}

func (opt *LockOption) configureStateReplaceProvider(conf *stateReplaceProviderConfig) {
	conf.lock = opt.lock
}

func (opt *LockTimeoutOption) configureStateReplaceProvider(conf *stateReplaceProviderConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *StateOption) configureStateReplaceProvider(conf *stateReplaceProviderConfig) {
	conf.state = opt.path
}

func (opt *StateOutOption) configureStateReplaceProvider(conf *stateReplaceProviderConfig) {
	conf.stateOut = opt.path
}

// StateMv represents the terraform state mv subcommand.
func (tf *Terraform) StateReplaceProvider(ctx context.Context, fromProviderFqn string, toProviderFqn string, opts ...StateReplaceProviderCmdOption) error {
	cmd, err := tf.stateReplaceProviderCmd(ctx, fromProviderFqn, toProviderFqn, opts...)
	if err != nil {
		return err
	}
	return tf.runTerraformCmd(ctx, cmd)
}

func (tf *Terraform) stateReplaceProviderCmd(ctx context.Context, fromProviderFqn string, toProviderFqn string, opts ...StateReplaceProviderCmdOption) (*exec.Cmd, error) {
	c := defaultStateReplaceProviderOptions

	for _, o := range opts {
		o.configureStateReplaceProvider(&c)
	}

	args := []string{"state", "replace-provider", "-no-color", "-auto-approve"}

	// string opts: only pass if set
	if c.backup != "" {
		args = append(args, "-backup="+c.backup)
	}
	if c.lockTimeout != "" {
		args = append(args, "-lock-timeout="+c.lockTimeout)
	}
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}
	if c.stateOut != "" {
		args = append(args, "-state-out="+c.stateOut)
	}

	// boolean and numerical opts: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.lock))

	// positional arguments
	args = append(args, fromProviderFqn)
	args = append(args, toProviderFqn)

	return tf.buildTerraformCmd(ctx, nil, args...), nil
}
