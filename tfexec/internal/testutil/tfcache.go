package testutil

import (
	"sync"
	"testing"

	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/hashicorp/terraform-exec/tfinstall/gitref"
)

const (
	Latest011   = "0.11.15"
	Latest012   = "0.12.31"
	Latest013   = "0.13.7"
	Latest014   = "0.14.11"
	Latest015   = "0.15.5"
	Latest_v1   = "1.0.0"
	Latest_v1_1 = "1.1.0-alpha20210922"
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
		return gitref.Install(ref, "", dir)
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
