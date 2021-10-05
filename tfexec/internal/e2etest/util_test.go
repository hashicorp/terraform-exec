package e2etest

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	}
	if override := os.Getenv("TFEXEC_E2ETEST_VERSIONS"); override != "" {
		versions = strings.Split(override, ",")
	}

	runTestVersions(t, versions, fixtureName, cb)
}

func runTestVersions(t *testing.T, versions []string, fixtureName string, cb func(t *testing.T, tfVersion *version.Version, tf *tfexec.Terraform)) {
	t.Helper()

	alreadyRunVersions := map[string]bool{}
	for _, tfv := range versions {
		t.Run(fmt.Sprintf("%s-%s", fixtureName, tfv), func(t *testing.T) {
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

			// TODO: do this in a cleaner way than string comparison?
			var execPath string
			switch {
			case strings.HasPrefix(tfv, "refs/"):
				execPath = tfcache.GitRef(t, tfv)
			default:
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

			if fixtureName != "" {
				err = copyFiles(filepath.Join(testFixtureDir, fixtureName), td)
				if err != nil {
					t.Fatalf("error copying config file into test dir: %s", err)
				}
			}

			var stdouterr strings.Builder
			tf.SetStdout(&stdouterr)
			tf.SetStderr(&stdouterr)

			tf.SetLogger(&testingPrintfer{t})

			// TODO: capture panics here?
			cb(t, runningVersion, tf)

			t.Logf("CLI Output:\n%s", stdouterr.String())
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
