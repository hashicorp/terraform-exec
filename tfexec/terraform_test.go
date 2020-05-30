package tfexec

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

const testFixtureDir = "testdata"
const testConfigFileName = "main.tf"
const testStateJsonFileName = "state.json"
const testTerraformStateFileName = "terraform.tfstate"

func TestInitCmd(t *testing.T) {
	tf, err := NewTerraform("/dev/null", "")
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	initCmd := tf.InitCmd(context.Background())

	actual := strings.TrimPrefix(initCmd.String(), initCmd.Path+" ")

	expected := "init -no-color -force-copy -input=false -lock-timeout=0s -backend=true -get=true -get-plugins=true -lock=true -upgrade=false -verify-plugins=true"

	if actual != expected {
		t.Fatalf("expected default arguments of InitCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	initCmd = tf.InitCmd(context.Background(), Backend(false), BackendConfig("confpath1"), BackendConfig("confpath2"), FromModule("testsource"), Get(false), GetPlugins(false), Lock(false), LockTimeout("999s"), PluginDir("testdir1"), PluginDir("testdir2"), Reconfigure(true), Upgrade(true), VerifyPlugins(false))

	actual = strings.TrimPrefix(initCmd.String(), initCmd.Path+" ")

	expected = "init -no-color -force-copy -input=false -from-module=testsource -lock-timeout=999s -backend=false -get=false -get-plugins=false -lock=false -upgrade=true -verify-plugins=false -reconfigure -backend-config=confpath1 -backend-config=confpath2 -plugin-dir=testdir1 -plugin-dir=testdir2"

	if actual != expected {
		t.Fatalf("expected arguments of InitCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}

func TestPlanCmd(t *testing.T) {
	tf, err := NewTerraform("/dev/null", "")
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	planCmd := tf.PlanCmd(context.Background())

	actual := strings.TrimPrefix(planCmd.String(), planCmd.Path+" ")

	expected := "plan -no-color -input=false -lock-timeout=0s -lock=true -parallelism=10 -refresh=true"

	if actual != expected {
		t.Fatalf("expected default arguments of PlanCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	planCmd = tf.PlanCmd(context.Background(), Destroy(true), Lock(false), LockTimeout("22s"), Out("whale"), Parallelism(42), Refresh(false), State("marvin"), Target("zaphod"), Target("beeblebrox"), Var("android=paranoid"), Var("brain_size=planet"), VarFile("trillian"))

	actual = strings.TrimPrefix(planCmd.String(), planCmd.Path+" ")

	expected = "plan -no-color -input=false -lock-timeout=22s -out=whale -state=marvin -var-file=trillian -lock=false -parallelism=42 -refresh=false -destroy -target=zaphod -target=beeblebrox -var 'android=paranoid' -var 'brain_size=planet'"

	if actual != expected {
		t.Fatalf("expected arguments of PlanCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}

func TestDestroyCmd(t *testing.T) {
	tf, err := NewTerraform("/dev/null", "")
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	destroyCmd := tf.DestroyCmd(context.Background())

	actual := strings.TrimPrefix(destroyCmd.String(), destroyCmd.Path+" ")

	expected := "destroy -no-color -auto-approve -lock-timeout=0s -lock=true -parallelism=10 -refresh=true"

	if actual != expected {
		t.Fatalf("expected default arguments of DestroyCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	destroyCmd = tf.DestroyCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"))

	actual = strings.TrimPrefix(destroyCmd.String(), destroyCmd.Path+" ")

	expected = "destroy -no-color -auto-approve -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -parallelism=99 -refresh=false -target=target1 -target=target2 -var 'var1=foo' -var 'var2=bar'"

	if actual != expected {
		t.Fatalf("expected arguments of DestroyCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}

func TestImportCmd(t *testing.T) {
	tf, err := NewTerraform("/dev/null", "")
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	importCmd := tf.ImportCmd(context.Background())

	actual := strings.TrimPrefix(importCmd.String(), importCmd.Path+" ")

	expected := "import -no-color -input=false -lock-timeout=0s -lock=true"

	if actual != expected {
		t.Fatalf("expected default arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	importCmd = tf.ImportCmd(context.Background(),
		Backup("testbackup"),
		LockTimeout("200s"),
		State("teststate"),
		StateOut("teststateout"),
		VarFile("testvarfile"),
		Lock(false),
		Var("var1=foo"),
		Var("var2=bar"),
		AllowMissingConfig(true))

	actual = strings.TrimPrefix(importCmd.String(), importCmd.Path+" ")

	expected = "import -no-color -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -allow-missing-config -var 'var1=foo' -var 'var2=bar'"

	if actual != expected {
		t.Fatalf("expected arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}

func TestOutputCmd(t *testing.T) {
	tf, err := NewTerraform("/dev/null", "")
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	outputCmd := tf.OutputCmd(context.Background())

	actual := strings.TrimPrefix(outputCmd.String(), outputCmd.Path+" ")

	expected := "output -no-color -json"

	if actual != expected {
		t.Fatalf("expected default arguments of OutputCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	outputCmd = tf.OutputCmd(context.Background(),
		State("teststate"))

	actual = strings.TrimPrefix(outputCmd.String(), outputCmd.Path+" ")

	expected = "output -no-color -json -state=teststate"

	if actual != expected {
		t.Fatalf("expected arguments of ImportCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}

func TestShow(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, "")
	if err != nil {
		t.Fatal(err)
	}

	// copy state and config files into test dir
	err = copyFile(filepath.Join(testFixtureDir, testTerraformStateFileName), filepath.Join(td, testTerraformStateFileName))
	if err != nil {
		t.Fatalf("error copying state file into test dir: %s", err)
	}
	err = copyFile(filepath.Join(testFixtureDir, testConfigFileName), filepath.Join(td, testConfigFileName))
	if err != nil {
		t.Fatalf("error copying config file into test dir: %s", err)
	}

	expected := tfjson.State{
		FormatVersion:    "0.1",
		TerraformVersion: "0.12.24",
		Values: &tfjson.StateValues{
			RootModule: &tfjson.StateModule{
				Resources: []*tfjson.StateResource{&tfjson.StateResource{
					Address: "null_resource.foo",
					AttributeValues: map[string]interface{}{
						"id":       "5510719323588825107",
						"triggers": nil,
					},
					Mode:         tfjson.ManagedResourceMode,
					Type:         "null_resource",
					Name:         "foo",
					ProviderName: "null",
				}},
			},
		},
	}

	err = tf.Init(context.Background())
	if err != nil {
		t.Fatalf("error running Init in test directory: %s", err)
	}

	actual, err := tf.Show(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, &expected) {
		t.Fatalf("actual: %s\nexpected: %s", spew.Sdump(actual), spew.Sdump(expected))
	}
}

func TestShow_errInitRequired(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, "")
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, testTerraformStateFileName), filepath.Join(td, testTerraformStateFileName))

	// This test will break if the error output of `terraform init`
	// changes significantly. We tolerate this brittleness as a poor
	// man's canary for significant changes to Terraform CLI.
	// FIXME: Parse this in the actual command and return ErrNoInit
	expected := "Error: Could not satisfy plugin requirements"

	_, err = tf.Show(context.Background())
	if err == nil {
		t.Fatal("expected Show to error, but it did not")
	} else {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("expected error %s to contain %s", err, expected)
		}
	}

}

func TestApply(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, "")
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(filepath.Join(testFixtureDir, testConfigFileName), filepath.Join(td, testConfigFileName))
	if err != nil {
		t.Fatalf("error copying config file into test dir: %s", err)
	}

	err = tf.Init(context.Background())
	if err != nil {
		t.Fatalf("error running Init in test directory: %s", err)
	}

	err = tf.Apply(context.Background())
	if err != nil {
		t.Fatalf("error running Apply: %s", err)
	}
}

func TestApplyCmd(t *testing.T) {
	tf, err := NewTerraform("/dev/null", "")
	if err != nil {
		t.Fatal(err)
	}

	applyCmd := tf.ApplyCmd(context.Background(), Backup("testbackup"), LockTimeout("200s"), State("teststate"), StateOut("teststateout"), VarFile("testvarfile"), Lock(false), Parallelism(99), Refresh(false), Target("target1"), Target("target2"), Var("var1=foo"), Var("var2=bar"), DirOrPlan("testfile"))

	actual := strings.TrimPrefix(applyCmd.String(), applyCmd.Path+" ")

	expected := "apply -no-color -auto-approve -input=false -backup=testbackup -lock-timeout=200s -state=teststate -state-out=teststateout -var-file=testvarfile -lock=false -parallelism=99 -refresh=false -target=target1 -target=target2 -var 'var1=foo' -var 'var2=bar' testfile"

	if actual != expected {
		t.Fatalf("expected arguments of ApplyCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}

func testTempDir(t *testing.T) string {
	d, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("error creating temporary test directory: %s", err)
	}

	return d
}

func copyFile(path string, dstPath string) error {
	srcF, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcF.Close()

	dstF, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}
