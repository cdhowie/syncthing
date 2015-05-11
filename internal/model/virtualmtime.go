// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"encoding/json"
	"time"

	"github.com/syncthing/syncthing/internal/db"
	"github.com/syndtr/goleveldb/leveldb"
)

type virtualMtimeRepo struct {
	ns *db.NamespacedKV
}

type virtualMtimeRecord struct {
	DiskMtime   time.Time
	ActualMtime time.Time
}

func newVirtualMtimeRepo(ldb *leveldb.DB, folder string) *virtualMtimeRepo {
	prefix := string(db.KeyTypeVirtualMtime) + folder

	return &virtualMtimeRepo{
		ns: db.NewNamespacedKV(ldb, prefix),
	}
}

func (r *virtualMtimeRepo) UpdateMtime(path string, diskMtime, actualMtime time.Time) {
	// TODO need error handling?
	data, _ := json.Marshal(virtualMtimeRecord{
		DiskMtime:   diskMtime,
		ActualMtime: actualMtime,
	})

	r.ns.PutString(path, string(data))
}

func (r *virtualMtimeRepo) GetMtime(path string, diskMtime time.Time) time.Time {
	data, exists := r.ns.String(path)

	if exists {
		record := virtualMtimeRecord{}

		err := json.Unmarshal([]byte(data), &record)

		if err == nil && record.DiskMtime.Equal(diskMtime) {
			diskMtime = record.ActualMtime
		}
	}

	return diskMtime
}

func (r *virtualMtimeRepo) DeleteMtime(path string) {
	r.ns.Delete(path)
}
