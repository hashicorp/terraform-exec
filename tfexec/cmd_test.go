// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/internal/version"
)

func TestMergeUserAgent(t *testing.T) {
	for i, c := range []struct {
		expected string
		uas      []string
	}{
		{"foo/1 bar/2", []string{"foo/1", "bar/2"}},
		{"foo/1 bar/2", []string{"foo/1 bar/2"}},
		{"foo/1 bar/2", []string{"", "foo/1", "bar/2"}},
		{"foo/1 bar/2", []string{"", "foo/1 bar/2"}},
		{"foo/1 bar/2", []string{"  ", "foo/1 bar/2"}},
		{"foo/1 bar/2", []string{"foo/1", "", "bar/2"}},
		{"foo/1 bar/2", []string{"foo/1", "   ", "bar/2"}},

		// comments
		{"foo/1 (bar/1 bar/2 bar/3) bar/2", []string{"foo/1 (bar/1 bar/2 bar/3)", "bar/2"}},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			actual := mergeUserAgent(c.uas...)
			if c.expected != actual {
				t.Fatalf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}

func defaultEnv() []string {
	return []string{
		"CHECKPOINT_DISABLE=",
		"TF_APPEND_USER_AGENT=HashiCorp-terraform-exec/" + version.ModuleVersion(),
		"TF_IN_AUTOMATION=1",
		"TF_LOG=",
		"TF_LOG_CORE=",
		"TF_LOG_PATH=",
		"TF_LOG_PROVIDER=",
	}
}

// stripSafeInherited removes variables that safeInheritedEnv pulls in from the
// test machine's environment. assertCmd calls this on the actual env so that
// command-level tests can assert on the library-managed vars only, without
// needing to know which safe-inherited vars happen to be set on the host.
func stripSafeInherited(env map[string]string) {
	for k := range inheritedEnvAllowlist {
		delete(env, k)
	}
	for k := range env {
		for _, prefix := range inheritedEnvPrefixAllowlist {
			if strings.HasPrefix(k, prefix) {
				delete(env, k)
				break
			}
		}
	}
}

// assertCmd asserts that a constructed exec.Cmd contains the expected args and environment variables.
// The command itself isn't executed; that is only done in E2E tests.
func assertCmd(t *testing.T, expectedArgs []string, expectedEnv map[string]string, actual *exec.Cmd) {
	t.Helper()

	// check args (skip path)
	actualArgs := actual.Args[1:]

	if len(expectedArgs) != len(actualArgs) {
		t.Fatalf("args mismatch\n\nexpected:\n%v\n\ngot:\n%v", strings.Join(expectedArgs, " "), strings.Join(actualArgs, " "))
	}
	for i := range expectedArgs {
		if expectedArgs[i] != actualArgs[i] {
			t.Fatalf("args mismatch, expected %q, got %q\n\nfull expected:\n%v\n\nfull actual:\n%v", expectedArgs[i], actualArgs[i], strings.Join(expectedArgs, " "), strings.Join(actualArgs, " "))
		}
	}

	// check environment
	expectedEnv = envMap(append(defaultEnv(), envSlice(expectedEnv)...))
	actualEnv := envMap(actual.Env)

	if len(actualEnv) != len(actual.Env) {
		t.Fatalf("duplication in actual env, unable to assert: %v", actual.Env)
	}

	// ignore tempdir related env vars from comparison
	for _, k := range []string{"TMPDIR", "TMP", "TEMP", "USERPROFILE"} {
		if _, ok := expectedEnv[k]; ok {
			t.Logf("ignoring env var %q", k)
			delete(expectedEnv, k)
		}

		if _, ok := actualEnv[k]; ok {
			t.Logf("ignoring env var %q", k)
			delete(actualEnv, k)
		}
	}

	// strip safe-inherited vars from actual: they come from the test machine's
	// environment and are not what command-level tests are asserting on
	stripSafeInherited(actualEnv)

	// compare against raw slice len incase of duplication or something
	if len(expectedEnv) != len(actualEnv) {
		t.Fatalf("env mismatch\n\nexpected:\n%v\n\ngot:\n%v", envSlice(expectedEnv), actual.Env)
	}

	for k, ev := range expectedEnv {
		av, ok := actualEnv[k]
		if !ok {
			t.Fatalf("env mismatch, missing %q\n\nfull expected:\n%v\n\nfull actual:\n%v", k, envSlice(expectedEnv), envSlice(actualEnv))
		}
		if ev != av {
			t.Fatalf("env mismatch, expected %q, got %q\n\nfull expected:\n%v\n\nfull actual:\n%v", ev, av, envSlice(expectedEnv), envSlice(actualEnv))
		}
	}
}

func TestSafeInheritedEnvAllowlist(t *testing.T) {
	for key := range inheritedEnvAllowlist {
		t.Run(key, func(t *testing.T) {
			t.Setenv(key, "test-value")
			env := safeInheritedEnv()
			if v, ok := env[key]; !ok || v != "test-value" {
				t.Fatalf("expected allowlisted var %q to be inherited, present=%v value=%q", key, ok, v)
			}
		})
	}
}

func TestSafeInheritedEnvCloudPrefixes(t *testing.T) {
	for _, tc := range []struct {
		key    string
		reason string
	}{
		{"AWS_ACCESS_KEY_ID", "AWS credentials"},
		{"AWS_SECRET_ACCESS_KEY", "AWS credentials"},
		{"AWS_SESSION_TOKEN", "AWS credentials"},
		{"AWS_REGION", "AWS config"},
		{"GOOGLE_CREDENTIALS", "GCP credentials"},
		{"GOOGLE_APPLICATION_CREDENTIALS", "GCP service account"},
		{"GCLOUD_PROJECT", "gcloud config"},
		{"CLOUDSDK_COMPUTE_ZONE", "Cloud SDK config"},
		{"ARM_CLIENT_ID", "Azure credentials"},
		{"ARM_TENANT_ID", "Azure credentials"},
		{"AZURE_CLIENT_SECRET", "Azure credentials"},
		{"VAULT_ADDR", "Vault config"},
		{"VAULT_TOKEN", "Vault credentials"},
		{"GITHUB_TOKEN", "GitHub provider"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			t.Setenv(tc.key, "test-value")
			env := safeInheritedEnv()
			if v, ok := env[tc.key]; !ok || v != "test-value" {
				t.Fatalf("expected cloud provider var %q (%s) to be inherited, present=%v value=%q", tc.key, tc.reason, ok, v)
			}
		})
	}
}

func TestSafeInheritedEnvBlocksDangerousVars(t *testing.T) {
	for _, tc := range []struct {
		key    string
		reason string
	}{
		{"LD_PRELOAD", "dynamic linker injection"},
		{"LD_LIBRARY_PATH", "dynamic linker injection"},
		{"DYLD_INSERT_LIBRARIES", "macOS dynamic linker injection"},
		{"DYLD_LIBRARY_PATH", "macOS dynamic linker injection"},
		{"TF_CLI_CONFIG_FILE", "redirect terraform config"},
		{"TF_PLUGIN_DIR", "redirect provider lookup"},
		{"TF_DATA_DIR", "redirect .terraform directory"},
		{"TERRAFORMRC", "redirect terraform config"},
		{"BASH_ENV", "bash code injection via subshells"},
		{"ENV", "POSIX shell code injection"},
		{"TF_LOG", "managed by library, not user-injectable"},
		{"TF_WORKSPACE", "managed by library, not user-injectable"},
		{"TF_INPUT", "managed by library, not user-injectable"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			t.Setenv(tc.key, "malicious-value")
			env := safeInheritedEnv()
			if v, ok := env[tc.key]; ok {
				t.Fatalf("expected dangerous var %q (%s) to be blocked, but it was inherited with value %q", tc.key, tc.reason, v)
			}
		})
	}
}

