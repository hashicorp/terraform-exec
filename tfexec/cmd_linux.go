package tfexec

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// syncWriter is an io.Writer protected by a sync.Mutex.
type syncWriter struct {
	sync.Mutex
	w io.Writer
}

// Write implements io.Writer.
func (w *syncWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	return w.w.Write(p)
}

func (tf *Terraform) runTerraformCmd(ctx context.Context, cmd *exec.Cmd) error {
	var errBuf strings.Builder
	// ensure we don't mix up stdout and stderr
	sw := &syncWriter{w: &errBuf}
	cmd.Stdout = mergeWriters(cmd.Stdout, tf.stdout, sw)
	cmd.Stderr = mergeWriters(cmd.Stderr, tf.stderr, sw)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// kill children if parent is dead
		Pdeathsig: syscall.SIGKILL,
		// set process group ID
		Setpgid: true,
	}

	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded || ctx.Err() == context.Canceled {
			if cmd != nil && cmd.Process != nil && cmd.ProcessState != nil {
				// send SIGINT to process group
				err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
				if err != nil {
					tf.logger.Printf("error from SIGINT: %s", err)
				}
			}

			// TODO: send a kill if it doesn't respond for a bit?
		}
	}()

	// check for early cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err := cmd.Run()
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		return tf.wrapExitError(ctx, err, errBuf.String())
	}

	return nil
}
