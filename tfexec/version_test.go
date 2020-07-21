package tfexec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
)

func TestParseVersionOutput(t *testing.T) {
	var mustVer = func(s string) *version.Version {
		v, err := version.NewVersion(s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}

	for i, c := range []struct {
		expectedV         *version.Version
		expectedProviders map[string]*version.Version

		stdout string
	}{
		// 0.13 tests
		{
			mustVer("0.13.0-dev"), nil, `
Terraform v0.13.0-dev`,
		},
		{
			mustVer("0.13.0-dev"), map[string]*version.Version{
				"registry.terraform.io/hashicorp/null": mustVer("2.1.2"),
				"registry.terraform.io/paultyng/null":  mustVer("0.1.0"),
			}, `
Terraform v0.13.0-dev
+ provider registry.terraform.io/hashicorp/null v2.1.2
+ provider registry.terraform.io/paultyng/null v0.1.0`,
		},
		{
			mustVer("0.13.0-dev"), nil, `
Terraform v0.13.0-dev

Your version of Terraform is out of date! The latest version
is 0.13.1. You can update by downloading from https://www.terraform.io/downloads.html`,
		},
		{
			mustVer("0.13.0-dev"), map[string]*version.Version{
				"registry.terraform.io/hashicorp/null": mustVer("2.1.2"),
				"registry.terraform.io/paultyng/null":  mustVer("0.1.0"),
			}, `
Terraform v0.13.0-dev
+ provider registry.terraform.io/hashicorp/null v2.1.2
+ provider registry.terraform.io/paultyng/null v0.1.0

Your version of Terraform is out of date! The latest version
is 0.13.1. You can update by downloading from https://www.terraform.io/downloads.html`,
		},

		// 0.12 tests
		{
			mustVer("0.12.26"), nil, `
Terraform v0.12.26
`,
		},
		{
			mustVer("0.12.26"), map[string]*version.Version{
				"null": mustVer("2.1.2"),
			}, `
Terraform v0.12.26
+ provider.null v2.1.2
`,
		},
		{
			mustVer("0.12.18"), nil, `
Terraform v0.12.18

Your version of Terraform is out of date! The latest version
is 0.12.26. You can update by downloading from https://www.terraform.io/downloads.html
`,
		},
		{
			mustVer("0.12.18"), map[string]*version.Version{
				"null": mustVer("2.1.2"),
			}, `
Terraform v0.12.18
+ provider.null v2.1.2

Your version of Terraform is out of date! The latest version
is 0.12.26. You can update by downloading from https://www.terraform.io/downloads.html
`,
		},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expectedV), func(t *testing.T) {
			actualV, actualProv, err := parseVersionOutput(c.stdout)
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
