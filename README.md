[![PkgGoDev](https://pkg.go.dev/badge/github.com/hashicorp/terraform-exec)](https://pkg.go.dev/github.com/hashicorp/terraform-exec)

# terraform-exec

A Go module for constructing and running [Terraform](https://terraform.io) CLI commands. Structured return values use the data types defined in [terraform-json](https://github.com/hashicorp/terraform-json).

The [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) is the canonical Go interface for Terraform plugins using the gRPC protocol. This library is intended for use in Go programs that make use of Terraform's other interface, the CLI. Importing this library is preferable to importing `github.com/hashicorp/terraform/command`, because the latter is not intended for use outside Terraform Core.

While terraform-exec is already widely used, please note that this module is **not yet at v1.0.0**, and that therefore breaking changes may occur in minor releases.

We strictly follow [semantic versioning](https://semver.org).

## Usage

The `Terraform` struct must be initialised with `NewTerraform(workingDir, execPath)`. 

Top-level Terraform commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.


### Example


```go
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"log"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

func main() {
	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		log.Fatalf("error creating temp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	execPath, err := tfinstall.Find(context.Background(), tfinstall.LatestVersion(tmpDir, false))
	if err != nil {
		log.Fatalf("error locating Terraform binary: %s", err)
	}

	workingDir := "/path/to/working/dir"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	fmt.Println(state.FormatVersion) // "0.1"
}
```

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md).
