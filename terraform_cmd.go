package tfexec

import (
	"os"
	"os/exec"
)

func buildTerraformCmd(workingDir string, args ...string) exec.Cmd {
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
		Path: FindTerraform(),
		Env:  env,
		Args: allArgs,
		Dir:  workingDir,
	}
}

func InitCmd(workingDir string, args ...string) exec.Cmd {
	allArgs := []string{"init"}
	allArgs = append(allArgs, args...)

	return buildTerraformCmd(workingDir, allArgs...)
}

func ShowCmd(workingDir string, args ...string) exec.Cmd {
	allArgs := []string{"show", "-json"}
	allArgs = append(allArgs, args...)

	return buildTerraformCmd(workingDir, allArgs...)
}

func ProvidersSchemaCmd(workingDir string, args ...string) exec.Cmd {
	allArgs := []string{"providers", "schema", "-json"}
	allArgs = append(allArgs, args...)

	return buildTerraformCmd(workingDir, allArgs...)
}
