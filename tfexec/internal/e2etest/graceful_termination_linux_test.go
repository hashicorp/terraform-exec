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

func Test_gracefulTerminationRunTerraformCmd_linux(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest_v1_1}, "infinite_loop", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		var bufStdout bytes.Buffer
		var bufStderr bytes.Buffer
		tf.SetStderr(&bufStdout)
		tf.SetStdout(&bufStderr)

		ctx, cancel := context.WithCancel(context.Background())
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}
		doneCh := make(chan error)
		go func() {
			doneCh <- tf.Apply(ctx)
		}()
		time.Sleep(3 * time.Second)
		cancel()
		err = <-doneCh
		close(doneCh)
		if err != nil {
			t.Log(err)
		}
		output := bufStderr.String() + bufStdout.String()
		if !strings.Contains(output, "Gracefully shutting down...") {
			t.Log(output)
			t.Fatal("canceling context should gracefully shut terraform down")
		}
	})

}
