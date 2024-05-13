// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-exec/internal/version"
)

const (
	checkpointDisableEnvVar  = "CHECKPOINT_DISABLE"
	cliArgsEnvVar            = "TF_CLI_ARGS"
	inputEnvVar              = "TF_INPUT"
	automationEnvVar         = "TF_IN_AUTOMATION"
	logEnvVar                = "TF_LOG"
	logCoreEnvVar            = "TF_LOG_CORE"
	logPathEnvVar            = "TF_LOG_PATH"
	logProviderEnvVar        = "TF_LOG_PROVIDER"
	reattachEnvVar           = "TF_REATTACH_PROVIDERS"
	appendUserAgentEnvVar    = "TF_APPEND_USER_AGENT"
	workspaceEnvVar          = "TF_WORKSPACE"
	disablePluginTLSEnvVar   = "TF_DISABLE_PLUGIN_TLS"
	skipProviderVerifyEnvVar = "TF_SKIP_PROVIDER_VERIFY"

	varEnvVarPrefix    = "TF_VAR_"
	cliArgEnvVarPrefix = "TF_CLI_ARGS_"
)

var prohibitedEnvVars = []string{
	cliArgsEnvVar,
	inputEnvVar,
	automationEnvVar,
	logEnvVar,
	logCoreEnvVar,
	logPathEnvVar,
	logProviderEnvVar,
	reattachEnvVar,
	appendUserAgentEnvVar,
	workspaceEnvVar,
	disablePluginTLSEnvVar,
	skipProviderVerifyEnvVar,
}

var prohibitedEnvVarPrefixes = []string{
	varEnvVarPrefix,
	cliArgEnvVarPrefix,
}

func manualEnvVars(env map[string]string, cb func(k string)) {
	for k := range env {
		for _, p := range prohibitedEnvVars {
			if p == k {
				cb(k)
				goto NextEnvVar
			}
		}
		for _, prefix := range prohibitedEnvVarPrefixes {
			if strings.HasPrefix(k, prefix) {
				cb(k)
				goto NextEnvVar
			}
		}
	NextEnvVar:
	}
}

// ProhibitedEnv returns a slice of environment variable keys that are not allowed
// to be set manually from the passed environment.
func ProhibitedEnv(env map[string]string) []string {
	var p []string
	manualEnvVars(env, func(k string) {
		p = append(p, k)
	})
	return p
}

// CleanEnv removes any prohibited environment variables from an environment map.
func CleanEnv(dirty map[string]string) map[string]string {
	clean := dirty
	manualEnvVars(clean, func(k string) {
		delete(clean, k)
	})
	return clean
}

func envMap(environ []string) map[string]string {
	env := map[string]string{}
	for _, ev := range environ {
		parts := strings.SplitN(ev, "=", 2)
		if len(parts) == 0 {
			continue
		}
		k := parts[0]
		v := ""
		if len(parts) == 2 {
			v = parts[1]
		}
		env[k] = v
	}
	return env
}

func envSlice(environ map[string]string) []string {
	env := []string{}
	for k, v := range environ {
		env = append(env, k+"="+v)
	}
	return env
}

func (tf *Terraform) buildEnv(mergeEnv map[string]string) []string {
	// set Terraform level env, if env is nil, fall back to os.Environ
	var env map[string]string
	if tf.env == nil {
		env = envMap(os.Environ())
	} else {
		env = make(map[string]string, len(tf.env))
		for k, v := range tf.env {
			env[k] = v
		}
	}

	// override env with any command specific environment
	for k, v := range mergeEnv {
		env[k] = v
	}

	// always propagate CHECKPOINT_DISABLE env var unless it is
	// explicitly overridden with tf.SetEnv or command env
	if _, ok := env[checkpointDisableEnvVar]; !ok {
		env[checkpointDisableEnvVar] = os.Getenv(checkpointDisableEnvVar)
	}

	// always override user agent
	ua := mergeUserAgent(
		os.Getenv(appendUserAgentEnvVar),
		tf.appendUserAgent,
		fmt.Sprintf("HashiCorp-terraform-exec/%s", version.ModuleVersion()),
	)
	env[appendUserAgentEnvVar] = ua

	// always override logging
	if tf.logPath == "" {
		// so logging can't pollute our stderr output
		env[logEnvVar] = ""
		env[logCoreEnvVar] = ""
		env[logPathEnvVar] = ""
		env[logProviderEnvVar] = ""
	} else {
		env[logEnvVar] = tf.log
		env[logCoreEnvVar] = tf.logCore
		env[logPathEnvVar] = tf.logPath
		env[logProviderEnvVar] = tf.logProvider
	}

	// constant automation override env vars
	env[automationEnvVar] = "1"

	// force usage of workspace methods for switching
	delete(env, workspaceEnvVar)

	if tf.disablePluginTLS {
		env[disablePluginTLSEnvVar] = "1"
	}

	if tf.skipProviderVerify {
		env[skipProviderVerifyEnvVar] = "1"
	}

	return envSlice(env)
}

func (tf *Terraform) buildTerraformCmd(ctx context.Context, mergeEnv map[string]string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, tf.execPath, args...)

	cmd.Env = tf.buildEnv(mergeEnv)
	cmd.Dir = tf.workingDir

	tf.logger.Printf("[INFO] running Terraform command: %s", cmd.String())

	return cmd
}

