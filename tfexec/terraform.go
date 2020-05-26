package tfexec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type Terraform struct {
	execPath   string
	workingDir string
	Env        []string
}

// NewTerraform returns a Terraform struct with default values for all fields.
// If a blank execPath is supplied, NewTerraform will attempt to locate an
// appropriate binary on the system PATH.
func NewTerraform(workingDir string, execPath string) (*Terraform, error) {
	var err error
	if workingDir == "" {
		return nil, fmt.Errorf("Terraform cannot be initialised with empty workdir")
	}

	if _, err := os.Stat(workingDir); err != nil {
		return nil, fmt.Errorf("error initialising Terraform with workdir %s: %s", workingDir, err)
	}

	if execPath == "" {
		execPath, err = FindTerraform()
		if err != nil {
			return nil, err
		}
	}
	// TODO for maximum helpfulness, check execPath looks like a terraform binary

	passthroughEnv := os.Environ()

	return &Terraform{
		execPath:   execPath,
		workingDir: workingDir,
		Env:        passthroughEnv,
	}, nil
}

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
	initCmd := t.InitCmd(ctx, opts...)

	var errBuf strings.Builder
	initCmd.Stderr = &errBuf

	err := initCmd.Run()
	if err != nil {
		return errors.New(errBuf.String())
	}

	return nil
}

type applyConfig struct {
	backup    string
	dirOrPlan string
	lock      bool

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

var defaultApplyOptions = applyConfig{
	lock:        true,
	parallelism: 10,
	refresh:     true,
}

type ApplyOption interface {
	configureApply(*applyConfig)
}

func (opt *ParallelismOption) configureApply(conf *applyConfig) {
	conf.parallelism = opt.parallelism
}

func (opt *BackupOption) configureApply(conf *applyConfig) {
	conf.backup = opt.backup
}

func (opt *TargetOption) configureApply(conf *applyConfig) {
	conf.targets = append(conf.targets, opt.target)
}

func (opt *LockTimeoutOption) configureApply(conf *applyConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *StateOption) configureApply(conf *applyConfig) {
	conf.state = opt.path
}

func (opt *StateOutOption) configureApply(conf *applyConfig) {
	conf.stateOut = opt.path
}

func (opt *VarFileOption) configureApply(conf *applyConfig) {
	conf.varFile = opt.path
}

func (opt *LockOption) configureApply(conf *applyConfig) {
	conf.lock = opt.lock
}

func (opt *RefreshOption) configureApply(conf *applyConfig) {
	conf.refresh = opt.refresh
}

func (opt *VarOption) configureApply(conf *applyConfig) {
	conf.vars = append(conf.vars, opt.assignment)
}

func (opt *DirOrPlanOption) configureApply(conf *applyConfig) {
	conf.dirOrPlan = opt.path
}

func (tf *Terraform) Apply(ctx context.Context, opts ...ApplyOption) error {
	applyCmd := tf.ApplyCmd(ctx, opts...)

	var errBuf strings.Builder
	applyCmd.Stderr = &errBuf

	err := applyCmd.Run()
	if err != nil {
		return errors.New(errBuf.String())
	}

	return nil
}

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
	conf.backup = opt.backup
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
	destroyCmd := tf.DestroyCmd(ctx, opts...)

	var errBuf strings.Builder
	destroyCmd.Stderr = &errBuf

	err := destroyCmd.Run()
	if err != nil {
		return errors.New(errBuf.String())
	}

	return nil
}

type planConfig struct {
	destroy     bool
	lock        bool
	lockTimeout string
	out         string
	parallelism int
	refresh     bool
	state       string
	targets     []string
	vars        []string
	varFile     string
}

var defaultPlanOptions = planConfig{
	destroy:     false,
	lock:        true,
	lockTimeout: "0s",
	parallelism: 10,
	refresh:     true,
}

type PlanOption interface {
	configurePlan(*planConfig)
}

func (opt *VarFileOption) configurePlan(conf *planConfig) {
	conf.varFile = opt.path
}

func (opt *VarOption) configurePlan(conf *planConfig) {
	conf.vars = append(conf.vars, opt.assignment)
}

func (opt *TargetOption) configurePlan(conf *planConfig) {
	conf.targets = append(conf.targets, opt.target)
}

func (opt *StateOption) configurePlan(conf *planConfig) {
	conf.state = opt.path
}

func (opt *RefreshOption) configurePlan(conf *planConfig) {
	conf.refresh = opt.refresh
}

func (opt *ParallelismOption) configurePlan(conf *planConfig) {
	conf.parallelism = opt.parallelism
}

func (opt *OutOption) configurePlan(conf *planConfig) {
	conf.out = opt.path
}

func (opt *LockTimeoutOption) configurePlan(conf *planConfig) {
	conf.lockTimeout = opt.timeout
}

func (opt *LockOption) configurePlan(conf *planConfig) {
	conf.lock = opt.lock
}

func (opt *DestroyFlagOption) configurePlan(conf *planConfig) {
	conf.destroy = opt.destroy
}

func (tf *Terraform) Plan(ctx context.Context, opts ...PlanOption) error {
	planCmd := tf.PlanCmd(ctx, opts...)

	var errBuf strings.Builder
	planCmd.Stderr = &errBuf

	err := planCmd.Run()
	if err != nil {
		return errors.New(errBuf.String())
	}

	return nil
}

func (t *Terraform) Show(ctx context.Context, args ...string) (*tfjson.State, error) {
	var ret tfjson.State

	var errBuf strings.Builder
	var outBuf bytes.Buffer

	showCmd := t.ShowCmd(ctx, args...)

	showCmd.Stderr = &errBuf
	showCmd.Stdout = &outBuf

	err := showCmd.Run()
	if err != nil {
		if tErr, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("terraform failed: %s\n\nstderr:\n%s", tErr.ProcessState.String(), errBuf.String())
		}
		return nil, err
	}

	err = json.Unmarshal(outBuf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (t *Terraform) ProvidersSchema(ctx context.Context, args ...string) (*tfjson.ProviderSchemas, error) {
	var ret tfjson.ProviderSchemas

	var errBuf strings.Builder
	var outBuf bytes.Buffer

	schemaCmd := t.ProvidersSchemaCmd(ctx, args...)

	schemaCmd.Stderr = &errBuf
	schemaCmd.Stdout = &outBuf

	err := schemaCmd.Run()
	if err != nil {
		if tErr, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("terraform failed: %s\n\nstderr:\n%s", tErr.ProcessState.String(), errBuf.String())
		}
		return nil, err
	}

	err = json.Unmarshal(outBuf.Bytes(), ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
