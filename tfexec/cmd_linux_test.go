package tfexec

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
	"time"
)

func Test_runTerraformCmd_linux(t *testing.T) {
	// Checks runTerraformCmd for race condition when using
	// go test -race -run Test_runTerraformCmd_linux ./tfexec -tags=linux
	var buf bytes.Buffer

	tf := &Terraform{
		logger:   log.New(&buf, "", 0),
		execPath: "echo",
	}

	ctx, cancel := context.WithCancel(context.Background())

	cmd := tf.buildTerraformCmd(nil, "hello tf-exec!")
	err := tf.runTerraformCmd(ctx, cmd)
	if err != nil {
		t.Fatal(err)
	}

	// Cancel stops the leaked go routine which logs an error
	cancel()
	time.Sleep(time.Second)
	if strings.Contains(buf.String(), "error from kill") {
		t.Fatal("canceling context should not lead to logging an error")
	}
}
