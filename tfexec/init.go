package tfexec

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type initConfig struct {
	backend       bool
	backendConfig []string
	forceCopy     bool
	fromModule    string
	get           bool
	getPlugins    bool
	lock          bool
	lockTimeout   string
	pluginDir     []string
	reconfigure   bool
	upgrade       bool
	verifyPlugins bool
}

var defaultInitOptions = initConfig{
	backend:       true,
	forceCopy:     false,
	get:           true,
	getPlugins:    true,
	lock:          true,
	lockTimeout:   "0s",
	reconfigure:   false,
	upgrade:       false,
	verifyPlugins: true,
}

type InitOption interface {
	configureInit(*initConfig)
}

func (opt *BackendOption) configureInit(conf *initConfig) {
	conf.backend = opt.backend
}

func (opt *BackendConfigOption) configureInit(conf *initConfig) {
	conf.backendConfig = append(conf.backendConfig, opt.path)
}

func (opt *FromModuleOption) configureInit(conf *initConfig) {
	conf.fromModule = opt.source
}

func (opt *GetOption) configureInit(conf *initConfig) {
	conf.get = opt.get
}

func (opt *GetPluginsOption) configureInit(conf *initConfig) {
	conf.getPlugins = opt.getPlugins
}

func (opt *LockOption) configureInit(conf *initConfig) {
	conf.lock = opt.lock
}

func (opt *LockTimeoutOption) configureInit(conf *initConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *PluginDirOption) configureInit(conf *initConfig) {
	conf.pluginDir = append(conf.pluginDir, opt.pluginDir)
}

func (opt *ReconfigureOption) configureInit(conf *initConfig) {
	conf.reconfigure = opt.reconfigure
}

func (opt *UpgradeOption) configureInit(conf *initConfig) {
	conf.upgrade = opt.upgrade
}

func (opt *VerifyPluginsOption) configureInit(conf *initConfig) {
	conf.verifyPlugins = opt.verifyPlugins
}

func (t *Terraform) Init(ctx context.Context, opts ...InitOption) error {
	initCmd := t.initCmd(ctx, opts...)

	var errBuf strings.Builder
	initCmd.Stderr = &errBuf

	err := initCmd.Run()
	if err != nil {
		return parseError(err, errBuf.String())
	}

	return nil
}

func (tf *Terraform) initCmd(ctx context.Context, opts ...InitOption) *exec.Cmd {
	c := defaultInitOptions

	for _, o := range opts {
		o.configureInit(&c)
	}

	args := []string{"init", "-no-color", "-force-copy", "-input=false"}

	// string opts: only pass if set
	if c.fromModule != "" {
		args = append(args, "-from-module="+c.fromModule)
	}
	if c.lockTimeout != "" {
		args = append(args, "-lock-timeout="+c.lockTimeout)
	}

	// boolean opts: always pass
	args = append(args, "-backend="+fmt.Sprint(c.backend))
	args = append(args, "-get="+fmt.Sprint(c.get))
	args = append(args, "-get-plugins="+fmt.Sprint(c.getPlugins))
	args = append(args, "-lock="+fmt.Sprint(c.lock))
	args = append(args, "-upgrade="+fmt.Sprint(c.upgrade))
	args = append(args, "-verify-plugins="+fmt.Sprint(c.verifyPlugins))

	// unary flags: pass if true
	if c.reconfigure {
		args = append(args, "-reconfigure")
	}

	// string slice opts: split into separate args
	if c.backendConfig != nil {
		for _, bc := range c.backendConfig {
			args = append(args, "-backend-config="+bc)
		}
	}
	if c.pluginDir != nil {
		for _, pd := range c.pluginDir {
			args = append(args, "-plugin-dir="+pd)
		}
	}

	return tf.buildTerraformCmd(ctx, args...)
}
