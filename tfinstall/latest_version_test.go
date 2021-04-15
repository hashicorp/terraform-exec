package tfinstall

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
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

	tfpath, err := Find(context.Background(), LatestVersion(tmpDir, false))
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(tfpath, "version", "-json")

	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	lowerBoundVersion, err := version.NewVersion("0.15.0")
	if err != nil {
		t.Fatal(err)
	}

	type versionOutput struct {
		TerraformVersion string `json:"terraform_version"`
	}
	vOut := versionOutput{}
	err = json.Unmarshal(out, &vOut)
	if err != nil {
		t.Fatal(err)
	}

	actualVersion, err := version.NewVersion(vOut.TerraformVersion)
	if err != nil {
		t.Fatal(err)
	}

	if actualVersion.LessThan(lowerBoundVersion) {
		t.Fatalf("ran terraform version, expected version to be greater than %s, but got %s", lowerBound, out)
	}
}
