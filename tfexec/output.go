package tfexec

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type outputConfig struct {
	state string
	json  bool
}

var defaultOutputOptions = outputConfig{}

// OutputOption represents options used in the Output method.
type OutputOption interface {
	configureOutput(*outputConfig)
}

func (opt *StateOption) configureOutput(conf *outputConfig) {
	conf.state = opt.path
}

// OutputMeta represents the JSON output of 'terraform output -json',
// which resembles state format version 3 due to a historical accident.
// Please see hashicorp/terraform/command/output.go.
// TODO KEM: Should this type be in terraform-json?
type OutputMeta struct {
	Sensitive bool            `json:"sensitive"`
	Type      json.RawMessage `json:"type"`
	Value     json.RawMessage `json:"value"`
}

// Output represents the terraform output subcommand.
func (tf *Terraform) Output(ctx context.Context, opts ...OutputOption) (map[string]OutputMeta, error) {
	outputCmd := tf.outputCmd(ctx, opts...)

	outputs := map[string]OutputMeta{}
	err := tf.runTerraformCmdJSON(ctx, outputCmd, &outputs)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

// Output represents the terraform output subcommand.
func (tf *Terraform) OutputRaw(ctx context.Context, FieldToOutput string) string {
	fmt.Println("FieldToOutput:", FieldToOutput)
	outputCmd := tf.outputCmdRaw(ctx)
	outputs := tf.runTerraformCmdString(ctx, outputCmd)
	fmt.Println("Outputraw:", outputs)
	return outputs
}

func (tf *Terraform) outputCmdRaw(ctx context.Context) *exec.Cmd {

	args := []string{"output", "-no-color", "-raw"}
	return tf.buildTerraformCmd(ctx, nil, args...)
}

func (tf *Terraform) outputCmd(ctx context.Context, opts ...OutputOption) *exec.Cmd {
	c := defaultOutputOptions

	for _, o := range opts {
		o.configureOutput(&c)
	}

	args := []string{"output", "-no-color", "-json"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	return tf.buildTerraformCmd(ctx, nil, args...)
}
