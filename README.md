# terraform-exec
A Go module for constructing and running [Terraform](https://terraform.io) CLI commands. Structured return values use the data types defined in [terraform-json](https://github.com/hashicorp/terraform-json).

The [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) is the canonical Go interface for Terraform plugins using the gRPC protocol. This library is intended for use in Go programs that make use of Terraform's other interface, the CLI. Importing this library is preferable to importing `github.com/hashicorp/terraform/command`, because the latter is not intended for use outside Terraform Core.

This is not an official HashiCorp project.

## Usage

Top-level Terraform commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.

Convenience functions are provided for obtaining a raw `exec.Cmd` or command-line `string` for each supported Terraform command.

Currently supported Terraform commands:
* init
* show

### Example


```go
package main

import (
        "fmt"
        "os"

        tfexec "github.com/kmoe/terraform-exec"
)

func main() {
        workingDir := os.Getenv("GOPATH") + "/src/github.com/kmoe/terraform-exec/testdata"

        // Run `terraform init` so that the working directories state can be initialized.
        err := tfexec.Init(workingDir)
        if err != nil {
                panic(err)
        }

        // Run `terraform show` against the state defined in the working directory.
        state, err := tfexec.Show(workingDir)
        if err != nil {
                panic(err)
        }

        // Print all returned values from the `terraform show` command (of type *tfjson.State)
        fmt.Println(state.FormatVersion) // "0.1"
        fmt.Println(state.TerraformVersion)
        fmt.Println(state.Values)
}
```

### `Init(workingDir string, args ...string) error`

Runs `terraform init` in the given directory.

### `Show(workingDir string, args ...string) (*tfjson.State, error)`

Returns the output of `terraform show -json`, represented as `tfjson.State`.


### `exec.Cmd` functions 

#### `InitCmd(workingDir string, args ...string) exec.Cmd`

#### `ShowCmd(workingDir string, args ...string) exec.Cmd`

### `string` functions

#### `InitString(args ...string) string`

#### `ShowString(args ...string) string`
