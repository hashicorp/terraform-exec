package tfexec

import (
	"context"
	"strings"
)

// WorkspaceShow represents the workspace show subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceShow(ctx context.Context) (string, error) {
	workspaceShowCmd := tf.buildTerraformCmd(ctx, nil, "workspace", "show", "-no-color")

	var outBuffer strings.Builder
	workspaceShowCmd.Stdout = &outBuffer

	err := tf.runTerraformCmd(ctx, workspaceShowCmd)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(outBuffer.String()), nil
}
