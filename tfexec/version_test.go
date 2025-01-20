// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func mustVersion(t *testing.T, s string) *version.Version {
	v, err := version.NewVersion(s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestParsePlaintextVersionOutput(t *testing.T) {
	for i, c := range []struct {
		expectedV         *version.Version
		expectedProviders map[string]*version.Version

		stdout string
	}{
		// 0.13 tests
		{
			mustVersion(t, "0.13.0-dev"), nil, `
Terraform v0.13.0-dev`,
		},
		{
			mustVersion(t, "0.13.0-dev"), map[string]*version.Version{
				"registry.terraform.io/hashicorp/null": mustVersion(t, "2.1.2"),
				"registry.terraform.io/paultyng/null":  mustVersion(t, "0.1.0"),
			}, `
Terraform v0.13.0-dev
+ provider registry.terraform.io/hashicorp/null v2.1.2
+ provider registry.terraform.io/paultyng/null v0.1.0`,
		},
		{
			mustVersion(t, "0.13.0-dev"), nil, `
Terraform v0.13.0-dev

Your version of Terraform is out of date! The latest version
is 0.13.1. You can update by downloading from https://www.terraform.io/downloads.html`,
		},
		{
			mustVersion(t, "0.13.0-dev"), map[string]*version.Version{
				"registry.terraform.io/hashicorp/null": mustVersion(t, "2.1.2"),
				"registry.terraform.io/paultyng/null":  mustVersion(t, "0.1.0"),
			}, `
Terraform v0.13.0-dev
+ provider registry.terraform.io/hashicorp/null v2.1.2
+ provider registry.terraform.io/paultyng/null v0.1.0

Your version of Terraform is out of date! The latest version
is 0.13.1. You can update by downloading from https://www.terraform.io/downloads.html`,
		},

		// 0.12 tests
		{
			mustVersion(t, "0.12.26"), nil, `
Terraform v0.12.26
`,
		},
		{
			mustVersion(t, "0.12.26"), map[string]*version.Version{
				"null": mustVersion(t, "2.1.2"),
			}, `
Terraform v0.12.26
+ provider.null v2.1.2
`,
		},
		{
			mustVersion(t, "0.12.18"), nil, `
Terraform v0.12.18

Your version of Terraform is out of date! The latest version
is 0.12.26. You can update by downloading from https://www.terraform.io/downloads.html
`,
		},
		{
			mustVersion(t, "0.12.18"), map[string]*version.Version{
				"null": mustVersion(t, "2.1.2"),
			}, `
Terraform v0.12.18
+ provider.null v2.1.2

Your version of Terraform is out of date! The latest version
is 0.12.26. You can update by downloading from https://www.terraform.io/downloads.html
`,
		},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expectedV), func(t *testing.T) {
			actualV, actualProv, err := parsePlaintextVersionOutput(c.stdout)
			if err != nil {
				t.Fatal(err)
			}

			if !c.expectedV.Equal(actualV) {
				t.Fatalf("expected %s, got %s", c.expectedV, actualV)
			}

			for k, v := range c.expectedProviders {
				if actual := actualProv[k]; actual == nil || !v.Equal(actual) {
					t.Fatalf("expected %s for %s, got %s", v, k, actual)
				}
			}

			if len(c.expectedProviders) != len(actualProv) {
				t.Fatalf("expected %d providers, got %d", len(c.expectedProviders), len(actualProv))
			}
		})
	}
}

func TestParseJsonVersionOutput(t *testing.T) {
	testStdout := []byte(`{
  "terraform_version": "0.15.0-beta1",
  "platform": "darwin_amd64",
  "provider_selections": {
    "registry.terraform.io/hashicorp/aws": "3.31.0",
    "registry.terraform.io/hashicorp/google": "3.58.0"
  },
  "terraform_outdated": false
}
`)
	tfVersion, pvs, err := parseJsonVersionOutput(testStdout)
	if err != nil {
		t.Fatal(err)
	}
	expectedTfVer := mustVersion(t, "0.15.0-beta1")

	if !expectedTfVer.Equal(tfVersion) {
		t.Fatalf("version doesn't match (%q != %q)",
			expectedTfVer.String(), tfVersion.String())
	}

	expectedPvs := map[string]*version.Version{
		"registry.terraform.io/hashicorp/aws":    mustVersion(t, "3.31.0"),
		"registry.terraform.io/hashicorp/google": mustVersion(t, "3.58.0"),
	}
	if diff := cmp.Diff(expectedPvs, pvs); diff != "" {
		t.Fatalf("provider versions don't match: %s", diff)
	}
}

