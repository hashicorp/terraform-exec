package tfexec

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

func (tf *Terraform) runTerraformCmd(parentCtx context.Context, cmd *exec.Cmd) error {
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
	// TODO: select
	cmdDoneCh := make(chan error)
	returnCh := make(chan error)
	defer close(returnCh)
	go func() {
		select {
		case <-parentCtx.Done(): // wait for context cancelled
			tf.logger.Printf("[WARN] The context was cancelled, we'll let Terraform finish by sending SIGINT signal")
			cmd.Process.Signal(os.Interrupt)
			if err != nil {
				tf.logger.Printf("[WARN] Error sending SIGINT to terraform: %v", err)
			}
			// give 10 seconds to the process before force killing it
			// TODO: make it configurable
			select {
			case <-time.After(10 * time.Second):
				cmd.Process.Signal(os.Kill) // to kill the process
				cancel()                    // to cancel stdout/stderr writers
				tf.logger.Printf("[ERROR] terraform forcefully killed after graceful shutdown timeout")
				returnCh <- errors.New("terraform forcefully killed after graceful shutdown timeout")
			case err := <-cmdDoneCh:
				returnCh <- err
				tf.logger.Printf("[INFO] terraform successfully interrupted")
			}
		case err := <-cmdDoneCh:
			tf.logger.Printf("[DEBUG] terraform run finished")
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
