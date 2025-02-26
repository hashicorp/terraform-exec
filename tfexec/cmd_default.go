// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !unix
// +build !unix

package tfexec

import "syscall"

var defaultSysProcAttr = &syscall.SysProcAttr{}