func (tf *Terraform) buildTerraformCmdWithGracefulShutdown(ctx context.Context, gracefulShutdownPeriod time.Duration, mergeEnv map[string]string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, tf.execPath, args...)

	cmd.Env = tf.buildEnv(mergeEnv)
	cmd.Dir = tf.workingDir

	tf.logger.Printf("[INFO] configuring gracefully cancellation for command: %s", cmd.String())
	cmd.WaitDelay = gracefulShutdownPeriod
	cmd.Cancel = func() error {
		tf.logger.Printf("[INFO] gracefully cancelling Terraform command: %s", cmd.String())
		// os.Interrupt doesn't work on Windows so it will result into an immediate kill.
		err := cmd.Process.Signal(os.Interrupt)
		if err != nil && !errors.Is(err, os.ErrProcessDone) {
			tf.logger.Printf("[ERROR] gracefully cancelling failed due to %s, forcefully cancelling Terraform command: %s", err.Error(), cmd.String())
			err = cmd.Process.Kill()
		}
		return err
	}

	tf.logger.Printf("[INFO] running Terraform command: %s", cmd.String())

	return cmd
}

func (tf *Terraform) runTerraformCmdJSON(ctx context.Context, cmd *exec.Cmd, v interface{}) error {
	var outbuf = bytes.Buffer{}
	cmd.Stdout = mergeWriters(cmd.Stdout, &outbuf)

	err := tf.runTerraformCmd(ctx, cmd)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(&outbuf)
	dec.UseNumber()
	return dec.Decode(v)
}

// mergeUserAgent does some minor deduplication to ensure we aren't
// just using the same append string over and over.
func mergeUserAgent(uas ...string) string {
	included := map[string]bool{}
	merged := []string{}
	for _, ua := range uas {
		ua = strings.TrimSpace(ua)

		if ua == "" {
			continue
		}
		if included[ua] {
			continue
		}
		included[ua] = true
		merged = append(merged, ua)
	}
	return strings.Join(merged, " ")
}

func mergeWriters(writers ...io.Writer) io.Writer {
	compact := []io.Writer{}
	for _, w := range writers {
		if w != nil {
			compact = append(compact, w)
		}
	}
	if len(compact) == 0 {
		return ioutil.Discard
	}
	if len(compact) == 1 {
		return compact[0]
	}
	return io.MultiWriter(compact...)
}

func writeOutput(ctx context.Context, r io.ReadCloser, w io.Writer) error {
	// ReadBytes will block until bytes are read, which can cause a delay in
	// returning even if the command's context has been canceled. Use a separate
	// goroutine to prompt ReadBytes to return on cancel
	closeCtx, closeCancel := context.WithCancel(ctx)
	defer closeCancel()
	go func() {
		select {
		case <-ctx.Done():
			r.Close()
		case <-closeCtx.Done():
			return
		}
	}()

	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadBytes('\n')
		if len(line) > 0 {
			if _, err := w.Write(line); err != nil {
				return err
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}
	}
}

func (tf *Terraform) runTerraformCmd(ctx context.Context, cmd *exec.Cmd) error {
	var errBuf strings.Builder

	cmd.SysProcAttr = defaultSysProcAttr

	// check for early cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Read stdout / stderr logs from pipe instead of setting cmd.Stdout and
	// cmd.Stderr because it can cause hanging when killing the command
	// https://github.com/golang/go/issues/23019
	stdoutWriter := mergeWriters(cmd.Stdout, tf.stdout)
	stderrWriter := mergeWriters(tf.stderr, &errBuf)

	cmd.Stderr = nil
	cmd.Stdout = nil

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if ctx.Err() != nil {
		return cmdErr{
			err:    err,
			ctxErr: ctx.Err(),
		}
	}
	if err != nil {
		return err
	}

	var errStdout, errStderr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		errStdoutCtx := ctx
		// When custom cancellation is enabled, avoid piping cancellation to prevent
		// cancellation from being interrupted due to broken pipe.
		if cmd.Cancel != nil {
			errStdoutCtx = context.WithoutCancel(ctx)
		}
		errStdout = writeOutput(errStdoutCtx, stdoutPipe, stdoutWriter)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errStderrCtx := ctx
		// When custome cancellation is enabled, avoid piping cancellation to prevent
		// cancellation from being interrupted due to broken pipe.
		if cmd.Cancel != nil {
			errStderrCtx = context.WithoutCancel(ctx)
		}
		errStderr = writeOutput(errStderrCtx, stderrPipe, stderrWriter)
	}()

	// Reads from pipes must be completed before calling cmd.Wait(). Otherwise
	// can cause a race condition
	wg.Wait()

	err = cmd.Wait()
	if ctx.Err() != nil {
		return cmdErr{
			err:    err,
			ctxErr: ctx.Err(),
		}
	}
	if err != nil {
		return fmt.Errorf("%w\n%s", err, errBuf.String())
	}

	// Return error if there was an issue reading the std out/err
	if errStdout != nil && ctx.Err() != nil {
		return fmt.Errorf("%w\n%s", errStdout, errBuf.String())
	}
	if errStderr != nil && ctx.Err() != nil {
		return fmt.Errorf("%w\n%s", errStderr, errBuf.String())
	}

	return nil
}
