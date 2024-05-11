// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import "syscall"

var defaultSysProcAttr = &syscall.SysProcAttr{
	// kill children if parent is dead
	Pdeathsig: syscall.SIGKILL,
	// set process group ID
	Setpgid: true,
}
