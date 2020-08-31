package e2etest

import (
	"context"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	jsonEncodeIntMinVersion = version.Must(version.NewVersion("0.12.0"))
)

func TestConsole(t *testing.T) {
	runTest(t, "empty_with_tf_file", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if runtime.GOOS == "windows" {
			t.Skip("terraform console does not support windows currently: https://github.com/hashicorp/terraform/issues/18242")
		}

		ctx := context.Background()
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		out, err := tf.Console(ctx, "1+5")
		if err != nil {
			t.Fatal(err)
		}

		out = strings.TrimSpace(out)

		if out != "6" {
			t.Fatalf("expected 6, got %q", out)
		}
	})
}

func consoleJSON(t *testing.T, tf *tfexec.Terraform, expr string, v interface{}) {
	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := tf.ConsoleJSON(ctx, expr, v)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConsoleJSON(t *testing.T) {
	runTest(t, "empty_with_tf_file", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if runtime.GOOS == "windows" {
			t.Skip("terraform console does not support windows currently: https://github.com/hashicorp/terraform/issues/18242")
		}

		t.Run("int", func(t *testing.T) {
			if tfv.LessThan(jsonEncodeIntMinVersion) {
				t.Skip("jsonencode does not work on ints in 0.11 for some reason")
			}
			var out int
			consoleJSON(t, tf, `6`, &out)
			if out != 6 {
				t.Fatalf("unexpected value, got %d", out)
			}
		})

		t.Run("int expression", func(t *testing.T) {
			if tfv.LessThan(jsonEncodeIntMinVersion) {
				t.Skip("jsonencode does not work on ints in 0.11 for some reason")
			}
			var out int
			consoleJSON(t, tf, `5+1`, &out)
			if out != 6 {
				t.Fatalf("unexpected value, got %d", out)
			}
		})

		t.Run("string", func(t *testing.T) {
			var out string
			consoleJSON(t, tf, `"this is a string"`, &out)
			if out != "this is a string" {
				t.Fatalf("unexpected value, got %q", out)
			}
		})

		t.Run("string expression", func(t *testing.T) {
			var out string
			consoleJSON(t, tf, `join(" ", split(" ", "this is a string"))`, &out)
			if out != "this is a string" {
				t.Fatalf("unexpected value, got %q", out)
			}
		})

		t.Run("array of strings", func(t *testing.T) {
			var out []string
			consoleJSON(t, tf, `split(" ", "this is a string")`, &out)
			if !reflect.DeepEqual([]string{"this", "is", "a", "string"}, out) {
				t.Fatalf("unexpected value, got %#v", out)
			}
		})

		// TODO: error for multiline expressions? HEREDOCs?
	})
}
