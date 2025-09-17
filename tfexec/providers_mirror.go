package tfexec

import (
	"context"
	"fmt"
	"os/exec"
)

type providersMirrorConfig struct {
	platforms []string
}

var defaultProviderMirrorOptions = providersMirrorConfig{}

// ProvidersMirrorCmdOption represents options used in the mirror method.
type ProvidersMirrorCmdOption interface {
	configureProvidersMirror(*providersMirrorConfig)
}

func (opt *PlatformOption) configureProvidersMirror(conf *providersMirrorConfig) {
	conf.platforms = append(conf.platforms, opt.platform)
}

// ProvidersMirror represents the terraform providers mirror subcommand.
func (tf *Terraform) ProvidersMirror(ctx context.Context, targetDir string, opts ...ProvidersMirrorCmdOption) error {
	err := tf.compatible(ctx, tf0_13_0, nil)
	if err != nil {
		return fmt.Errorf("mirror was first introduced in Terraform 0.13.0: %w", err)
	}

	mirrorCmd := tf.providersMirrorCmd(ctx, targetDir, opts...)

	return tf.runTerraformCmd(ctx, mirrorCmd)
}

func (tf *Terraform) providersMirrorCmd(ctx context.Context, targetDir string, opts ...ProvidersMirrorCmdOption) *exec.Cmd {
	c := defaultProviderMirrorOptions

	for _, o := range opts {
		o.configureProvidersMirror(&c)
	}

	args := []string{"providers", "mirror", targetDir}

	for _, p := range c.platforms {
		args = append(args, "-platform="+p)
	}

	return tf.buildTerraformCmd(ctx, nil, args...)
}
