package tfexec

import (
	"context"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

type statePullConfig struct {
	reattachInfo ReattachInfo
}

var defaultStatePullConfig = statePullConfig{}

func (tf *Terraform) StatePull(ctx context.Context) (*tfjson.State, error) {
	c := defaultStatePullConfig

	mergeEnv := map[string]string{}
	if c.reattachInfo != nil {
		reattachStr, err := c.reattachInfo.marshalString()
		if err != nil {
			return nil, err
		}
		mergeEnv[reattachEnvVar] = reattachStr
	}

	cmd := tf.statePullCmd(ctx)

	var ret tfjson.State
	ret.UseJSONNumber(true)
	err := tf.runTerraformCmdJSON(ctx, cmd, &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *Terraform) statePullCmd(ctx context.Context) *exec.Cmd {
	args := []string{"state", "pull"}

	return tf.buildTerraformCmd(ctx, nil, args...)
}
