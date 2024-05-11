//go:build unix && !linux

package tfexec

import "syscall"

var defaultSysProcAttr = &syscall.SysProcAttr{
	// set process group ID
	Setpgid: true,
}
