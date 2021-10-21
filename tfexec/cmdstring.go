package tfexec

import (
	"os/exec"
)

func cmdString(c *exec.Cmd) string {
	return c.String()
}
