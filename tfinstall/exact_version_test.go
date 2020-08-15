package tfinstall

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
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
