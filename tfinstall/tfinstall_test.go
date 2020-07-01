package tfinstall

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
)

// downloads terraform 0.12.26 from the live releases site
func TestFindExactVersion(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tfinstall-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tfpath, err := Find(ExactVersion("0.12.26", tmpDir))
	if err != nil {
		t.Fatal(err)
	}

	// run "terraform version" to check we've downloaded a terraform 0.12.26 binary
	cmd := exec.Command(tfpath, "version")

	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	expected := "Terraform v0.12.26"
	actual := string(out)
	if !strings.HasPrefix(actual, expected) {
		t.Fatalf("ran terraform version, expected %s, but got %s", expected, actual)
	}
}

// downloads terraform 0.13.0-beta1 from the live releases site
func TestFindExactVersionPrerelease(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tfinstall-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tfpath, err := Find(ExactVersion("0.13.0-beta1", tmpDir))
	if err != nil {
		t.Fatal(err)
	}

	// run "terraform version" to check we've downloaded a terraform 0.12.26 binary
	cmd := exec.Command(tfpath, "version")

	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	expected := "Terraform v0.13.0-beta1"
	actual := string(out)
	if !strings.HasPrefix(actual, expected) {
		t.Fatalf("ran terraform version, expected %s, but got %s", expected, actual)
	}
}

// test that Find returns an appropriate error when given an exact path
// which exists, but is not a terraform executable
func TestExactPath(t *testing.T) {
	// we just want the path to a local executable that definitely exists
	execPath, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Find(ExactPath(execPath))
	if err == nil {
		t.Fatalf("expected Find() to fail when given ExactPath(%s), but it did not", execPath)
	}

	expected := fmt.Sprintf("executable found at path %s is not terraform", execPath)
	if !strings.HasPrefix(err.Error(), expected) {
		t.Fatalf("expected Find() to return %s, but got %s", expected, err)
	}
}

// test that Find falls back to the next working strategy when the file at
// ExactPath does not exist
func TestExactVersionFallback(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "tfinstall-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tfpath, err := Find(ExactPath("/hopefully/completely/nonexistent/path"), ExactVersion("0.12.26", tmpDir))
	if err != nil {
		t.Fatal(err)
	}

	// run "terraform version" to check we've downloaded a terraform 0.12.26 binary
	cmd := exec.Command(tfpath, "version")

	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	expected := "Terraform v0.12.26"
	actual := string(out)
	if !strings.HasPrefix(actual, expected) {
		t.Fatalf("ran terraform version, expected %s, but got %s", expected, actual)
	}
}

// latest version calculation itself is handled by checkpoint, so the test can be straightforward -
// just test that we've managed to download a version of terraform later than 0.12.27
func TestLatestVersion(t *testing.T) {
	lowerBound := "0.12.27"

	tmpDir, err := ioutil.TempDir("", "tfinstall-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tfpath, err := Find(LatestVersion(tmpDir, false))
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(tfpath, "version")

	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	lowerBoundVersion, err := version.NewVersion("0.12.27")
	if err != nil {
		t.Fatal(err)
	}

	outVersion := strings.Trim(string(out), "\n")
	outVersion = strings.TrimLeft(outVersion, "Terraform v")

	actualVersion, err := version.NewVersion(outVersion)
	if err != nil {
		t.Fatal(err)
	}

	if actualVersion.LessThan(lowerBoundVersion) {
		t.Fatalf("ran terraform version, expected version to be greater than %s, but got %s", lowerBound, out)
	}
}
