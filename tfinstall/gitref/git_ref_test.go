package gitref_test

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/hashicorp/terraform-exec/tfinstall/gitref"
)

// ensure the option satisfies the interface
var _ tfinstall.ExecPathFinder = &gitref.Option{}

func TestGitRef(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping git ref tests for short run")
	}

	cmd := exec.Command("go", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error with Go installation: %s\n%s", err, string(out))
	}
	t.Logf("go version\n%s", string(out))

	for n, c := range map[string]struct {
		expectedVersion string
		ref             string
	}{
		"branch v0.12": {"Terraform v0.12.", "refs/heads/v0.12"},
		"tag v0.12.29": {"Terraform v0.12.29", "refs/tags/v0.12.29"},
		// "commit 83630a7": {"Terraform v0.12.29", "83630a7003fb8b868a3bf940798326634c3c6acc"},
		"empty": {"Terraform v1.", ""}, // should pull main, which is currently v1 dev
	} {
		c := c
		t.Run(n, func(t *testing.T) {
			// these are really long running due to the compilation, run them in parallel
			t.Parallel()

			ctx := context.Background()

			// hacking this tmpdir to local dir due to circle perms, should be env var
			tmpDir, err := ioutil.TempDir("", "tfinstall-test")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				os.RemoveAll(tmpDir)
			})

			t.Logf("finding / building ref %q...", c.ref)
			tfpath, err := tfinstall.Find(ctx, gitref.Install(c.ref, "", tmpDir))
			if err != nil {
				t.Fatalf("%T %s", err, err)
			}

			t.Logf("testing version cmd...")
			cmd := exec.Command(tfpath, "version")

			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("%s\n\n%s", err, out)
			}

			actual := string(out)
			if !strings.Contains(actual, c.expectedVersion) {
				t.Fatalf("ran terraform version, expected %s, but got %s", c.expectedVersion, actual)
			}
		})
	}
}
