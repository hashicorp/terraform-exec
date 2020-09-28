package tfexec

import (
	"context"
	"errors"
	"os"
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
