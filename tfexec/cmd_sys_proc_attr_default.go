//go:build !linux
// +build !linux

package tfexec

import "syscall"

var defaultSysProcAttr = &syscall.SysProcAttr{}
