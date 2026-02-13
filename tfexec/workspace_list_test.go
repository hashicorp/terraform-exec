// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestWorkspaceListCmd(t *testing.T) {
	tf, err := NewTerraform(t.TempDir(), tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		workspaceListCmd, err := tf.workspaceListCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "list",
			"-no-color",
		}, nil, workspaceListCmd)
	})

	t.Run("reattach config", func(t *testing.T) {
		workspaceListCmd, err := tf.workspaceListCmd(context.Background(), Reattach(map[string]ReattachConfig{
			"registry.terraform.io/hashicorp/examplecloud": {
				Protocol:        "grpc",
				ProtocolVersion: 6,
				Pid:             1234,
				Test:            true,
				Addr: ReattachConfigAddr{
					Network: "unix",
					String:  "/fake_folder/T/plugin123",
				},
			},
		}))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"workspace", "list",
			"-no-color",
		}, map[string]string{
			"TF_REATTACH_PROVIDERS": `{"registry.terraform.io/hashicorp/examplecloud":{"Protocol":"grpc","ProtocolVersion":6,"Pid":1234,"Test":true,"Addr":{"Network":"unix","String":"/fake_folder/T/plugin123"}}}`,
		}, workspaceListCmd)
	})
}

func TestParseWorkspaceList(t *testing.T) {
	for i, c := range []struct {
		expected        []string
		expectedCurrent string
		stdout          string
	}{
		{
			[]string{"default"},
			"default",
			`* default

`,
		},
		{
			[]string{"default", "foo", "bar"},
			"foo",
			`  default
* foo
  bar

`,
		},

		// linux new lines
		{
			[]string{"default", "foo"},
			"foo",
			"  default\n* foo\n\n",
		},
		// windows new lines
		{
			[]string{"default", "foo"},
			"foo",
			"  default\r\n* foo\r\n\r\n",
		},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actualList, actualCurrent := parseWorkspaceList(c.stdout)

			if actualCurrent != c.expectedCurrent {
				t.Fatalf("expected selected %q, got %q", c.expectedCurrent, actualCurrent)
			}

			if !reflect.DeepEqual(c.expected, actualList) {
				t.Fatalf("expected %#v, got %#v", c.expected, actualList)
			}
		})
	}
}
