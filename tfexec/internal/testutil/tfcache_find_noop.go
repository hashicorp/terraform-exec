// +build !go1.14

package testutil

import (
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

func (tf *TFCache) find(t *testing.T, key string, finder func(dir string) tfinstall.ExecPathFinder) string {
	panic("not implemented for " + runtime.Version())
}
