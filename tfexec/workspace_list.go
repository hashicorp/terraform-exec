// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"os/exec"
	"strings"
)

type workspaceListConfig struct {
	reattachInfo ReattachInfo
}

var defaultWorkspaceListOptions = workspaceListConfig{}

type WorkspaceListOption interface {
	configureWorkspaceList(*workspaceListConfig)
}

func (opt *ReattachOption) configureWorkspaceList(conf *workspaceListConfig) {
	conf.reattachInfo = opt.info
}

// WorkspaceList represents the workspace list subcommand to the Terraform CLI.
func (tf *Terraform) WorkspaceList(ctx context.Context, opts ...WorkspaceListOption) ([]string, string, error) {
	wlCmd, err := tf.workspaceListCmd(ctx, opts...)
	if err != nil {
		return nil, "", err
	}

	var outBuf strings.Builder
	wlCmd.Stdout = &outBuf

	err = tf.runTerraformCmd(ctx, wlCmd)
	if err != nil {
		return nil, "", err
	}

	// Parse human output into a list of workspaces and the current workspace
	ws, current := parseWorkspaceListHumanOutput(outBuf.String())

	return ws, current, nil
}

func (tf *Terraform) workspaceListCmd(ctx context.Context, opts ...WorkspaceListOption) (*exec.Cmd, error) {
	c := defaultWorkspaceListOptions

	for _, o := range opts {
		o.configureWorkspaceList(&c)
	}

	args := []string{"workspace", "list", "-no-color"}
	return tf.buildWorkspaceListCmd(ctx, c, args)
}

// parseWorkspaceListHumanOutput parses human output from the workspace list command
// to return a slice of workspace names and the current workspace name
func parseWorkspaceListHumanOutput(stdout string) ([]string, string) {
	var currentWorkspacePrefix string = "* "

	lines := strings.Split(stdout, "\n")

	current := ""
	workspaces := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, currentWorkspacePrefix) {
			line = strings.TrimPrefix(line, currentWorkspacePrefix)
			current = line
		}
		workspaces = append(workspaces, line)
	}

	return workspaces, current
}

func (tf *Terraform) buildWorkspaceListCmd(ctx context.Context, c workspaceListConfig, args []string) (*exec.Cmd, error) {
	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}
