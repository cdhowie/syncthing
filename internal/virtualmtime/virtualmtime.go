// Copyright (C) 2015 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package virtualmtime

import (
	"time"
)

type VirtualMtimeRepo interface {
	UpdateMtime(path string, diskMtime, actualMtime time.Time)
	GetMtime(path string, diskMtime time.Time) time.Time
	DeleteMtime(path string)
	Drop()
}
