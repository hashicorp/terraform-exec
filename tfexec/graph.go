package tfexec

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type graphConfig struct {
	plan        string
	drawCycles  bool
	graphType   string
	moduleDepth int
}

var defaultGraphOptions = graphConfig{
	moduleDepth: -1,
}

type GraphOption interface {
	configureGraph(*graphConfig)
}

func (opt *GraphPlanOption) configureGraph(conf *graphConfig) {
	conf.plan = opt.file
}

func (opt *DrawCyclesOption) configureGraph(conf *graphConfig) {
	conf.drawCycles = opt.drawCycles
}

func (opt *GraphTypeOption) configureGraph(conf *graphConfig) {
	conf.graphType = opt.graphType
}

func (opt *ModuleDepthOption) configureGraph(conf *graphConfig) {
	conf.moduleDepth = opt.moduleDepth
}

func (tf *Terraform) Graph(ctx context.Context, opts ...GraphOption) (string, error) {
	graphCmd := tf.graphCmd(ctx, opts...)
	var outBuf strings.Builder
	graphCmd.Stdout = &outBuf
	err := tf.runTerraformCmd(ctx, graphCmd)
	if err != nil {
		return "", err
	}

	return outBuf.String(), nil

}

func (tf *Terraform) graphCmd(ctx context.Context, opts ...GraphOption) *exec.Cmd {
	c := defaultGraphOptions

	for _, o := range opts {
		o.configureGraph(&c)
	}

	args := []string{"graph"}

	if c.plan != "" {
		args = append(args, "-plan="+c.plan)
	}

	if c.drawCycles {
		args = append(args, "-draw-cycles")
	}

	if c.graphType != "" {
		args = append(args, "-type="+c.graphType)
	}

	// -1 is the default value set in terraform CLI for module depth
	if c.moduleDepth != -1 {
		args = append(args, fmt.Sprintf("-module-depth=%d", c.moduleDepth))
	}

	return tf.buildTerraformCmd(ctx, nil, args...)
}
