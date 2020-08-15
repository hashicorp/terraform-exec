package tfinstall

import "os"

type ExactPathOption struct {
	execPath string
}

func ExactPath(execPath string) *ExactPathOption {
	opt := &ExactPathOption{
		execPath: execPath,
	}
	return opt
}

func (opt *ExactPathOption) ExecPath() (string, error) {
	if _, err := os.Stat(opt.execPath); err != nil {
		// fall through to the next strategy if the local path does not exist
		return "", nil
	}
	return opt.execPath, nil
}
