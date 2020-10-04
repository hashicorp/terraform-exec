// an example application using tfexec
package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

func main() {
	execPath, err := tfinstall.Find(context.Background(), tfinstall.LookPath(), tfinstall.ExactVersion("0.13.3", "/tmp"))
	if err != nil {
		panic(err)
	}

	workingDir := "/tmp"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		panic(err)
	}

	err = tf.Init(context.Background())
	if err != nil {
		panic(err)
	}

	tfVersion, _, err := tf.Version(context.Background(), true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("successfully initialised Terraform version %s", tfVersion)
}
