// list taken from https://github.com/golang/go/blob/91ef076562dfcf783074dbd84ad7c6db60fdd481/src/go/build/syslist.go#L38-L51
//go:build aix || android || darwin || dragonfly || freebsd || hurd || illumos || ios || linux || netbsd || openbsd || solaris
// +build aix android darwin dragonfly freebsd hurd illumos ios linux netbsd openbsd solaris

package e2etest

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func Test_gracefulTerminationRunTerraformCmd(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest_v1_1}, "infinite_loop", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		var bufStdout bytes.Buffer
		var bufStderr bytes.Buffer
		tf.SetStderr(&bufStdout)
		tf.SetStdout(&bufStderr)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := tf.Init(ctx)
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		doneCh := make(chan error)
		shutdown := make(chan struct{})
		go func() {
			doneCh <- tf.Apply(ctx, tfexec.InterruptChannel(shutdown))
		}()

		time.Sleep(3 * time.Second)
		close(shutdown)
		err = <-doneCh
		close(doneCh)
		if err != nil {
			t.Log(err)
		}
		output := bufStderr.String() + bufStdout.String()
		t.Log(output)
		if !strings.Contains(output, "Gracefully shutting down...") {
			t.Fatal("canceling context should gracefully shut terraform down")
		}
	})
}

func Test_gracefulTerminationRunTerraformCmdWithNoGracefulShutdownTimeout(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest_v1_1}, "infinite_loop", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		var bufStdout bytes.Buffer
		var bufStderr bytes.Buffer
		tf.SetStderr(&bufStdout)
		tf.SetStdout(&bufStderr)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := tf.Init(ctx)
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		doneCh := make(chan error)
		go func() {
			doneCh <- tf.Apply(ctx, tfexec.InterruptChannel(make(chan struct{})))
		}()

		time.Sleep(3 * time.Second)
		cancel()
		err = <-doneCh
		close(doneCh)
		if err != nil {
			t.Log(err)
		}
		output := bufStderr.String() + bufStdout.String()
		t.Log(output)
		if strings.Contains(output, "Gracefully shutting down...") {
			t.Fatal("canceling context with no graceful shutdown timeout should immediately kill the process and not start a graceful cancellation")
		}
	})
}
