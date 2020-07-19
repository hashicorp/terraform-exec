package tfexec

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestApply(t *testing.T) {
	ctx := context.Background()

	for _, c := range []struct {
		version     string
		configDir   string
		checkOutput bool
	}{
		{"0.11.14", "basic", false},
		{"0.12.28", "basic", false},
		{"0.13.0-beta3", "basic", false},

		{"0.12.28", "var", true},
		{"0.13.0-beta3", "var", true},
	} {
		testName := fmt.Sprintf(fmt.Sprintf("%s %s", c.version, c.configDir))
		t.Run(testName, func(t *testing.T) {
			td := testTempDir(t)
			defer os.RemoveAll(td)

			err := copyFiles(filepath.Join(testFixtureDir, c.configDir), td)
			if err != nil {
				t.Fatal(err)
			}

			tf, err := NewTerraform(td, tfVersion(t, c.version))
			if err != nil {
				t.Fatal(err)
			}

			err = tf.Init(ctx, Lock(false))
			if err != nil {
				t.Fatal(err)
			}

			opts := []ApplyOption{}
			if c.checkOutput {
				opts = append(opts, Var("in="+testName))
			}

			err = tf.Apply(ctx, opts...)
			if err != nil {
				t.Fatalf("error running Apply: %s", err)
			}

			outputs, err := tf.Output(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if !c.checkOutput {
				return
			}

			if out, ok := outputs["out"]; ok {
				var vs string
				err = json.Unmarshal(out.Value, &vs)
				if err != nil {
					t.Fatal(err)
				}

				if vs != testName {
					t.Fatalf("expected %q, got %q", testName, vs)
				}

				return
			}

			t.Fatalf("output %q not found", "out")
		})
	}
}

func TestApplyCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, "0.12.28"))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("basic", func(t *testing.T) {
		applyCmd := tf.applyCmd(context.Background(),
			Backup("testbackup"),
			LockTimeout("200s"),
			State("teststate"),
			StateOut("teststateout"),
			VarFile("testvarfile"),
			Lock(false),
			Parallelism(99),
			Refresh(false),
			Target("target1"),
			Target("target2"),
			Var("var1=foo"),
			Var("var2=bar"),
			DirOrPlan("testfile"),
		)

		assertCmd(t, []string{
			"apply",
			"-no-color",
			"-auto-approve",
			"-input=false",
			"-backup=testbackup",
			"-lock-timeout=200s",
			"-state=teststate",
			"-state-out=teststateout",
			"-var-file=testvarfile",
			"-lock=false",
			"-parallelism=99",
			"-refresh=false",
			"-target=target1",
			"-target=target2",
			"-var", "var1=foo",
			"-var", "var2=bar",
			"testfile",
		}, nil, applyCmd)
	})
}
