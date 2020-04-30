package tfexec

import (
	"strings"
)

func buildTerraformArgs(args ...string) []string {
	tfPath := FindTerraform()
	allArgs := []string{tfPath}
	allArgs = append(allArgs, args...)
	allArgs = append(allArgs, "-no-color")
	return append(allArgs)
}

func initString(args ...string) []string {
	allArgs := append([]string{"init"}, args...)
	return buildTerraformArgs(allArgs...)
}

func InitString(args ...string) string {
	return strings.Join(initString(args...), " ")
}

func showString(args ...string) []string {
	allArgs := append([]string{"show", "-json"}, args...)
	return buildTerraformArgs(allArgs...)
}

func ShowString(args ...string) string {
	return strings.Join(showString(args...), " ")
}
