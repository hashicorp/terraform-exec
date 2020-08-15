package tfinstall

import "github.com/hashicorp/go-version"

type ExactVersionOption struct {
	tfVersion  string
	installDir string
}

func ExactVersion(tfVersion string, installDir string) *ExactVersionOption {
	opt := &ExactVersionOption{
		tfVersion:  tfVersion,
		installDir: installDir,
	}

	return opt
}

func (opt *ExactVersionOption) ExecPath() (string, error) {
	// validate version
	_, err := version.NewVersion(opt.tfVersion)
	if err != nil {
		return "", err
	}

	return downloadWithVerification(opt.tfVersion, opt.installDir)
}
