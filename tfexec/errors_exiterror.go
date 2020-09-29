// +build !go1.12,!go1.13

package tfexec

import (
	"errors"
	"fmt"
	"os/exec"
)

func (e *ExitError) Error() string {
	var out string
	ee, ok := e.err.(*exec.ExitError)
	if ok {
		out = fmt.Sprintf("%q (pid %d) exited (code %d): %s",
			e.args,
			ee.Pid(),
			ee.ExitCode(),
			ee.ProcessState.String())
		if e.ctxErr != nil {
			out += fmt.Sprintf("\n%s", e.ctxErr)
		}
	} else {
		out = fmt.Sprintf("%q exited: %s", e.args, e.err.Error())
		if e.ctxErr != nil && !errors.Is(e.err, e.ctxErr) {
			out += fmt.Sprintf("\n%s", e.ctxErr)
		}
	}

	if e.stderr != "" {
		out += fmt.Sprintf("\nstderr: %q", e.stderr)
	}

	return out
}
