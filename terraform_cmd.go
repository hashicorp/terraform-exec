package tfexec

import (
	"os"
	"os/exec"
)

func (t *Terraform) buildTerraformCmd(args ...string) exec.Cmd {
	allArgs := []string{"terraform"}
	allArgs = append(allArgs, args...)
	allArgs = append(allArgs, "-no-color")

	var env []string
	for _, e := range os.Environ() {
		env = append(env, e)
	}

	env = append(env, "TF_LOG=") // so logging can't pollute our stderr output
	env = append(env, "TF_INPUT=0")

	return exec.Cmd{
		Path: t.execPath,
		Env:  t.Env,
		Args: allArgs,
		Dir:  t.workingDir,
	}
}

func (t *Terraform) InitCmd(args ...string) exec.Cmd {
	allArgs := []string{"init"}
	allArgs = append(allArgs, args...)

	return t.buildTerraformCmd(allArgs...)
}

func (t *Terraform) ShowCmd(args ...string) exec.Cmd {
	allArgs := []string{"show", "-json"}
	allArgs = append(allArgs, args...)

	return t.buildTerraformCmd(allArgs...)
}

func (t *Terraform) ProvidersSchemaCmd(args ...string) exec.Cmd {
	allArgs := []string{"providers", "schema", "-json"}
	allArgs = append(allArgs, args...)

	return t.buildTerraformCmd(allArgs...)
}
