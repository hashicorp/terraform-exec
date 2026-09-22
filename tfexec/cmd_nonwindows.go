// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package tfexec

import "os/exec"

func (tf *Terraform) configureCommandWindow(cmd *exec.Cmd) {}
