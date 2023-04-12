// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package e2etest contains end-to-end acceptance tests for the tfexec
// package. It aims to cover as many realistic use cases for tfexec as possible;
// to serve as a smoke test for the incidental usage of hc-install with tfexec;
// and, eventually, to define the known and expected behaviour of the entire
// Terraform CLI.
// Test files inside the tfexec package are intended as unit tests covering the
// behaviour of *Cmd functions.
package e2etest
