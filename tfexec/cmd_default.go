//go:build !linux
// +build !linux

package tfexec

import (
	"context"
	"os/exec"
	"strings"
)

func (tf *Terraform) runTerraformCmd(ctx context.Context, cmd *exec.Cmd) error {
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

	err = cmd.Start()
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		return tf.wrapExitError(ctx, err, "")
	}

	exitChLen := 2
	exitCh := make(chan error, exitChLen)
	go func() {
		exitCh <- writeOutput(stdoutPipe, stdoutWriter)
	}()
	go func() {
		exitCh <- writeOutput(stderrPipe, stderrWriter)
	}()

	err = cmd.Wait()
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		return tf.wrapExitError(ctx, err, errBuf.String())
	}

	// Wait for the logs to finish writing
	counter := 0
	for {
		counter++
		err := <-exitCh
		if err != nil && err != context.Canceled {
			return tf.wrapExitError(ctx, err, errBuf.String())
		}
		if counter >= exitChLen {
			return ctx.Err()
		}
	}
}
