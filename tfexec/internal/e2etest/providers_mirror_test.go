package e2etest

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfexec/internal/testutil"
)

func TestProvidersMirror(t *testing.T) {
	runTestVersions(t, []string{testutil.Latest013, testutil.Latest014, testutil.Latest015, testutil.Latest_v1}, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		err := tf.Init(context.Background())
		if err != nil {
			t.Fatalf("error running Init in test directory: %s", err)
		}

		err = tf.ProvidersMirror(context.Background(), "targetDir")
		if err != nil {
			t.Fatalf("error running provider mirror: %s", err)
		}
	})

}
