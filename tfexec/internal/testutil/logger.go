// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestLogger() *log.Logger {
	if testing.Verbose() {
		return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}
	return log.New(ioutil.Discard, "", 0)
}