func TestVersionInRange(t *testing.T) {
	for i, c := range []struct {
		expected bool
		min      string
		tfv      string
		max      string
	}{
		{true, "", "0.12.26", ""},
		{true, "", "0.13.0-beta3", ""},

		{false, "", "0.12.26", "0.12.25"},
		{false, "", "0.12.26", "0.12.26"},
		{false, "0.12.27", "0.12.26", ""},
		{true, "", "0.12.26", "0.13.0"},
		{true, "0.12.25", "0.12.26", ""},
		{true, "0.12.26", "0.12.26", ""},
		{true, "0.12.26", "0.12.26", "0.12.27"},
		{true, "0.12.26", "0.12.26", "0.13.0"},

		{false, "0.12.26", "0.13.0-beta3", "0.13.0"},
		{true, "0.12.26", "0.13.0-beta3", ""},
		{true, "0.13.0", "0.13.0-beta3", ""},
		{true, "0.13.0", "0.13.0-beta3", "0.14.0"},
		{true, "", "0.13.0-beta3", "0.14.0"},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tfv, err := version.NewVersion(c.tfv)
			if err != nil {
				t.Fatal(err)
			}

			var min *version.Version
			if c.min != "" {
				min, err = version.NewVersion(c.min)
				if err != nil {
					t.Fatal(err)
				}
			}

			var max *version.Version
			if c.max != "" {
				max, err = version.NewVersion(c.max)
				if err != nil {
					t.Fatal(err)
				}
			}

			actual := versionInRange(tfv, min, max)
			if actual != c.expected {
				t.Fatalf("expected %v, got %v: %s <= %s < %s", c.expected, actual, min, tfv, max)
			}
		})
	}
}

func TestCompatible(t *testing.T) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		t.Skip("Terraform for darwin/arm64 is not available until v1")
	}

	ev := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("0.12.26")),
	}
	ev.SetLogger(testutil.TestLogger())

	ctx := context.Background()
	t.Cleanup(func() { ev.Remove(ctx) })

	tf01226, err := ev.Install(ctx)
	if err != nil {
		t.Fatal(err)
	}

	ev = &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("0.13.0-beta3")),
	}
	ev.SetLogger(testutil.TestLogger())
	t.Cleanup(func() { ev.Remove(ctx) })

	tf013beta3, err := ev.Install(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for i, c := range []struct {
		expected bool
		min      string
		max      string
		binPath  string
	}{
		{false, "0.12.27", "", tf01226},
		{false, "0.12.26", "0.13.0", tf013beta3},

		{true, "0.12.25", "", tf01226},
		{true, "0.12.26", "0.13.0", tf01226},
		{true, "", "0.12.27", tf01226},

		{true, "0.12.26", "", tf013beta3},
		{true, "0.13.0", "", tf013beta3},
		{true, "0.13.0", "0.14.0", tf013beta3},
		{true, "", "0.14.0", tf013beta3},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			tf, err := NewTerraform(filepath.Dir(c.binPath), c.binPath)
			if err != nil {
				t.Fatal(err)
			}

			var min *version.Version
			if c.min != "" {
				min, err = version.NewVersion(c.min)
				if err != nil {
					t.Fatal(err)
				}
			}

			var max *version.Version
			if c.max != "" {
				max, err = version.NewVersion(c.max)
				if err != nil {
					t.Fatal(err)
				}
			}
			var mismatch *ErrVersionMismatch
			err = tf.compatible(context.Background(), min, max)
			switch {
			case c.expected && err != nil:
				t.Fatal(err)
			case !c.expected && err == nil:
				t.Fatal("expected version mismatch error, no error returned")
			case !c.expected && !errors.As(err, &mismatch):
				t.Fatal(err)
			}
		})
	}
}

func TestExperimentsEnabled(t *testing.T) {
	testCases := map[string]struct {
		tfVersion     *version.Version
		expectedError error
	}{
		"experiments-enabled-in-alphas": {
			tfVersion: version.Must(version.NewVersion(testutil.Latest_Alpha_v1_9)),
		},
		"experiments-disabled-in-betas": {
			tfVersion:     version.Must(version.NewVersion(testutil.Latest_Beta_v1_8)),
			expectedError: fmt.Errorf("experiments are not enabled in version %s, as it's not an alpha or dev build", testutil.Latest_Beta_v1_8),
		},
		"experiments-disabled-in-stable": {
			tfVersion:     version.Must(version.NewVersion(testutil.Latest_v1_5)),
			expectedError: fmt.Errorf("experiments are not enabled in version %s, as it's not an alpha or dev build", testutil.Latest_v1_5),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			ev := &releases.ExactVersion{
				Product: product.Terraform,
				Version: testCase.tfVersion,
			}
			ev.SetLogger(testutil.TestLogger())

			ctx := context.Background()
			t.Cleanup(func() { ev.Remove(ctx) })

			tfBinPath, err := ev.Install(ctx)
			if err != nil {
				t.Fatal(err)
			}

			tf, err := NewTerraform(filepath.Dir(tfBinPath), tfBinPath)
			if err != nil {
				t.Fatal(err)
			}

			err = tf.experimentsEnabled(context.Background())
			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("expected no error, got: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("got no error, expected: %s", testCase.expectedError)
			}
		})
	}
}
