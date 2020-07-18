package tfexec

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Terraform struct {
	execPath    string
	workingDir  string
	execVersion string
	env         map[string]string

	logger  *log.Logger
	logPath string
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
		err := fmt.Errorf("NewTerraform: please supply the path to a Terraform executable using execPath, e.g. using the tfinstall package.")
		return nil, &ErrNoSuitableBinary{err: err}

	}
	tf := Terraform{
		execPath:   execPath,
		workingDir: workingDir,
		env:        nil, // explicit nil means copy os.Environ
		logger:     log.New(ioutil.Discard, "", 0),
	}

	execVersion, err := tf.version()
	if err != nil {
		return nil, &ErrNoSuitableBinary{err: fmt.Errorf("error running 'terraform version': %s", err)}
	}

	tf.execVersion = execVersion

	return &tf, nil
}

// SetEnv allows you to override environment variables, this should not be used for any well known
// Terraform environment variables that are already covered in options. Pass nil to copy the values
// from os.Environ. Attempting to set environment variables that should be managed manually will
// result in ErrManualEnvVar being returned.
func (tf *Terraform) SetEnv(env map[string]string) error {
	for k := range env {
		if strings.HasPrefix(k, varEnvVarPrefix) {
			return fmt.Errorf("variables should be passed using the Var option: %w", &ErrManualEnvVar{k})
		}
		for _, p := range prohibitedEnvVars {
			if p == k {
				return &ErrManualEnvVar{k}
			}
		}
	}

	tf.env = env
	return nil
}

func (tf *Terraform) SetLogger(logger *log.Logger) {
	tf.logger = logger
}

// SetLogPath sets the TF_LOG_PATH environment variable for Terraform CLI
// execution.
func (tf *Terraform) SetLogPath(path string) error {
	tf.logPath = path
	return nil
}

func (tf *Terraform) version() (string, error) {
	versionCmd := tf.buildTerraformCmd(context.Background(), "version")

	var errBuf strings.Builder
	var outBuf bytes.Buffer
	versionCmd.Stderr = &errBuf
	versionCmd.Stdout = &outBuf

	err := versionCmd.Run()
	if err != nil {
		fmt.Println(errBuf.String())
		return "", fmt.Errorf("%s, %s", err, errBuf.String())
	}

	return outBuf.String(), nil
}
