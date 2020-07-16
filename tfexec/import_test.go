package tfexec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImport(t *testing.T) {
	const (
		expectedID      = "asdlfjksdlfkjsdlfk"
		resourceAddress = "random_string.random_string"
	)
	ctx := context.Background()

	for _, tfv := range []string{
		"0.11.14", // doesn't support show JSON output, but does support import
		"0.12.28",
		"0.13.0-beta3",
	} {
		t.Run(tfv, func(t *testing.T) {
			td := testTempDir(t)
			defer os.RemoveAll(td)

			err := copyFiles(filepath.Join(testFixtureDir, "import"), td)
			if err != nil {
				t.Fatal(err)
			}

			tf, err := NewTerraform(td, tfVersion(t, tfv))
			if err != nil {
				t.Fatal(err)
			}

			err = tf.Init(ctx, Lock(false))
			if err != nil {
				t.Fatal(err)
			}

			// Config is unnecessary here since its already the working dir, but just testing an additional flag
			err = tf.Import(ctx, resourceAddress, expectedID, DisableBackup(), Lock(false), Config(td))
			if err != nil {
				t.Fatal(err)
			}

			if strings.HasPrefix(tfv, "0.11.") {
				t.Logf("skipping state assertion for 0.11")
				return
			}

			state, err := tf.Show(ctx)
			if err != nil {
				t.Fatal(err)
			}

			for _, r := range state.Values.RootModule.Resources {
				if r.Address != resourceAddress {
					continue
				}

				raw, ok := r.AttributeValues["id"]
				if !ok {
					t.Fatal("value not found for \"id\" attribute")
				}
				actual, ok := raw.(string)
				if !ok {
					t.Fatalf("unable to cast %T to string: %#v", raw, raw)
				}

				if actual != expectedID {
					t.Fatalf("expected %q, got %q", expectedID, actual)
				}

				// success
				return
			}

			t.Fatalf("imported resource %q not found", resourceAddress)
		})
	}
}

func TestImportCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	importCmd := tf.importCmd(context.Background(), "my-addr", "my-id")

	actual := strings.TrimPrefix(cmdString(importCmd), importCmd.Path+" ")

	expected := "import -no-color -input=false -lock-timeout=0s -lock=true my-addr my-id"

	if actual != expected {
		t.Fatalf("expected default arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	importCmd = tf.importCmd(context.Background(), "my-addr2", "my-id2",
		Backup("testbackup"),
		LockTimeout("200s"),
		State("teststate"),
		StateOut("teststateout"),
		VarFile("testvarfile"),
		Lock(false),
		Var("var1=foo"),
		Var("var2=bar"),
		AllowMissingConfig(true),
	)

	actual = strings.TrimPrefix(cmdString(importCmd), importCmd.Path+" ")

	expected = "import -no-color -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -allow-missing-config -var 'var1=foo' -var 'var2=bar' my-addr2 my-id2"

	if actual != expected {
		t.Fatalf("expected arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
