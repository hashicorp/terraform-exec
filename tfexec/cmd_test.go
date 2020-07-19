package tfexec

import (
	"os/exec"
	"strings"
	"testing"
)

var defaultEnv = []string{
	"TF_LOG=",
	"TF_LOG_PATH=",
	"TF_IN_AUTOMATION=1",
	"CHECKPOINT_DISABLE=",
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
	expectedEnv = envMap(append(defaultEnv, envSlice(expectedEnv)...))

	// compare against raw slice len incase of duplication or something
	if len(expectedEnv) != len(actual.Env) {
		t.Fatalf("env mismatch\n\nexpected:\n%v\n\ngot:\n%v", envSlice(expectedEnv), actual.Env)
	}

	actualEnv := envMap(actual.Env)

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
