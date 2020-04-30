package tfexec

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	getter "github.com/hashicorp/go-getter"
)

const releaseHost = "https://releases.hashicorp.com"

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

func tfURL(version, osName, archName string) string {
	return fmt.Sprintf(
		"%s/terraform/%s/terraform_%s_%s_%s.zip",
		releaseHost, version, version, osName, archName,
	)
}

// InstallTerraform downloads and decompresses a Terraform CLI executable with
// the specified version, downloaded from the HashiCorp releases page over HTTP.
//
// The version string must match an existing Terraform release semver version,
// e.g. 0.12.5.
//
// The terraform executable is installed to a temporary folder.
// TODO: method for cleaning up this temporary folder.
func InstallTerraform(tfVersion string) (string, error) {
	osName := runtime.GOOS
	archName := runtime.GOARCH

	var tfDir string
	var err error

	tfDir, err = ioutil.TempDir("", "tfexec")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %s", err)
	}

	url := tfURL(tfVersion, osName, archName)

	client := getter.Client{
		Src: url,
		Dst: tfDir,

		Mode: getter.ClientModeDir,
	}

	err = client.Get()
	if err != nil {
		return "", fmt.Errorf("failed to download terraform from %s: %s", url, err)
	}

	return filepath.Join(tfDir, "terraform"), nil
}
