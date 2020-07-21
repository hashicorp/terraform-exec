package e2etest

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var err error
		installDir, err = ioutil.TempDir("", "tfinstall")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(installDir)

		return m.Run()
	}())
}
