package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

const (
	Latest011 = "0.11.14"
	Latest012 = "0.12.29"
	Latest013 = "0.13.0"
)

type TFCache struct {
	sync.Mutex

	dir   string
	execs map[string]string
}

func NewTFCache(dir string) *TFCache {
	return &TFCache{
		dir:   dir,
		execs: map[string]string{},
	}
}

func (tf *TFCache) GitRef(t *testing.T, ref string) string {
	t.Helper()
	return tf.find(t, "gitref:"+ref, func(dir string) tfinstall.ExecPathFinder {
		return tfinstall.GitRef(ref, "", dir)
	})
}

func (tf *TFCache) Version(t *testing.T, v string) string {
	t.Helper()
	return tf.find(t, "v:"+v, func(dir string) tfinstall.ExecPathFinder {
		return tfinstall.ExactVersion(v, dir)
	})
}

func (tf *TFCache) find(t *testing.T, key string, finder func(dir string) tfinstall.ExecPathFinder) string {

	t.Helper()

	if tf.dir == "" {
		// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
		panic("installDir not yet configured")
	}

	tf.Lock()
	defer tf.Unlock()

	path, ok := tf.execs[key]
	if !ok {
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

		path, err = tfinstall.Find(context.Background(), finder(dir))
		if err != nil {
			// panic instead of t.fatal as this is going to affect all downstream tests reusing the cache entry
			panic(fmt.Sprintf("error installing terraform %q: %s", key, err))
		}
		tf.execs[key] = path
	}

	return path
}
