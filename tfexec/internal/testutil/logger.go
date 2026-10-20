// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestLogger() *log.Logger {
	if testing.Verbose() {
		return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}
	return log.New(io.Discard, "", 0)
}
