// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func (tf *TFCache) find(t *testing.T, key string, execPathFunc func(context.Context) (string, error)) string {
	t.Helper()

	if tf.dir == "" {
		// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
		panic("installDir not yet configured")
	}

	tf.Lock()
	defer tf.Unlock()

	if path, ok := tf.execs[key]; ok {
		return path
	}

	keyDir := key
	keyDir = strings.ReplaceAll(keyDir, ":", "-")
	keyDir = strings.ReplaceAll(keyDir, "/", "-")

	dir := filepath.Join(tf.dir, keyDir)

	t.Logf("caching exec %q in dir %q", key, dir)

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
		panic(fmt.Sprintf("unable to mkdir %q: %s", dir, err))
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	execPath, err := execPathFunc(ctx)
	if err != nil {
		// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
		panic(fmt.Sprintf("error installing terraform %q: %s", key, err))
	}

	tf.execs[key] = execPath

	return execPath
}