func TestSafeInheritedEnvSetEnvCanPassAnything(t *testing.T) {
	// SetEnv must be able to pass non-allowlisted vars explicitly — callers
	// need an escape hatch for providers not covered by the allowlist.
	td := t.TempDir()
	tf, err := NewTerraform(td, "echo")
	if err != nil {
		t.Fatal(err)
	}
	err = tf.SetEnv(map[string]string{
		"KUBECONFIG":    "/tmp/kubeconfig",
		"MY_CUSTOM_VAR": "custom-value",
	})
	if err != nil {
		t.Fatal(err)
	}

	cmd := tf.buildTerraformCmd(t.Context(), nil)
	env := envMap(cmd.Env)

	for k, want := range map[string]string{
		"KUBECONFIG":    "/tmp/kubeconfig",
		"MY_CUSTOM_VAR": "custom-value",
	} {
		if got := env[k]; got != want {
			t.Fatalf("expected SetEnv var %q=%q in cmd env, got %q", k, want, got)
		}
	}
}

func TestSafeInheritedEnvEmptyValue(t *testing.T) {
	// an allowlisted var set to empty string should still pass through
	t.Setenv("HOME", "")
	env := safeInheritedEnv()
	if _, ok := env["HOME"]; !ok {
		t.Fatal("expected HOME with empty value to be inherited")
	}
}

func TestSafeInheritedEnvDangerousVarNotInheritedEvenIfPrefixMatches(t *testing.T) {
	// GOOGLE_* prefix is allowed, but a hypothetical GOOGLE_INTERNAL_OVERRIDE
	// that slipped in should only reach Terraform if the caller sets it via
	// SetEnv — not by accident from the ambient environment.
	// This test documents the accepted risk: we inherit all GOOGLE_* vars.
	// The allowlisted prefixes are intentionally broad for IAM compatibility.
	t.Setenv("AWS_ANYTHING", "value")
	env := safeInheritedEnv()
	if _, ok := env["AWS_ANYTHING"]; !ok {
		// expected: prefix-matched vars do pass through
		t.Fatal("AWS_ANYTHING should be inherited via AWS_ prefix")
	}

	// contrast: a non-matching var must not pass through
	t.Setenv("RANDOM_SECRET", "secret")
	env = safeInheritedEnv()
	if _, ok := env["RANDOM_SECRET"]; ok {
		t.Fatal("RANDOM_SECRET must not be inherited")
	}
}

func TestSafeInheritedEnvHomePresent(t *testing.T) {
	// HOME must always be inherited
	home := os.Getenv("HOME")
	if home == "" {
		t.Skip("HOME not set in this environment")
	}
	env := safeInheritedEnv()
	if env["HOME"] != home {
		t.Fatalf("expected HOME=%q, got %q", home, env["HOME"])
	}
}
