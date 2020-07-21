package e2etest

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

const testFixtureDir = "testdata"
const testTerraformStateFileName = "terraform.tfstate"

func testTempDir(t *testing.T) string {
	d, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("error creating temporary test directory: %s", err)
	}

	return d
}

func copyFiles(path string, dstPath string) error {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			// TODO: make recursive with filepath.Walk?
			continue
		}
		err = copyFile(filepath.Join(path, info.Name()), dstPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(path string, dstPath string) error {
	srcF, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcF.Close()

	di, err := os.Stat(dstPath)
	if err != nil {
		return err
	}
	if di.IsDir() {
		_, file := filepath.Split(path)
		dstPath = filepath.Join(dstPath, file)
	}

	dstF, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}

type installedVersion struct {
	path string
	err  error
}

var (
	installDir           string
	installedVersionLock sync.Mutex
	installedVersions    = map[string]installedVersion{}
)

func tfVersion(t *testing.T, v string) string {
	if installDir == "" {
		t.Fatalf("installDir not yet configured, TestMain must run first")
	}

	installedVersionLock.Lock()
	defer installedVersionLock.Unlock()

	iv, ok := installedVersions[v]
	if !ok {
		dir := filepath.Join(installDir, v)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			t.Fatal(err)
		}
		iv.path, iv.err = tfinstall.Find(tfinstall.ExactVersion(v, dir))
		installedVersions[v] = iv
	}

	if iv.err != nil {
		t.Fatalf("error installing terraform version %q: %s", v, iv.err)
	}

	return iv.path
}
