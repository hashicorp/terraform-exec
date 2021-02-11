# Contributing to terraform-exec

While terraform-exec is already widely used, please note that this module is **not yet at v1.0.0**, and that therefore breaking changes may occur in minor releases.

We strictly follow [semantic versioning](https://semver.org).

## Repository structure

Three packages comprise the public API of the terraform-exec Go module:

### `tfexec`

Package `github.com/hashicorp/terraform-exec/tfexec` exposes functionality for constructing and running Terraform CLI commands. Structured return values use the data types defined in the [hashicorp/terraform-json](https://github.com/hashicorp/terraform-json) package.

#### Adding a new Terraform CLI command to `tfexec`

Each Terraform CLI first- or second-level subcommand (e.g. `terraform refresh`, or `terraform workspace new`) is implemented in a separate Go file. This file defines a public function on the `Terraform` struct, which consumers use to call the CLI command, and a private `*Cmd` version of this function which returns an `*exec.Cmd`, in order to facilitate unit testing (see Testing below).

For example:
```go
func (tf *Terraform) Refresh(ctx context.Context, opts ...RefreshCmdOption) error {
	cmd, err := tf.refreshCmd(ctx, opts...)
	if err != nil {
		return err
	}
	return tf.runTerraformCmd(cmd)
}

func (tf *Terraform) refreshCmd(ctx context.Context, opts ...RefreshCmdOption) (*exec.Cmd, error) {
	...
  	return tf.buildTerraformCmd(ctx, mergeEnv, args...), nil
}
```

Command options are implemented using the functional variadic options pattern. For further reading on this pattern, please see [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) by Dave Cheney.

### `tfinstall`

Package `github.com/hashicorp/terraform-exec/tfinstall` offers multiple strategies for finding and/or installing a binary version of Terraform. Some of the strategies can also authenticate the source of the binary as an official HashiCorp release.

**The public API of `tfinstall` is experimental and we welcome suggestions and ideas for its improvement, and accounts of current or proposed use cases.**

As in `tfexec`, a variadic options pattern is used to supply a list of strategies for locating the Terraform executable. Each is implemented in a separate Go file. Our aim is to present a flexible and reasonably intuitive interface for ensuring a Terraform executable is available, with or without version constraints.

### `cmd/tfinstall`

Package `github.com/hashicorp/terraform-exec/cmd/tfinstall` has no public Go API and is distributed as a compiled binary which can be run to download and install a version of Terraform. 

## Testing

We aim for full test coverage of all Terraform CLI commands implemented in `tfexec`, with as many combinations of command-line options as possible. New command implementations will not be merged without both unit and end-to-end tests.

### Unit tests

Unit tests live alongside command implementations in `tfexec/`. A unit test asserts that the *string* version of the `exec.Cmd` returned by the `*Cmd` function (e.g. `refreshCmd`) is as expected. Minimally, commands must be tested with no options passed ("defaults"), and with all options set to non-default values. The `assertCmd()` helper can be used for this purpose. Please see `tfexec/init_test.go` for a reasonable starting point.

### End-to-end tests

End-to-end tests test both `tfinstall` and `tfexec`, using the former to install Terraform binaries according to various version constraints, and exercising the latter in as many combinations as possible, after real-world use cases.

By default, each test is run against the latest versions of Terraform 0.11, Terraform 0.12, and Terraform 0.13. Copy an existing test and use the `runTest()` helper for this purpose.

If the command implemented differs in any way between Terraform versions (e.g. a flag is added or removed, or the subcommand does not exist in earlier versions), the `runTestVersions()` helper should be used to make both positive and negative assertions against the Terraform version being run, or skip the test as appropriate.

## Versioning

The `github.com/hashicorp/terraform-exec` Go module in its entirety is versioned according to [Go module versioning](https://golang.org/ref/mod#versions) with Git tags. The latest version is automatically written to `internal/version/version.go` during the release process.

## Releases

Releases are made on a reasonably regular basis by the Terraform Plugin SDK team, using our custom CI workflows. There is currently no set release schedule and no requirement for contributors to write CHANGELOG entries.

## Security vulnerabilities

Please disclose security vulnerabilities by following the procedure
described at https://www.hashicorp.com/security#vulnerability-reporting.


## Cosmetic changes, code formatting, and typos

In general we do not accept PRs containing only the following changes:

 - Correcting spelling or typos
 - Code formatting, including whitespace
 - Other cosmetic changes that do not affect functionality
 
While we appreciate the effort that goes into preparing PRs, there is always a tradeoff between benefit and cost. The costs involved in accepting such contributions include the time taken for thorough review, the noise created in the git history, and the increased number of GitHub notifications that maintainers must attend to.

In the case of `terraform-plugin-sdk`, the repo's close relationship to the `terraform` repo means that maintainers will sometimes port changes from `terraform` to `terraform-plugin-sdk`. Cosmetic changes to the SDK repo make this much more time-consuming as they will cause merge conflicts. This is the major hidden cost of cosmetic PRs, and the main reason we do not accept them at this time.

#### Exceptions

We belive that one should "leave the campsite cleaner than you found it", so you are welcome to clean up cosmetic issues in the neighbourhood when submitting a patch that makes functional changes or fixes.
