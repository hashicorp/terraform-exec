package testutil

import (
	"sync"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

const (
	Latest011 = "0.11.14"
	Latest012 = "0.12.29"
	Latest013 = "0.13.0"
)

const appendUserAgent = "tfexec-testutil"

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
		f := tfinstall.ExactVersion(v, dir)
		f.UserAgent = appendUserAgent
		return f
	})
}
