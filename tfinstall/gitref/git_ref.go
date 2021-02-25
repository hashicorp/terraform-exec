package gitref

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Option struct {
	installDir string
	repoURL    string
	ref        string
}

func Install(ref, repo, installDir string) *Option {
	return &Option{
		installDir: installDir,
		repoURL:    repo,
		ref:        ref,
	}
}

func (opt *Option) ExecPath(ctx context.Context) (string, error) {
	installDir, err := ensureInstallDir(opt.installDir)
	if err != nil {
		return "", err
	}

	ref := plumbing.ReferenceName(opt.ref)
	if opt.ref == "" {
		ref = plumbing.ReferenceName("refs/heads/main")
	}

	repoURL := opt.repoURL
	if repoURL == "" {
		repoURL = "https://github.com/hashicorp/terraform.git"
	}

	_, err = git.PlainClone(installDir, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: ref,

		Depth: 1,
		Tags:  git.NoTags,
	})
	if err != nil {
		return "", fmt.Errorf("unable to clone %q: %w", repoURL, err)
	}

	var binName string
	{
		// TODO: maybe there is a better way to make sure this filename is available?
		// I guess we could locate it in a different dir, or nest the git underneath
		// the root tmp dir, etc.
		binPattern := "terraform"
		if runtime.GOOS == "windows" {
			binPattern = "terraform*.exe"
		}
		binFile, err := ioutil.TempFile(installDir, binPattern)
		if err != nil {
			return "", fmt.Errorf("unable to create bin file: %w", err)
		}
		binName = binFile.Name()
		binFile.Close()
	}

	goArgs := []string{"build", "-o", binName}

	vendorDir := filepath.Join(installDir, "vendor")
	if fi, err := os.Stat(vendorDir); err == nil && fi.IsDir() {
		goArgs = append(goArgs, "-mod", "vendor")
	}

	cmd := exec.CommandContext(ctx, "go", goArgs...)
	cmd.Dir = installDir
	out, err := cmd.CombinedOutput()
	log.Print(string(out))
	if err != nil {
		return "", fmt.Errorf("unable to build Terraform: %w\n%s", err, out)
	}

	return binName, nil
}

func ensureInstallDir(installDir string) (string, error) {
	if installDir == "" {
		return ioutil.TempDir("", "tfexec")
	}

	if _, err := os.Stat(installDir); err != nil {
		return "", fmt.Errorf("could not access directory %s for installing Terraform: %w", installDir, err)
	}

	return installDir, nil
}
