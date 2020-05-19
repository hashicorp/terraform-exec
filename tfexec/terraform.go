package tfexec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type Terraform struct {
	execPath   string
	workingDir string

	// Log each Terraform command with fmt.Println. For use in tests and debugging.
	echo bool
	Env  []string
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

type applyConfig struct {
	Backup    string
	DirOrPlan string
	Lock      bool

	// LockTimeout must be a string with time unit, e.g. '10s'
	LockTimeout string
	Parallelism int
	Refresh     bool
	State       string
	StateOut    string
	Targets     []string

	// Vars: each var must be supplied as a single string, e.g. 'foo=bar'
	Vars    []string
	VarFile string
}

var defaultApplyOptions = applyConfig{
	Lock:        true,
	Parallelism: 10,
	Refresh:     true,
}

type ApplyOption interface {
	configureApply(*applyConfig)
}

type ParallelismOption struct {
	parallelism int
}

type BackupOption struct {
	backup string
}

type TargetOption struct {
	target string
}

func (opt *ParallelismOption) configureApply(conf *applyConfig) {
	conf.Parallelism = opt.parallelism
}

func (opt *BackupOption) configureApply(conf *applyConfig) {
	conf.Backup = opt.backup
}

func (opt *TargetOption) configureApply(conf *applyConfig) {
	conf.Targets = append(conf.Targets, opt.target)
}

func Parallelism(n int) *ParallelismOption {
	return &ParallelismOption{n}
}

func Backup(path string) *BackupOption {
	return &BackupOption{path}
}

func Target(resource string) *TargetOption {
	return &TargetOption{resource}
}

func (t *Terraform) Apply(opts ...ApplyOption) error {
	c := &defaultApplyOptions

	for _, o := range opts {
		o.configureApply(c)
	}

	args := []string{}

	// string args: only pass if set
	if c.Backup != "" {
		args = append(args, "-backup="+c.Backup)
	}
	if c.LockTimeout != "" {
		args = append(args, "-lock-timeout="+c.LockTimeout)
	}
	if c.State != "" {
		args = append(args, "-state="+c.State)
	}
	if c.StateOut != "" {
		args = append(args, "-state-out="+c.StateOut)
	}
	if c.VarFile != "" {
		args = append(args, "-var-file="+c.VarFile)
	}

	// boolean and numerical args: always pass
	args = append(args, "-lock="+strconv.FormatBool(c.Lock))

	args = append(args, "-parallelism="+fmt.Sprint(c.Parallelism))
	args = append(args, "-refresh="+strconv.FormatBool(c.Refresh))

	// string slice args: pass as separate args
	if c.Targets != nil {
		for _, ta := range c.Targets {
			args = append(args, "-target="+ta)
		}
	}

	if c.Vars != nil {
		for _, v := range c.Vars {
			args = append(args, "-var '"+v+"'")
		}
	}

	applyCmd := t.ApplyCmd(args...)

	if t.echo {
		fmt.Println(applyCmd.String())
	}

	var errBuf strings.Builder
	applyCmd.Stderr = &errBuf

	err := applyCmd.Run()
	if err != nil {
		return errors.New(errBuf.String())
	}

	return nil
}

type planConfig struct {
	Destroy     bool
	Lock        bool
	LockTimeout string
	Out         string
	Parallelism int
	Refresh     bool
	State       string
	Targets     []string
	Vars        []string
	VarFile     string
}

type PlanOption interface {
	configurePlan(*planConfig)
}

func (opt *ParallelismOption) configurePlan(conf *planConfig) {
	conf.Parallelism = opt.parallelism
}

func (t *Terraform) Plan(opts ...PlanOption) error {
	return nil
}

func (t *Terraform) Init(args ...string) error {
	initCmd := t.InitCmd(args...)

	var errBuf strings.Builder
	initCmd.Stderr = &errBuf

	err := initCmd.Run()
	if err != nil {
		return errors.New(errBuf.String())
	}

	return nil
}

func (t *Terraform) Show(args ...string) (*tfjson.State, error) {
	var ret tfjson.State

	var errBuf strings.Builder
	var outBuf bytes.Buffer

	showCmd := t.ShowCmd(args...)

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

func (t *Terraform) ProvidersSchema(args ...string) (*tfjson.ProviderSchemas, error) {
	var ret tfjson.ProviderSchemas

	var errBuf strings.Builder
	var outBuf bytes.Buffer

	schemaCmd := t.ProvidersSchemaCmd(args...)

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
