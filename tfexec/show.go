package tfexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

func (tf *Terraform) Show(ctx context.Context) (*tfjson.State, error) {
	err := tf.compatible(ctx, tf0_12_0, nil)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json was added in 0.12.0: %w", err)
	}

	var ret tfjson.State

	var errBuf strings.Builder
	var outBuf bytes.Buffer

	showCmd := tf.showCmd(ctx)

	showCmd.Stderr = &errBuf
	showCmd.Stdout = &outBuf

	err = showCmd.Run()
	if err != nil {
		return nil, parseError(errBuf.String())
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

func (tf *Terraform) showCmd(ctx context.Context, args ...string) *exec.Cmd {
	allArgs := []string{"show", "-json", "-no-color"}
	allArgs = append(allArgs, args...)

	return tf.buildTerraformCmd(ctx, allArgs...)
}
