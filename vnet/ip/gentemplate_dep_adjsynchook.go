// autogenerated: do not edit!
// generated from gentemplate [gentemplate -id adjSyncHook -d Package=ip -d DepsType=adjSyncCounterHookVec -d Type=adjSyncCounterHook -d Data=adjSyncCounterHooks github.com/platinasystems/go/elib/dep/dep.tmpl]

// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ip

import (
	"github.com/platinasystems/go/elib/dep"
)

type adjSyncCounterHookVec struct {
	deps                dep.Deps
	adjSyncCounterHooks []adjSyncCounterHook
}

func (t *adjSyncCounterHookVec) Len() int {
	return t.deps.Len()
}

func (t *adjSyncCounterHookVec) Get(i int) adjSyncCounterHook {
	return t.adjSyncCounterHooks[t.deps.Index(i)]
}

func (t *adjSyncCounterHookVec) Add(x adjSyncCounterHook, ds ...*dep.Dep) {
	if len(ds) == 0 {
		t.deps.Add(&dep.Dep{})
	} else {
		t.deps.Add(ds[0])
	}
	t.adjSyncCounterHooks = append(t.adjSyncCounterHooks, x)
}
