package testutil

import (
	"context"
	"os"
	"path/filepath"
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

	dir      string
	versions map[string]string
}

func NewTFCache(dir string) *TFCache {
	return &TFCache{
		dir:      dir,
		versions: map[string]string{},
	}
}

func (tf *TFCache) Version(t *testing.T, v string) string {
	t.Helper()

	if tf.dir == "" {
		t.Fatalf("installDir not yet configured")
	}

	tf.Lock()
	defer tf.Unlock()

	path, ok := tf.versions[v]
	if !ok {
		t.Logf("caching version %s", v)

		dir := filepath.Join(tf.dir, v)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			t.Fatal(err)
		}

		path, err = tfinstall.Find(context.Background(), tfinstall.ExactVersion(v, dir))
		if err != nil {
			t.Fatalf("error installing terraform version %q: %s", v, err)
		}
		tf.versions[v] = path
	}

	return path
}
