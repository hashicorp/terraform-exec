// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testutil

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/build"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

const (
	Latest011          = "0.11.15"
	Latest012          = "0.12.31"
	Latest013          = "0.13.7"
	Latest014          = "0.14.11"
	Latest015          = "0.15.5"
	Latest_v1          = "1.0.11"
	Latest_v1_1        = "1.1.9"
	Latest_v1_2        = "1.2.9"
	Latest_v1_3        = "1.3.10"
	Latest_v1_4        = "1.4.7"
	Latest_v1_5        = "1.5.3"
	Latest_v1_6        = "1.6.6"
	Latest_v1_7        = "1.7.5"
	Latest_v1_8        = "1.8.5"
	Latest_Beta_v1_8   = "1.8.0-beta1"
	Latest_v1_9        = "1.9.8"
	Latest_Alpha_v1_9  = "1.9.0-alpha20240516"
	Latest_v1_10       = "1.10.5"
	Latest_Alpha_v1_10 = "1.10.0-alpha20240926"
	Latest_v1_11       = "1.11.4"
	Latest_v1_12       = "1.12.2"
	Latest_Alpha_v1_14 = "1.14.0-alpha20250903"
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

	key := "gitref:" + ref

	return tf.find(t, key, func(ctx context.Context) (string, error) {
		gr := &build.GitRevision{
			Product: product.Terraform,
			Ref:     ref,
		}
		gr.SetLogger(TestLogger())

		return gr.Build(ctx)
	})
}

func (tf *TFCache) Version(t *testing.T, v string) string {
	t.Helper()

	key := "v:" + v

	return tf.find(t, key, func(ctx context.Context) (string, error) {
		ev := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion(v)),
		}
		ev.SetLogger(TestLogger())

		return ev.Install(ctx)
	})
}
