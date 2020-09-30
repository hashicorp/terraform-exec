package tfexec

import (
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"
)

func TestMock_sleep(t *testing.T) {
	tf := mockTerraform(t, &MockCall{
		Args:          []string{"version"},
		SleepDuration: 5 * time.Millisecond,
		Stdout:        "Terraform v0.12.0\n",
	})

	timeout := 1 * time.Millisecond
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancelFunc)

	expectedErr := &exec.ExitError{}
	_, _, err := tf.Version(ctx, true)
	if err != nil {
		if errors.As(err, &expectedErr) {
			return
		}

		t.Fatalf("errors don't match.\nexpected: %#v\ngiven:    %#v\n",
			expectedErr, err)
	}

	t.Fatalf("expected timeout error: %#v", expectedErr)
}

func TestMock_singleCall(t *testing.T) {
	tf := mockTerraform(t, &MockCall{
		Args:     []string{"version"},
		Stdout:   "Terraform v0.12.0\n",
		ExitCode: 0,
	})
	v, _, err := tf.Version(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "0.12.0" {
		t.Fatalf("output does not match: %#v", v)
	}
}

func TestMock_multipleCalls(t *testing.T) {
	expectedOutput := "formatted config"
	tf := mockTerraform(t, &MockQueue{
		Q: []*MockItem{
			{
				Args:     []string{"version"},
				Stdout:   "Terraform v0.13.1",
				ExitCode: 0,
			},
			{
				Args:     []string{"fmt", "-no-color", "-"},
				Stdout:   string(expectedOutput),
				ExitCode: 0,
			},
		},
	})
	out, err := tf.FormatString(context.Background(), "unformatted")
	if err != nil {
		t.Fatal(err)
	}

	if out != expectedOutput {
		t.Fatalf("Expected output: %q\nGiven: %q\n",
			string(expectedOutput), string(out))
	}
}

func TestMock_jsonOutput(t *testing.T) {
	tf := mockTerraform(t, &MockCall{
		Args:     []string{"providers", "schema", "-json", "-no-color"},
		Stdout:   `{"format_version": "0.1"}`,
		ExitCode: 0,
	})

	ps, err := tf.ProvidersSchema(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	expectedVersion := "0.1"
	if ps.FormatVersion != expectedVersion {
		t.Fatalf("format version doesn't match.\nexpected: %q\ngiven: %q\n",
			expectedVersion, ps.FormatVersion)
	}
}

func mockTerraform(t *testing.T, md MockItemDispenser) *Terraform {
	tf, err := NewMockTerraform(md)
	if err != nil {
		t.Fatal(err)
	}
	return tf
}
