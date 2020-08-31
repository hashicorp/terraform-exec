package tfexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// Console represents the console subcommand to the Terraform CLI.
func (tf *Terraform) Console(ctx context.Context, expression string) (string, error) {
	if runtime.GOOS == "windows" {
		return "", fmt.Errorf("terraform console does not support stdin scripting on Windows currently, see https://github.com/hashicorp/terraform/issues/18242")
	}

	// TODO: [DIR] param option
	// TODO: var option
	// TODO: var-file option
	cmd := tf.buildTerraformCmd(ctx, "console")

	expression = strings.TrimSpace(expression)
	expression += "\n"
	cmd.Stdin = strings.NewReader(expression)

	tf.logger.Printf("console expression: %q", expression)

	outbuf := bytes.Buffer{}
	cmd.Stdout = &outbuf

	err := tf.runTerraformCmd(cmd)
	if err != nil {
		return "", err
	}

	return outbuf.String(), nil
}

// ConsoleJSON is a convenience method for invoking Console with the expression wrapped with
// jsonencode and the result unmarshaled in to the passed value.
func (tf *Terraform) ConsoleJSON(ctx context.Context, expression string, v interface{}) error {
	expression = fmt.Sprintf("jsonencode(%s)", expression)

	raw, err := tf.Console(ctx, expression)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(raw), v)
	if err != nil {
		return err
	}

	return nil
}
