package tfexec

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestPlanCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		planCmd, err := tf.PlanCmd(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"plan",
			"-no-color",
			"-input=false",
			"-detailed-exitcode",
			"-lock-timeout=0s",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
		}, nil, planCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		planCmd, err := tf.PlanCmd(context.Background(),
			Destroy(true),
			Lock(false),
			LockTimeout("22s"),
			Out("whale"),
			Parallelism(42),
			Refresh(false),
			Replace("ford.prefect"),
			Replace("arthur.dent"),
			State("marvin"),
			Target("zaphod"),
			Target("beeblebrox"),
			Var("android=paranoid"),
			Var("brain_size=planet"),
			VarFile("trillian"),
			Dir("earth"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"plan",
			"-no-color",
			"-input=false",
			"-detailed-exitcode",
			"-lock-timeout=22s",
			"-out=whale",
			"-state=marvin",
			"-var-file=trillian",
			"-lock=false",
			"-parallelism=42",
			"-refresh=false",
			"-replace=ford.prefect",
			"-replace=arthur.dent",
			"-destroy",
			"-target=zaphod",
			"-target=beeblebrox",
			"-var", "android=paranoid",
			"-var", "brain_size=planet",
			"earth",
		}, nil, planCmd)
	})
}
