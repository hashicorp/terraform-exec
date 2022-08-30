package tfexec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

const defaultGracefulShutdownTimeout = 0

func (tf *Terraform) runTerraformCmd(parentCtx context.Context, cmd *exec.Cmd) error {
	return tf.runTerraformCmdWithGracefulshutdownTimeout(parentCtx, cmd, defaultGracefulShutdownTimeout)
}

func (tf *Terraform) runTerraformCmdWithGracefulshutdownTimeout(parentCtx context.Context, cmd *exec.Cmd, gracefulShutdownTimeout time.Duration) error {
	var errBuf strings.Builder

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// kill children if parent is dead
		Pdeathsig: syscall.SIGKILL,
		// set process group ID
		Setpgid: true,
	}

	// check for early cancellation
	select {
	case <-parentCtx.Done():
		return parentCtx.Err()
	default:
	}
	// Context for the stdout and stderr writers so that they are not closed on parentCtx cancellation
	// and avoiding "broken pipe"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	cmdMu := sync.Mutex{}
	withCmdLock := func(fn func() error) error {
		cmdMu.Lock()
		defer cmdMu.Unlock()
		return fn()
	}
	cmdDoneCh := make(chan error, 1)
	returnCh := make(chan error, 1)
	defer close(returnCh)
	go func() {
		select {
		case <-parentCtx.Done(): // wait for context cancelled
			tf.logger.Printf("[WARN] The context was cancelled, we'll let Terraform finish by sending SIGINT signal")
			if err := withCmdLock(func() error { return cmd.Process.Signal(os.Interrupt) }); err != nil {
				tf.logger.Printf("[ERROR] Error sending SIGINT to terraform: %v", err)
			}
			// give some time to the process before forcefully killing it
			select {
			case <-time.After(gracefulShutdownTimeout):
				// Forcefully kill the process
				if err := withCmdLock(func() error { return cmd.Process.Signal(os.Kill) }); err != nil {
					tf.logger.Printf("[ERROR] Error sending SIGKILL to terraform: %v", err)
				}
				cancel() // to cancel stdout/stderr writers
				tf.logger.Printf("[ERROR] terraform forcefully killed after graceful shutdown timeout")
				returnCh <- fmt.Errorf("%w: terraform forcefully killed after graceful shutdown timeout", parentCtx.Err())
			case err := <-cmdDoneCh:
				returnCh <- fmt.Errorf("%w: %v", parentCtx.Err(), err)
				tf.logger.Printf("[INFO] terraform successfully interrupted")
			}
		case err := <-cmdDoneCh:
			tf.logger.Printf("[DEBUG] terraform run finished")
			returnCh <- err
		}
	}()
	err = withCmdLock(func() error { return cmd.Start() })
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

	if err == nil && parentCtx.Err() != nil {
		err = parentCtx.Err()
	}
	if err != nil {
		return tf.wrapExitError(parentCtx, err, errBuf.String())
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
