# terraform-exec

A Go module for constructing and running [Terraform](https://terraform.io) CLI commands. Structured return values use the data types defined in [terraform-json](https://github.com/hashicorp/terraform-json).

The [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) is the canonical Go interface for Terraform plugins using the gRPC protocol. This library is intended for use in Go programs that make use of Terraform's other interface, the CLI. Importing this library is preferable to importing `github.com/hashicorp/terraform/command`, because the latter is not intended for use outside Terraform Core.

This is not an official HashiCorp project.

## Usage

Top-level Terraform commands each have their own function, which will return either `error` or `(T, error)`, where `T` is a `terraform-json` type.

Convenience functions are provided for obtaining a raw `exec.Cmd` or command-line `string` for each supported Terraform command.


### Example


```go
package main

import (
    tfexec "github.com/kmoe/terraform-exec"
    
func main() {
    workingDir := "/path/to/working/dir"
    cfg := tfexec.Config{
        WorkingDir: workingDir,
    }

    state, err := cfg.Show()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(state.FormatVersion) // "0.1"
)
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
