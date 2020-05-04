package tfexec

import (
	"os"
	"os/exec"
)

// FindTerraform attempts to find a Terraform CLI executable.
//
// As a first preference it will look for the environment variable
// TFEXEC_TERRAFORM_PATH and return its value. If that variable is not set, it will
// look in PATH for a program named "terraform" and, if one is found, return
// its absolute path.
//
// If no Terraform executable can be found, the result is the empty string. In
// that case, the test program will usually fail outright.
func FindTerraform() string {
	if p := os.Getenv("TFEXEC_TERRAFORM_PATH"); p != "" {
		return p
	}
	p, err := exec.LookPath("terraform")
	if err != nil {
		return ""
	}
	return p
}
