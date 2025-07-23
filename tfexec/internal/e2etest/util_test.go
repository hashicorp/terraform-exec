// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

const testFixtureDir = "testdata"

var (
	showMinVersion = version.Must(version.NewVersion("0.12.0"))

	providerAddressMinVersion = version.Must(version.NewVersion("0.13.0"))

	metadataFunctionsMinVersion = version.Must(version.NewVersion("1.4.0"))
)

func runTest(t *testing.T, fixtureName string, cb func(t *testing.T, tfVersion *version.Version, tf *tfexec.Terraform)) {
	t.Helper()

	versions := []string{
		testutil.Latest011,
		testutil.Latest012,
		testutil.Latest013,
		testutil.Latest014,
		testutil.Latest015,
		testutil.Latest_v1,
		testutil.Latest_v1_1,
		testutil.Latest_v1_2,
		testutil.Latest_v1_3,
		testutil.Latest_v1_4,
		testutil.Latest_v1_5,
		testutil.Latest_v1_6,
		testutil.Latest_v1_7,
		testutil.Latest_v1_8,
		testutil.Latest_v1_9,
		testutil.Latest_v1_10,
		testutil.Latest_v1_11,
		testutil.Latest_v1_12,
	}
	if override := os.Getenv("TFEXEC_E2ETEST_VERSIONS"); override != "" {
		versions = strings.Split(override, ",")
	}

	// If the env var TFEXEC_E2ETEST_TERRAFORM_PATH is set to the path of a
	// valid Terraform executable, only tests appropriate to that
	// executable's version will be run.
	if localBinPath := os.Getenv("TFEXEC_E2ETEST_TERRAFORM_PATH"); localBinPath != "" {
		// By convention, every new Terraform struct is given a clean
		// temp dir, even if we are only invoking tf.Version(). This
		// prevents any possible confusion that could result from
		// reusing an os.TempDir() (for example) that already contained
		// Terraform files.
		td, err := ioutil.TempDir("", "tf")
		if err != nil {
			t.Fatalf("error creating temporary test directory: %s", err)
		}
		t.Cleanup(func() {
			os.RemoveAll(td)
		})
		ltf, err := tfexec.NewTerraform(td, localBinPath)
		if err != nil {
			t.Fatal(err)
		}

		ltf.SetAppendUserAgent("tfexec-e2etest")

		lVersion, _, err := ltf.Version(context.Background(), false)
		if err != nil {
			t.Fatalf("unable to determine version of Terraform binary at %s: %s", localBinPath, err)
		}

		versions = []string{lVersion.String()}
	}

	runTestWithVersions(t, versions, fixtureName, cb)
}

func runTestWithVersions(t *testing.T, versions []string, fixtureName string, cb func(t *testing.T, tfVersion *version.Version, tf *tfexec.Terraform)) {
	t.Helper()

	alreadyRunVersions := map[string]bool{}
	for _, tfv := range versions {
		t.Run(fmt.Sprintf("%s-%s", fixtureName, tfv), func(t *testing.T) {
			if !strings.HasPrefix(tfv, "refs/") && runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
				v, err := version.NewVersion(tfv)
				if err != nil {
					t.Fatal(err)
				}
				if v.LessThan(version.Must(version.NewVersion("1.0.2"))) {
					t.Skipf("Terraform not available for darwin/arm64 < 1.0.2 (%s)", v)
				}
			}

			if alreadyRunVersions[tfv] {
				t.Skipf("already run version %q", tfv)
			}
			alreadyRunVersions[tfv] = true

			td, err := ioutil.TempDir("", "tf")
			if err != nil {
				t.Fatalf("error creating temporary test directory: %s", err)
			}
			t.Cleanup(func() {
				os.RemoveAll(td)
			})

			var execPath string
			if localBinPath := os.Getenv("TFEXEC_E2ETEST_TERRAFORM_PATH"); localBinPath != "" {
				execPath = localBinPath
			} else if strings.HasPrefix(tfv, "refs/") {
				execPath = tfcache.GitRef(t, tfv)
			} else {
				execPath = tfcache.Version(t, tfv)
			}

			tf, err := tfexec.NewTerraform(td, execPath)
			if err != nil {
				t.Fatal(err)
			}

			tf.SetAppendUserAgent("tfexec-e2etest")

			runningVersion, _, err := tf.Version(context.Background(), false)
			if err != nil {
				t.Fatalf("unable to determine running version (expected %q): %s", tfv, err)
			}

			// Check that the runningVersion matches the expected
			// test version. This ensures non-matching tests are
			// skipped when using a local Terraform executable.
			if !strings.HasPrefix(tfv, "refs/") {
				testVersion, err := version.NewVersion(tfv)
				if err != nil {
					t.Fatalf("unable to parse version %s: %s", testVersion, err)
				}
				if !testVersion.Equal(runningVersion) {
					t.Skipf("test applies to version %s, but local executable is version %s", tfv, runningVersion)
				}
			}

			if fixtureName != "" {
				err = copyFiles(filepath.Join(testFixtureDir, fixtureName), td)
				if err != nil {
					t.Fatalf("error copying config file into test dir: %s", err)
				}
			}

			// Separate strings.Builder because it's not concurrent safe
			var stdout strings.Builder
			tf.SetStdout(&stdout)
			var stderr strings.Builder
			tf.SetStderr(&stderr)

			tf.SetLogger(&testingPrintfer{t})

			// TODO: capture panics here?
			cb(t, runningVersion, tf)

			t.Logf("CLI Output:\n%s", stdout.String())
			if len(stderr.String()) > 0 {
				t.Logf("CLI Error:\n%s", stderr.String())
			}
		})
	}
}

type testingPrintfer struct {
	t *testing.T
}

func (t *testingPrintfer) Printf(format string, v ...interface{}) {
	t.t.Logf(format, v...)
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

// filesEqual asserts that two files have the same contents.
func textFilesEqual(t *testing.T, expected, actual string) {
	eb, err := ioutil.ReadFile(expected)
	if err != nil {
		t.Fatal(err)
	}

	ab, err := ioutil.ReadFile(actual)
	if err != nil {
		t.Fatal(err)
	}

	es := string(eb)
	es = strings.ReplaceAll(es, "\r\n", "\n")

	as := string(ab)
	as = strings.ReplaceAll(as, "\r\n", "\n")

	if as != es {
		t.Fatalf("expected:\n%s\n\ngot:\n%s\n", es, as)
	}
}

func checkSum(t *testing.T, filename string) uint32 {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return crc32.ChecksumIEEE(b)
}
