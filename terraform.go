package tfexec

import (
	"bytes"
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
