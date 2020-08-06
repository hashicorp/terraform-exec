package e2etest

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
)

const testFixtureDir = "testdata"

func setupFixture(t *testing.T, version, name string) (*tfexec.Terraform, func()) {
	t.Helper()

	td, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("error creating temporary test directory: %s", err)
	}
	// TODO: make this a t.Cleanup once we no longer support Go 1.13
	cleanup := func() {
		os.RemoveAll(td)
	}

	tf, err := tfexec.NewTerraform(td, tfcache.Version(t, version))
	if err != nil {
		t.Fatal(err)
	}

	if name != "" {
		err = copyFiles(filepath.Join(testFixtureDir, name), td)
		if err != nil {
			t.Fatalf("error copying config file into test dir: %s", err)
		}
	}

	return tf, cleanup
}

func copyFiles(path string, dstPath string) error {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, info := range infos {
		srcPath := filepath.Join(path, info.Name())
		if info.IsDir() {
			newDir := filepath.Join(dstPath, info.Name())
			err = os.MkdirAll(newDir, info.Mode())
			if err != nil {
				return err
			}
			err = copyFiles(srcPath, newDir)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
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

// filesEqual returns true iff the two files have the same contents.
func filesEqual(file1, file2 string) (bool, error) {
	sf, err := os.Open(file1)
	if err != nil {
		return false, err
	}

	df, err := os.Open(file2)
	if err != nil {
		return false, err
	}

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return true, nil
		}
	}

	return false, nil
}
