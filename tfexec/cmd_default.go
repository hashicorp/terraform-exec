//go:build !linux
// +build !linux

package tfexec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const defaultGracefulShutdownTimeout = 0

func (tf *Terraform) runTerraformCmd(parentCtx context.Context, cmd *exec.Cmd) error {
	return tf.runTerraformCmdWithGracefulshutdownTimeout(parentCtx, cmd, defaultGracefulShutdownTimeout)
}

func (tf *Terraform) runTerraformCmdWithGracefulshutdownTimeout(ctx context.Context, cmd *exec.Cmd, _ time.Duration) error {
	var errBuf strings.Builder

	// check for early cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Read stdout / stderr logs from pipe instead of setting cmd.Stdout and
	// cmd.Stderr because it can cause hanging when killing the command
	// https://github.com/golang/go/issues/23019
	stdoutWriter := mergeWriters(cmd.Stdout, tf.stdout)
	stderrWriter := mergeWriters(tf.stderr, &errBuf)

	cmd.Stderr = nil
	cmd.Stdout = nil

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmdDoneCh := make(chan error, 1)
	returnCh := make(chan error, 1)
	defer close(cmdDoneCh)
	go func() {
		defer close(returnCh)
		select {
		case <-ctx.Done(): // wait for context cancelled
			cmd.Process.Signal(os.Kill)
			returnCh <- fmt.Errorf("%w: terraform forcefully killed", ctx.Err())
		case err := <-cmdDoneCh:
			returnCh <- err
		}
	}()
	err = cmd.Start()
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		return tf.wrapExitError(ctx, err, "")
	}

	var errStdout, errStderr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		errStdout = writeOutput(ctx, stdoutPipe, stdoutWriter)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errStderr = writeOutput(ctx, stderrPipe, stderrWriter)
	}()

	// Reads from pipes must be completed before calling cmd.Wait(). Otherwise
	// can cause a race condition
	wg.Wait()

	cmdDoneCh <- cmd.Wait()
	err = <-returnCh

	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		return tf.wrapExitError(ctx, err, errBuf.String())
	}

	// Return error if there was an issue reading the std out/err
	if errStdout != nil && ctx.Err() != nil {
		return tf.wrapExitError(ctx, errStdout, errBuf.String())
	}
	if errStderr != nil && ctx.Err() != nil {
		return tf.wrapExitError(ctx, errStderr, errBuf.String())
	}

	return nil
}
