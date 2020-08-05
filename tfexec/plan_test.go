package tfexec

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestPlanCmd(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest012))
	if err != nil {
		t.Fatal(err)
	}

	// defaults
	planCmd := tf.planCmd(context.Background())

	actual := strings.TrimPrefix(cmdString(planCmd), planCmd.Path+" ")

	expected := "plan -no-color -input=false -lock-timeout=0s -lock=true -parallelism=10 -refresh=true"

	if actual != expected {
		t.Fatalf("expected default arguments of PlanCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}

	// override all defaults
	planCmd = tf.planCmd(context.Background(), Destroy(true), Lock(false), LockTimeout("22s"), Out("whale"), Parallelism(42), Refresh(false), State("marvin"), Target("zaphod"), Target("beeblebrox"), Var("android=paranoid"), Var("brain_size=planet"), VarFile("trillian"))

	actual = strings.TrimPrefix(cmdString(planCmd), planCmd.Path+" ")

	expected = "plan -no-color -input=false -lock-timeout=22s -out=whale -state=marvin -var-file=trillian -lock=false -parallelism=42 -refresh=false -destroy -target=zaphod -target=beeblebrox -var 'android=paranoid' -var 'brain_size=planet'"

	if actual != expected {
		t.Fatalf("expected arguments of PlanCmd:\n%s\n actual arguments:\n%s\n", expected, actual)
	}
}
