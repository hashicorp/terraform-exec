package e2etest

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func TestFormatString(t *testing.T) {
	runTest(t, "", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		unformatted := strings.TrimSpace(`
resource     "foo"      "bar" {
	baz = 1
		qux      =        2
}
`)

		expected := strings.TrimSpace(`
resource "foo" "bar" {
  baz = 1
  qux = 2
}
`)

		actual, err := tf.FormatString(context.Background(), unformatted)
		if err != nil {
			t.Fatal(err)
		}

		actual = strings.TrimSpace(actual)

		if actual != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s\n", expected, actual)
		}
	})
}

func TestFormatCheck(t *testing.T) {
	runTest(t, "unformatted", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		checksums := map[string]uint32{
			"file1.tf": checkSum(t, filepath.Join(tf.WorkingDir(), "file1.tf")),
			"file2.tf": checkSum(t, filepath.Join(tf.WorkingDir(), "file2.tf")),
		}

		formatted, files, err := tf.FormatCheck(context.Background())
		if err != nil {
			t.Fatalf("error from FormatCheck: %T %q", err, err)
		}

		if formatted {
			t.Fatal("expected unformatted")
		}

		if !reflect.DeepEqual(files, []string{"file1.tf", "file2.tf"}) {
			t.Fatalf("unexpected files list: %#v", files)
		}

		for file, checksum := range checksums {
			if checksum != checkSum(t, filepath.Join(tf.WorkingDir(), file)) {
				t.Fatalf("%s should not have changed", file)
			}
		}
	})
}

func TestFormatWrite(t *testing.T) {
	runTest(t, "unformatted", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.FormatWrite(context.Background())
		if err != nil {
			t.Fatalf("error from FormatWrite: %T %q", err, err)
		}

		for file, golden := range map[string]string{
			"file1.tf": "file1.golden.txt",
			"file2.tf": "file2.golden.txt",
		} {
			textFilesEqual(t, filepath.Join(tf.WorkingDir(), golden), filepath.Join(tf.WorkingDir(), file))
		}
	})
}

func TestFormat(t *testing.T) {
	runTest(t, "", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		unformatted := strings.TrimSpace(`
resource     "foo"      "bar" {
	baz = 1
		qux      =        2
}
`)

		expected := strings.TrimSpace(`
resource "foo" "bar" {
  baz = 1
  qux = 2
}
`)

		start := time.Now()
		var actual strings.Builder
		err := tf.Format(context.Background(), strings.NewReader(unformatted), &actual)
		if err != nil {
			t.Fatal(err)
		}
		duration := time.Since(start)
		t.Logf("formatting took %dms", duration.Milliseconds())

		actualString := strings.TrimSpace(actual.String())
		if actualString != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s\n", expected, actualString)
		}
	})
}

func TestFormat_warmFormatter(t *testing.T) {
	runTest(t, "", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		unformatted := strings.TrimSpace(`
resource     "foo"      "bar" {
	baz = 1
		qux      =        2
}
`)

		expected := strings.TrimSpace(`
resource "foo" "bar" {
  baz = 1
  qux = 2
}
`)

		inR, inW := io.Pipe()
		outR, outW := io.Pipe()

		go func() {
			err := tf.Format(context.Background(), inR, outW)
			if err != nil {
				outW.CloseWithError(err)
			}
			_ = outW.Close()
		}()

		t.Log("Sleeping while CLI is warmed...")
		time.Sleep(5 * time.Second)
		t.Log("Sending unformatted data...")
		start := time.Now()
		_, err := inW.Write([]byte(unformatted))
		if err != nil {
			t.Fatal(err)
		}
		err = inW.Close()
		if err != nil {
			t.Fatal(err)
		}

		actual, err := ioutil.ReadAll(outR)
		if err != nil {
			t.Fatal(err)
		}
		duration := time.Since(start)
		t.Logf("formatting took %dms", duration.Milliseconds())

		actualString := strings.TrimSpace(string(actual))
		if actualString != expected {
			t.Fatalf("expected:\n%s\ngot:\n%s\n", expected, actualString)
		}
	})
}
