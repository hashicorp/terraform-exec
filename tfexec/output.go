package tfexec

import (
	"context"
	"encoding/json"
	"os/exec"
)

type outputConfig struct {
	state string
	name  string
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

func (opt *OutputNameOption) configureOutput(conf *outputConfig) {
	conf.name = opt.name
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
	c := defaultOutputOptions

	for _, o := range opts {
		o.configureOutput(&c)
	}
	outputCmd := tf.outputCmd(ctx, c)

	outputs := map[string]OutputMeta{}
	if c.name != "" {
		var outputValue json.RawMessage
		err := tf.runTerraformCmdJSON(ctx, outputCmd, &outputValue)
		if err != nil {
			return nil, err
		}
		output := OutputMeta{}
		output.Value = outputValue
		outputs[c.name] = output
	} else {
		err := tf.runTerraformCmdJSON(ctx, outputCmd, &outputs)
		if err != nil {
			return nil, err
		}
	}

	return outputs, nil
}

func (tf *Terraform) outputCmd(ctx context.Context, c outputConfig) *exec.Cmd {

	args := []string{"output", "-no-color", "-json"}

	// string opts: only pass if set
	if c.state != "" {
		args = append(args, "-state="+c.state)
	}

	// output name: only pass if set
	if c.name != "" {
		args = append(args, c.name)
	}

	return tf.buildTerraformCmd(ctx, nil, args...)
}
