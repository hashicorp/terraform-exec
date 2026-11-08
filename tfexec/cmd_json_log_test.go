// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRunTerraformCmdJSONLogLargeMessage(t *testing.T) {
	if os.Getenv("TFEXEC_LARGE_LOG_HELPER") == "1" {
		data, err := json.Marshal(map[string]string{"@level": "info", "@message": strings.Repeat("x", 2*1024*1024), "type": "log"})
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		// Also exercise a final line without a newline.
		fmt.Print(`{"@level":"info","@message":"last","type":"log"}`)
		os.Exit(0)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	tf, err := NewTerraform(t.TempDir(), executable)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, "-test.run=^TestRunTerraformCmdJSONLogLargeMessage$")
	cmd.Env = append(os.Environ(), "TFEXEC_LARGE_LOG_HELPER=1")
	finished := make(chan error, 1)
	go func() {
		var messages []string
		for next := range tf.runTerraformCmdJSONLog(ctx, cmd) {
			if next.Err != nil {
				finished <- next.Err
				return
			}
			if next.Msg != nil {
				messages = append(messages, next.Msg.Message())
			}
		}
		if len(messages) != 2 || messages[0] != strings.Repeat("x", 2*1024*1024) || messages[1] != "last" {
			finished <- fmt.Errorf("expected large message followed by final unterminated message; got %d messages", len(messages))
			return
		}
		finished <- nil
	}()
	select {
	case err := <-finished:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("JSON log reader blocked on a large valid message")
	}
}
