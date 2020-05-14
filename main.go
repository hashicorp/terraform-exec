package tfexec

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// FindTerraform attempts to find a Terraform CLI executable.
//
// As a first preference it will look for the environment variable
// TFEXEC_TERRAFORM_PATH and return its value. If that variable is not set, it will
// look in PATH for a program named "terraform.exe" on Windows, or "terraform" otherwise,
// and, if one is found, return its absolute path.
func FindTerraform() (string, error) {
	if p := os.Getenv("TFEXEC_TERRAFORM_PATH"); p != "" {
		return p, nil
	}

	execName := "terraform"

	if runtime.GOOS == "windows" {
		execName = "terraform.exe"
	}

	p, err := exec.LookPath(execName)
	if err != nil {
		return "", fmt.Errorf("terraform executable could not be found: %s", err)
	}
	return p, nil
}
