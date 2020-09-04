package tfexec

import (
	"fmt"
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
		"TF_LOG_PATH=",
		"TF_LOG=",
		"TF_WORKSPACE=",
	}
}

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
