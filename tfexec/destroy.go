package tfexec

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type destroyConfig struct {
	backup string
	lock   bool

	// LockTimeout must be a string with time unit, e.g. '10s'
	lockTimeout string
	parallelism int
	refresh     bool
	state       string
	stateOut    string
	targets     []string

	// Vars: each var must be supplied as a single string, e.g. 'foo=bar'
	vars    []string
	varFile string
}

var defaultDestroyOptions = destroyConfig{
	lock:        true,
	lockTimeout: "0s",
	parallelism: 10,
	refresh:     true,
}

type DestroyOption interface {
	configureDestroy(*destroyConfig)
}

func (opt *ParallelismOption) configureDestroy(conf *destroyConfig) {
	conf.parallelism = opt.parallelism
}

func (opt *BackupOption) configureDestroy(conf *destroyConfig) {
	conf.backup = opt.path
}

func (opt *TargetOption) configureDestroy(conf *destroyConfig) {
	conf.targets = append(conf.targets, opt.target)
}

func (opt *LockTimeoutOption) configureDestroy(conf *destroyConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *StateOption) configureDestroy(conf *destroyConfig) {
	conf.state = opt.path
}

func (opt *StateOutOption) configureDestroy(conf *destroyConfig) {
	conf.stateOut = opt.path
}

func (opt *VarFileOption) configureDestroy(conf *destroyConfig) {
	conf.varFile = opt.path
}

func (opt *LockOption) configureDestroy(conf *destroyConfig) {
	conf.lock = opt.lock
}

func (opt *RefreshOption) configureDestroy(conf *destroyConfig) {
	conf.refresh = opt.refresh
}

func (opt *VarOption) configureDestroy(conf *destroyConfig) {
	conf.vars = append(conf.vars, opt.assignment)
}

func (tf *Terraform) Destroy(ctx context.Context, opts ...DestroyOption) error {
	destroyCmd := tf.destroyCmd(ctx, opts...)

	var errBuf strings.Builder
	destroyCmd.Stderr = &errBuf

	err := destroyCmd.Run()
	if err != nil {
		return parseError(err, errBuf.String())
	}

	return nil
}

func (tf *Terraform) destroyCmd(ctx context.Context, opts ...DestroyOption) *exec.Cmd {
	c := defaultDestroyOptions

	for _, o := range opts {
		o.configureDestroy(&c)
	}

	args := []string{"destroy", "-no-color", "-auto-approve", "-input=false"}

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
	if c.varFile != "" {
		args = append(args, "-var-file="+c.varFile)
	}

	// boolean and numerical opts: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.lock))
	args = append(args, "-parallelism="+fmt.Sprint(c.parallelism))
	args = append(args, "-refresh="+strconv.FormatBool(c.refresh))

	// string slice opts: split into separate args
	if c.targets != nil {
		for _, ta := range c.targets {
			args = append(args, "-target="+ta)
		}
	}
	if c.vars != nil {
		for _, v := range c.vars {
			args = append(args, "-var", v)
		}
	}

	return tf.buildTerraformCmd(ctx, args...)
}
