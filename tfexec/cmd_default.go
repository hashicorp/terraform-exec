//go:build !linux
// +build !linux

package tfexec

import (
	"context"
	"os/exec"
	"strings"
	"sync"
)

func (tf *Terraform) runTerraformCmd(ctx context.Context, cmd *exec.Cmd) error {
	var errBuf strings.Builder
	// ensure we don't mix up stdout and stderr
	sw := &syncWriter{w: &errBuf}
	cmd.Stdout = mergeWriters(cmd.Stdout, tf.stdout, sw)
	cmd.Stderr = mergeWriters(cmd.Stderr, tf.stderr, sw)

	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded || ctx.Err() == context.Canceled {
			if cmd != nil && cmd.Process != nil && cmd.ProcessState != nil {
				err := cmd.Process.Kill()
				if err != nil {
					tf.logger.Printf("error from kill: %s", err)
				}
			}
		}
	}()

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

	err = cmd.Wait()
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
