// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"os/exec"
	"syscall"
)

func (tf *Terraform) configureCommandWindow(cmd *exec.Cmd) {
	if tf.hideWindow {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
}
