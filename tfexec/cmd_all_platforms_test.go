package tfexec

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"
	"time"
)

func Test_runTerraformCmd(t *testing.T) {
	// Checks runTerraformCmd for race condition when using
	// go test -race -run Test_runTerraformCmd_default ./tfexec
	var buf bytes.Buffer

	tf := &Terraform{
		logger:   log.New(&buf, "", 0),
		execPath: "echo",
	}

	ctx, cancel := context.WithCancel(context.Background())

	cmd := tf.buildTerraformCmd(ctx, nil, "hello tf-exec!")
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

func Test_runTerraformCmdCancel(t *testing.T) {
	var buf bytes.Buffer

	tf := &Terraform{
		logger:   log.New(&buf, "", 0),
		execPath: "sleep",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := tf.buildTerraformCmd(ctx, nil, "10")
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	err := tf.runTerraformCmd(ctx, cmd)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %T %s", err, err)
	}
}
