// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !elog_gtk

package loop

import (
	"github.com/platinasystems/go/elib/elog"
)

func (l *Loop) ViewEventLog(v *elog.View) {
	l.Logf("event log graphical viewer not supported")
}
