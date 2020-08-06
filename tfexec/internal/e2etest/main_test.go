package e2etest

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

var tfcache *testutil.TFCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		installDir, err := ioutil.TempDir("", "tfinstall")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(installDir)

		tfcache = testutil.NewTFCache(installDir)
		return m.Run()
	}())
}
