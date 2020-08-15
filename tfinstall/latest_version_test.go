package tfinstall

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
)

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
