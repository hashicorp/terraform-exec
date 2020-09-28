package tfexec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestContext_alreadyTimedOut(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Microsecond)
	t.Cleanup(cancelFunc)

	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.13.1"))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = tf.Version(ctx, true)
	if err == nil {
		t.Fatal("expected error from version command")
	}

	ee := &ExitError{}
	isExitErr := errors.As(err, &ee)
	if !isExitErr {
		t.Fatalf("expected error to be ExitError compatible, given: %#v", err)
	}

	dee := context.DeadlineExceeded
	isCtxErr := errors.Is(err, dee)
	if !isCtxErr {
		t.Fatalf("expected error to be context.DeadlineExceeded compatible, given: %#v", err)
	}
}

func TestContext_timeout(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Millisecond)
	t.Cleanup(cancelFunc)

	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.13.1"))
	if err != nil {
		t.Fatal(err)
	}
	err = tf.sleep(ctx, 1*time.Second)
	if err == nil {
		t.Fatal("expected error from timed out sleep")
	}

	ee := &ExitError{}
	isExitErr := errors.As(err, &ee)
	if !isExitErr {
		t.Fatalf("expected error to be ExitError compatible, given: %#v", err)
	}

	isCtxErr := errors.Is(ee.ctxErr, context.DeadlineExceeded)
	if !isCtxErr {
		t.Fatalf("expected context error to be context.DeadlineExceeded")
	}
}

func (tf *Terraform) sleep(ctx context.Context, d time.Duration) error {
	seconds := fmt.Sprintf("%.0f", d.Seconds())
	env := map[string]string{}
	cmd := tf.buildTerraformCmd(ctx, env, seconds)

	sPath, err := exec.LookPath("sleep")
	if err != nil {
		return err
	}
	cmd.Path = sPath
	cmd.Args[0] = sPath

	return tf.runTerraformCmd(ctx, cmd)
}

func TestContext_alreadyCancelled(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()

	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.13.1"))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = tf.Version(ctx, true)
	if err == nil {
		t.Fatal("expected error from version command")
	}

	ee := &ExitError{}
	isExitErr := errors.As(err, &ee)
	if !isExitErr {
		t.Fatalf("expected error to be ExitError compatible, given: %#v", err)
	}

	dee := context.Canceled
	isCtxErr := errors.Is(err, dee)
	if !isCtxErr {
		t.Fatalf("expected error to be context.Canceled compatible, given: %#v", err)
	}
}
