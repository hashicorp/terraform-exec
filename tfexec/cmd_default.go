//go:build !linux
// +build !linux

package tfexec

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func (tf *Terraform) runTerraformCmd(ctx context.Context, cmd *exec.Cmd) error {
	var errBuf strings.Builder

	cmd.Stdout = mergeWriters(cmd.Stdout, tf.stdout)
	cmd.Stderr = mergeWriters(cmd.Stderr, tf.stderr, &errBuf)

	go func() {
		<-ctx.Done()
		fmt.Println("context done", ctx.Err(), cmd.Process, cmd.ProcessState)
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

	err := cmd.Run()
	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		return tf.wrapExitError(ctx, err, errBuf.String())
	}

	return nil
}
