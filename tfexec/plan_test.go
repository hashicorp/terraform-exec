// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestPlanCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_5))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		planCmd, err := tf.planCmd(context.Background())
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
		planCmd, err := tf.planCmd(context.Background(),
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

	t.Run("run a refresh-only plan", func(t *testing.T) {
		planCmd, err := tf.planCmd(context.Background(), RefreshOnly(true))
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
			"-refresh-only",
		}, nil, planCmd)
	})

	t.Run("run a generate-config-out plan", func(t *testing.T) {
		planCmd, err := tf.planCmd(context.Background(), GenerateConfigOut("generated.tf"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"plan",
			"-no-color",
			"-input=false",
			"-detailed-exitcode",
			"-generate-config-out=generated.tf",
			"-lock-timeout=0s",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
		}, nil, planCmd)
	})
}

func TestPlanJSONCmd(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_5))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("defaults", func(t *testing.T) {
		planCmd, err := tf.planJSONCmd(context.Background())
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
			"-json",
		}, nil, planCmd)
	})

	t.Run("override all defaults", func(t *testing.T) {
		planCmd, err := tf.planJSONCmd(context.Background(),
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
			"-json",
			"earth",
		}, nil, planCmd)
	})

	t.Run("generate-config-out", func(t *testing.T) {
		planCmd, err := tf.planJSONCmd(context.Background(), GenerateConfigOut("generated.tf"))
		if err != nil {
			t.Fatal(err)
		}

		assertCmd(t, []string{
			"plan",
			"-no-color",
			"-input=false",
			"-detailed-exitcode",
			"-generate-config-out=generated.tf",
			"-lock-timeout=0s",
			"-lock=true",
			"-parallelism=10",
			"-refresh=true",
			"-json",
		}, nil, planCmd)
	})
}

func TestPlanCmd_AllowDeferral(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_Alpha_v1_9))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("allow deferrals during plan", func(t *testing.T) {
		planCmd, err := tf.planCmd(context.Background(), AllowDeferral(true))
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
			"-allow-deferral",
		}, nil, planCmd)
	})
}

func TestPlanCmd_MinimalRefresh(t *testing.T) {
	td := t.TempDir()

	tf, err := NewTerraform(td, tfVersion(t, testutil.Latest_v1_17))
	if err != nil {
		t.Fatal(err)
	}

	// empty env, to avoid environ mismatch in testing
	tf.SetEnv(map[string]string{})

	t.Run("minimal-refresh plan", func(t *testing.T) {
		planCmd, err := tf.planCmd(context.Background(), MinimalRefresh(true))
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
			"-minimal-refresh",
		}, nil, planCmd)
	})

	t.Run("minimal-refresh unsupported version", func(t *testing.T) {
		tf, err := NewTerraform(td, tfVersion(t, "1.16.0"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = tf.planCmd(context.Background(), MinimalRefresh(true))
		if err == nil {
			t.Fatal("expected an error but received none")
		}

		expectedErr := "minimal-refresh option was introduced in Terraform 1.17.0"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("expected error to contain %q, got: %q", expectedErr, err)
		}
	})

	t.Run("minimal-refresh and refresh=false returns an error", func(t *testing.T) {
		_, err := tf.planCmd(context.Background(), Refresh(false), MinimalRefresh(true))
		if err == nil {
			t.Fatal("expected an error but received none")
		}

		expectedErr := "cannot use refresh=false with mimimal-refresh"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("expected error to contain %q, got: %q", expectedErr, err)
		}
	})

	t.Run("minimal-refresh and refresh-only returns an error", func(t *testing.T) {
		_, err := tf.planCmd(context.Background(), RefreshOnly(true), MinimalRefresh(true))
		if err == nil {
			t.Fatal("expected an error but received none")
		}

		expectedErr := "cannot use refresh-only with mimimal-refresh"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("expected error to contain %q, got: %q", expectedErr, err)
		}
	})

	t.Run("minimal-refresh and destroy returns an error", func(t *testing.T) {
		_, err := tf.planCmd(context.Background(), Destroy(true), MinimalRefresh(true))
		if err == nil {
			t.Fatal("expected an error but received none")
		}

		expectedErr := "cannot use destroy with mimimal-refresh"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Fatalf("expected error to contain %q, got: %q", expectedErr, err)
		}
	})
}
