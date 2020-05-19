# terraform-exec

A Go module for constructing and running [Terraform](https://terraform.io) CLI commands. Structured return values use the data types defined in [terraform-json](https://github.com/hashicorp/terraform-json).

The [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) is the canonical Go interface for Terraform plugins using the gRPC protocol. This library is intended for use in Go programs that make use of Terraform's other interface, the CLI. Importing this library is preferable to importing `github.com/hashicorp/terraform/command`, because the latter is not intended for use outside Terraform Core.

This is not an official HashiCorp project.

## Usage

The `Terraform` struct must be initialised with `NewTerraform(workingDir, execPath)`. 

Top-level Terraform commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.


### Example


```go
package main

import (
    "github.com/kmoe/terraform-exec/tfexec"
    
func main() {
    workingDir := "/path/to/working/dir"
    tf, err := tfexec.NewTerraform(workingDir)
    if err != nil {
        panic(err)
    }

    state, err := tf.Init()
    if err != nil {
        panic(err)
    }

    fmt.Println(state.FormatVersion) // "0.1"
)
```


### `(tf *Terraform) Init(opts ...InitOpts) (*tfjson.State, error)`

Runs `terraform init`.

### `(tf *Terraform) Show(opts ...ShowOpts) (*tfjson.State, error)`

Returns the output of `terraform show -json`, represented as `tfjson.State`.

### `exec.Cmd` functions 

#### `InitCmd(workingDir string, args ...string) exec.Cmd`

#### `ShowCmd(workingDir string, args ...string) exec.Cmd`
