// This file contains tests that only compile/work in Go 1.13 and forward
// +build go1.13

package e2etest

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	protocol5MinVersion = version.Must(version.NewVersion("0.12.0"))
)

func TestUnparsedError(t *testing.T) {
	// This simulates an unparsed error from the Cmd.Run method (in this case file not found). This
	// is to ensure we don't miss raising unexpected errors in addition to parsed / well known ones.
	runTest(t, "", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {

		// force delete the working dir to cause an os.PathError
		err := os.RemoveAll(tf.WorkingDir())
		if err != nil {
			t.Fatal(err)
		}

		err = tf.Init(context.Background())
		if err == nil {
			t.Fatalf("expected error running Init, none returned")
		}
		var e *os.PathError
		if !errors.As(err, &e) {
			t.Fatalf("expected os.PathError, got %T, %s", err, err)
		}
	})
}

func TestMissingVar(t *testing.T) {
	runTest(t, "var", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("err during init: %s", err)
		}

		// Variable names from testdata/var/main.tf
		shortVarName := "no_default"
		longVarName := "no_default_really_long_variable_name_that_will_line_wrap_tf_output"

		// Test for ErrMissingVar and properly formatted error message on shorter variable names
		_, err = tf.Plan(context.Background(), tfexec.Var(longVarName+"=foo"))
		if err == nil {
			t.Fatalf("expected error running Plan, none returned")
		}
		var e *tfexec.ErrMissingVar
		if !errors.As(err, &e) {
			t.Fatalf("expected ErrMissingVar, got %T, %s", err, err)
		}

		if e.VariableName != shortVarName {
			t.Fatalf("expected missing %s, got %q", shortVarName, e.VariableName)
		}

		// Test for ErrMissingVar and properly formatted error message on long variable names
		_, err = tf.Plan(context.Background(), tfexec.Var(shortVarName+"=foo"))
		if err == nil {
			t.Fatalf("expected error running Plan, none returned")
		}
		if !errors.As(err, &e) {
			t.Fatalf("expected ErrMissingVar, got %T, %s", err, err)
		}

		if e.VariableName != longVarName {
			t.Fatalf("expected missing %s, got %q", longVarName, e.VariableName)
		}

		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("expected exec.ExitError, got %T, %s", err, err)
		}

		// Test for no error when all variables have a value
		_, err = tf.Plan(context.Background(), tfexec.Var(shortVarName+"=foo"), tfexec.Var(longVarName+"=foo"))
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	})
}

func TestTFVersionMismatch(t *testing.T) {
	runTest(t, "tf99", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		// force cache version for error messaging
		_, _, err := tf.Version(context.Background(), true)
		if err != nil {
			t.Fatal(err)
		}

		err = tf.Init(context.Background())
		if err == nil {
			t.Fatal("expected error, but didn't find one")
		}

		var e *tfexec.ErrTFVersionMismatch
		if !errors.As(err, &e) {
			t.Fatalf("expected ErrTFVersionMismatch, got %T, %s", err, err)
		}

		// in 0.12, we just return "unknown" as the specifics are not included in the error messaging
		if e.Constraint != "unknown" && e.Constraint != ">99.0.0" {
			t.Fatalf("unexpected constraint %q", e.Constraint)
		}

		if e.TFVersion != tfv.String() {
			t.Fatalf("expected %q, got %q", tfv.String(), e.TFVersion)
		}

		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("expected exec.ExitError, got %T, %s", err, err)
		}
	})
}

func TestContext_alreadyPastDeadline(t *testing.T) {
	runTest(t, "", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
		defer cancel()

		_, _, err := tf.Version(ctx, true)
		if err == nil {
			t.Fatal("expected error from version command")
		}

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.DeadlineExceeded, got %T %s", err, err)
		}
	})
}

func TestContext_sleepNoCancellation(t *testing.T) {
	// this test is just to ensure that time_sleep works properly without cancellation
	runTest(t, "sleep", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		// only testing versions that can cancel mid apply
		if !tfv.GreaterThanOrEqual(protocol5MinVersion) {
			t.Skip("the ability to interrupt an apply was added in protocol 5.0 in Terraform 0.12, so test is not valid")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("err during init: %s", err)
		}

		ctx := context.Background()
		start := time.Now()
		err = tf.Apply(ctx, tfexec.Var(`create_duration=5s`))
		if err != nil {
			t.Fatalf("error during apply: %s", err)
		}
		elapsed := time.Now().Sub(start)
		if elapsed < 5*time.Second {
			t.Fatalf("expected runtime of at least 5s, got %s", elapsed)
		}
	})
}

func TestContext_sleepTimeoutExpired(t *testing.T) {
	runTest(t, "sleep", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		// only testing versions that can cancel mid apply
		if !tfv.GreaterThanOrEqual(protocol5MinVersion) {
			t.Skip("the ability to interrupt an apply was added in protocol 5.0 in Terraform 0.12, so test is not valid")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("err during init: %s", err)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = tf.Apply(ctx)
		if err == nil {
			t.Fatal("expected error, but didn't find one")
		}

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.DeadlineExceeded, got %T %s", err, err)
		}
	})
}

func TestContext_alreadyCancelled(t *testing.T) {
	runTest(t, "", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err := tf.Version(ctx, true)
		if err == nil {
			t.Fatal("expected error from version command")
		}

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %T %s", err, err)
		}
	})
}
