// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package e2etest

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

var (
	validateMinVersion = version.Must(version.NewVersion("0.12.0"))
)

func TestValidate(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(validateMinVersion) {
			t.Skip("terraform validate -json was added in Terraform 0.12, so test is not valid")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		validation, err := tf.Validate(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if !validation.Valid {
			t.Fatalf("expected valid, got %#v", validation)
		}
	})

	runTest(t, "invalid", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		if tfv.LessThan(validateMinVersion) {
			t.Skip("terraform validate -json was added in Terraform 0.12, so test is not valid")
		}

		err := tf.Init(context.Background())
		if err != nil {
			t.Logf("error initializing: %s", err)

			// 0.13 will error, 0.12 will not
			// unsure why 0.12 terraform init does not have a non-zero exit code for syntax problems
			if err == nil {
				t.Fatalf("expected error, but did not get one")
			}
		}

		var expectedDiags []tfjson.Diagnostic

		if tfv.Core().GreaterThanOrEqual(v0_15_0) {
			expectedDiags = []tfjson.Diagnostic{
				{
					Severity: "error",
					Summary:  "Unsupported block type",
					Detail:   "Blocks of type \"bad_block\" are not expected here.",
					Range: &tfjson.Range{
						Filename: "main.tf",
						Start: tfjson.Pos{
							Line:   1,
							Column: 1,
						},
						End: tfjson.Pos{
							Line:   1,
							Column: 10,
							Byte:   9,
						},
					},
					Snippet: &tfjson.DiagnosticSnippet{
						Code:                 "bad_block {",
						StartLine:            1,
						HighlightStartOffset: 0,
						HighlightEndOffset:   9,
						Values:               []tfjson.DiagnosticExpressionValue{},
					},
				},
				{
					Severity: "error",
					Summary:  "Unsupported argument",
					Detail:   "An argument named \"bad_attribute\" is not expected here.",
					Range: &tfjson.Range{
						Filename: "main.tf",
						Start: tfjson.Pos{
							Line:   5,
							Column: 5,
							Byte:   31,
						},
						End: tfjson.Pos{
							Line:   5,
							Column: 18,
							Byte:   44,
						},
					},
					Snippet: &tfjson.DiagnosticSnippet{
						Context:              ptrToString("terraform"),
						Code:                 "    bad_attribute = \"string\"",
						StartLine:            5,
						HighlightStartOffset: 4,
						HighlightEndOffset:   17,
						Values:               []tfjson.DiagnosticExpressionValue{},
					},
				},
			}
		} else {
			expectedDiags = []tfjson.Diagnostic{
				{
					Severity: "error",
					Summary:  "Unsupported block type",
					Detail:   "Blocks of type \"bad_block\" are not expected here.",
					Range: &tfjson.Range{
						Filename: "main.tf",
						Start: tfjson.Pos{
							Line:   1,
							Column: 1,
						},
						End: tfjson.Pos{
							Line:   1,
							Column: 10,
							Byte:   9,
						},
					},
				},
				{
					Severity: "error",
					Summary:  "Unsupported argument",
					Detail:   "An argument named \"bad_attribute\" is not expected here.",
					Range: &tfjson.Range{
						Filename: "main.tf",
						Start: tfjson.Pos{
							Line:   5,
							Column: 5,
							Byte:   31,
						},
						End: tfjson.Pos{
							Line:   5,
							Column: 18,
							Byte:   44,
						},
					},
				},
			}
		}

		actual, err := tf.Validate(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(expectedDiags, actual.Diagnostics); diff != "" {
			t.Fatalf("diags do not match: %s", diff)
		}
	})
}

func ptrToString(value string) *string {
	return &value
}
