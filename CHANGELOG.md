# 0.13.3 (Unreleased)

SECURITY:
 - `tfinstall`: The HashiCorp PGP signing key has been rotated (HCSEC-2021-12). This key is used to verify downloaded versions of Terraform. We recommend all users of terraform-exec upgrade to v0.13.3 for this security fix. [GH-166]

# 0.13.2 (April 06, 2021)

BUG FIXES:
 - Update `terraform-json` to support 0.15 changes in plan & config JSON output ([#153](https://github.com/hashicorp/terraform-exec/issues/153))
 - Update `go-getter` to prevent race conditions where consumers would require `go-cleanhttp` `>=0.5.2` (which tfexec itself didn't depend on until now) ([#154](https://github.com/hashicorp/terraform-exec/issues/154))

# 0.13.1 (March 29, 2021)

BUG FIXES:
 - Bump version of terraform-json library to handle latest Terraform 0.15 output format ([#143](https://github.com/hashicorp/terraform-exec/issues/143))

NOTES:
 - This release no longer supports Go 1.12 (1.13+ is required)

# 0.13.0 (February 05, 2021)

Please note that this is the first release of terraform-exec compatible with Terraform 0.15. Running Terraform 0.15 commands with previous versions of terraform-exec may produce unexpected results.

FEATURES:
 - Compatibility checks for CLI flags removed in Terraform 0.15 ([#120](https://github.com/hashicorp/terraform-exec/issues/120))
 - Introduce `StateRm` method ([#122](https://github.com/hashicorp/terraform-exec/issues/122))

# 0.12.0 (December 18, 2020)

BREAKING CHANGES:
 - Move Git ref installation to subpackage so that consumers can limit dependencies ([#98](https://github.com/hashicorp/terraform-exec/issues/98))

FEATURES:
 - Improve error handling for formatting command on unsupported version (`<0.7.7`) ([#88](https://github.com/hashicorp/terraform-exec/issues/88))
 - Introduce `Format` method with `io.Reader`/`io.Writer` interfaces ([#96](https://github.com/hashicorp/terraform-exec/issues/96))
 - Introduce `Validate` method with `tfjson` defined diagnostic types. Those types reflect exactly the types used in `terraform validate -json` output ([#68](https://github.com/hashicorp/terraform-exec/issues/68))
 - Introduce `StateMv` method ([#112](https://github.com/hashicorp/terraform-exec/issues/112))
 - Introduce `Upgrade012` method ([#105](https://github.com/hashicorp/terraform-exec/issues/105))

BUG FIXES:
 - Fix issue in tfinstall.GitRef where it assumed a `vendor` directory was present ([#89](https://github.com/hashicorp/terraform-exec/issues/89))
 - Use `json.Number` instead of `float64` when parsing state ([#113](https://github.com/hashicorp/terraform-exec/issues/113))
 - Support long variable names in `ErrMissingVar` ([#110](https://github.com/hashicorp/terraform-exec/issues/110))

# 0.11.0 (September 23, 2020)

FEATURES:
 - Added Terraform fmt support with the ability to format and write files/folders, check if files/folders need formatting, and format strings directly ([#82](https://github.com/hashicorp/terraform-exec/issues/82))
 - Added support for refs in the tfinstall CLI ([#80](https://github.com/hashicorp/terraform-exec/issues/80))

N.B. tfinstall binaries for all supported platforms are now available via GitHub Releases.

# 0.10.0 (September 15, 2020)

FEATURES:
 - Added the ability to customize the `User-Agent` header for some `tfinstall` finders ([#76](https://github.com/hashicorp/terraform-exec/issues/76))
 - Added well known error for a mismatch for `required_version` ([#66](https://github.com/hashicorp/terraform-exec/issues/66))
 - Added new `ShowPlanFileRaw` function to obtain the human-friendly output of a plan ([#83](https://github.com/hashicorp/terraform-exec/issues/83))

# 0.9.0 (September 09, 2020)

BREAKING CHANGES:
 - `context.Context` added to `tfinstall.Find` to allow for cancellation, timeouts, etc ([#51](https://github.com/hashicorp/terraform-exec/issues/51))
 - You can no longer use `TF_WORKSPACE` for workspace management, you must use `Terraform.WorkspaceSelect` ([#75](https://github.com/hashicorp/terraform-exec/issues/75))

FEATURES:
 - Add `ErrWorkspaceExists` for when workspaces with the same name already exist when calling `Terraform.WorkspaceNew` ([#67](https://github.com/hashicorp/terraform-exec/issues/67))
 - Added `tfinstall.GitRef` to support installation of Terraform from a git ref instead of by released version ([#51](https://github.com/hashicorp/terraform-exec/issues/51))
 - Created the **tfinstall** CLI utility (this is mostly for use in things like CI automation) ([#29](https://github.com/hashicorp/terraform-exec/issues/29))
 - Added `ReattachOption` for plugin reattach functionality ([#78](https://github.com/hashicorp/terraform-exec/issues/78))

# 0.8.0 (August 29, 2020)

BREAKING CHANGES:
 - Add `-detailed-exit-code` to `Terraform.Plan` calls, `Terraform.Plan` now also returns a bool indicating if any diff is present ([#55](https://github.com/hashicorp/terraform-exec/issues/55))
 
FEATURES:
 - Added `Terraform.SetAppendUserAgent` for User-Agent management from consuming applications ([#46](https://github.com/hashicorp/terraform-exec/issues/46))
 - Added `Terraform.WorkspaceList`, `Terraform.WorkspaceNew`, and `Terraform.WorkspaceSelect` along with the `ErrNoWorkspace` error to indicate a workspace does not exist ([#56](https://github.com/hashicorp/terraform-exec/issues/56))
 - Added support for using multiple `VarFile` options ([#61](https://github.com/hashicorp/terraform-exec/issues/61))
 
BUG FIXES:
 - Fix bug with checking for empty path before executing version command ([#62](https://github.com/hashicorp/terraform-exec/issues/62))

# 0.7.0 (August 20, 2020)

FEATURES:
 - Added `Terraform.Refresh` method ([#53](https://github.com/hashicorp/terraform-exec/issues/53))
 - Added `Terraform.ShowStateFile` and `Terraform.ShowPlanFile` ([#54](https://github.com/hashicorp/terraform-exec/issues/54))
 - Added support for `DIR` positional arg in init, destroy, and plan ([#52](https://github.com/hashicorp/terraform-exec/issues/52))
 - Relaxed logger interface ([#57](https://github.com/hashicorp/terraform-exec/issues/57))
 - Added error for missing required variable ([#57](https://github.com/hashicorp/terraform-exec/issues/57))

BUG FIXES:
 - Fixed logging issue for error cmd ([#57](https://github.com/hashicorp/terraform-exec/issues/57))

# 0.6.0 (August 14, 2020)

FEATURES:
 - Added `Terraform.SetStdout` and `Terraform.SetStderr` to let consumers log CLI output ([#49](https://github.com/hashicorp/terraform-exec/issues/49))

BUG FIXES:
 - Fixed miscategorization of `ErrNoInit` on Terraform 0.13 ([#48](https://github.com/hashicorp/terraform-exec/issues/48))

# 0.5.0 (August 14, 2020)

FEATURES:
 - Version compatibility testing for `terraform show` ([#41](https://github.com/hashicorp/terraform-exec/issues/41))

BUG FIXES:
 - Tolerate reversed `terraform version` output order ([#47](https://github.com/hashicorp/terraform-exec/issues/47))

# 0.4.0 (July 30, 2020)

FEATURES:
  - Added `Terraform.SetLogPath` method to set `TF_LOG_PATH` environment variable, and prevented manual setting of programmatically supported environment variables ([#32](https://github.com/hashicorp/terraform-exec/issues/32))
  - Added `Terraform.Version` method to get executable version information ([#7](https://github.com/hashicorp/terraform-exec/issues/7))

BUG FIXES:
  - Fixed `-var` handling issue ([#34](https://github.com/hashicorp/terraform-exec/issues/34))

# 0.3.0 (July 17, 2020)

BREAKING CHANGES:
  - Stop exporting `exec.Cmd` versions of methods ([#25](https://github.com/hashicorp/terraform-exec/issues/25))
  - Require `address` and `id` arguments in `Import()` ([#24](https://github.com/hashicorp/terraform-exec/issues/24))
  - Rename `StateShow()` to `Show()` ([#30](https://github.com/hashicorp/terraform-exec/issues/30))

BUG FIXES:
  - Fix bug in `Import()` config argument ([#28](https://github.com/hashicorp/terraform-exec/issues/28))

# 0.2.2 (July 13, 2020)

BUG FIXES:
  - Version number is now correctly reported by the tfinstall package. Please note that `tfinstall.Version` was incorrect between versions 0.1.1 and 0.2.1 inclusive.

# 0.2.1 (July 10, 2020)

BUG FIXES:
  - Minor code changes to allow for compilation in Go 1.12 ([#21](https://github.com/hashicorp/terraform-exec/pull/21))

# 0.2.0 (July 8, 2020)

FEATURES:
  - add `Import()` function ([#20](https://github.com/hashicorp/terraform-exec/pull/20))

# 0.1.1 (July 7, 2020)

BUG FIXES:
 - Downgrade `github.com/hashicorp/go-getter` dependency, which added a requirement for Go 1.13.

# 0.1.0 (July 3, 2020)

Initial release.

This Go module contains two packages, `github.com/hashicorp/terraform-exec/tfexec`, and `github.com/hashicorp/terraform-exec/tfinstall`, which share the same version.
