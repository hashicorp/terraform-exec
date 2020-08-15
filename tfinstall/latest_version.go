package tfinstall

import (
	"fmt"

	"github.com/hashicorp/go-checkpoint"
)

type LatestVersionOption struct {
	forceCheckpoint bool
	installDir      string
}

func LatestVersion(installDir string, forceCheckpoint bool) *LatestVersionOption {
	opt := &LatestVersionOption{
		forceCheckpoint: forceCheckpoint,
		installDir:      installDir,
	}

	return opt
}

func (opt *LatestVersionOption) ExecPath() (string, error) {
	v, err := latestVersion(opt.forceCheckpoint)
	if err != nil {
		return "", err
	}

	return downloadWithVerification(v, opt.installDir)
}

func latestVersion(forceCheckpoint bool) (string, error) {
	resp, err := checkpoint.Check(&checkpoint.CheckParams{
		Product: "terraform",
		Force:   forceCheckpoint,
	})
	if err != nil {
		return "", err
	}

	if resp.CurrentVersion == "" {
		return "", fmt.Errorf("could not determine latest version of terraform using checkpoint: CHECKPOINT_DISABLE may be set")
	}

	return resp.CurrentVersion, nil
}
