# 0.2.1 (Unreleased)

BUG FIXES:
  - Minor code changes to allow for compilation in Go 1.12 ([#21](https://github.com/hashicorp/terraform-exec/pull/21))

# 0.2.0 (July 8, 2020)

NEW FEATURES:
  - add `Import()` function ([#20](https://github.com/hashicorp/terraform-exec/pull/20))

# 0.1.1 (July 7, 2020)

BUG FIXES:
 - Downgrade `github.com/hashicorp/go-getter` dependency, which added a requirement for Go 1.13.

# 0.1.0 (July 3, 2020)

Initial release. 

This Go module contains two packages, `github.com/hashicorp/terraform-exec/tfexec`, and `github.com/hashicorp/terraform-exec/tfinstall`, which share the same version.
