// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"fmt"
	"time"

	"github.com/syncthing/syncthing/internal/db"
	"github.com/syndtr/goleveldb/leveldb"
)

type virtualMtimeRepo struct {
	ns *db.NamespacedKV
}

func newVirtualMtimeRepo(ldb *leveldb.DB, folder string) *virtualMtimeRepo {
	prefix := string(db.KeyTypeVirtualMtime) + folder

	return &virtualMtimeRepo{
		ns: db.NewNamespacedKV(ldb, prefix),
	}
}

func (r *virtualMtimeRepo) UpdateMtime(path string, diskMtime, actualMtime time.Time) {
	diskBytes, _ := diskMtime.MarshalBinary()
	actualBytes, _ := actualMtime.MarshalBinary()

	data := append(diskBytes, actualBytes...)

	r.ns.PutBytes(path, data)
}

func (r *virtualMtimeRepo) GetMtime(path string, diskMtime time.Time) time.Time {
	data, exists := r.ns.String(path)

	if exists {
		var mtime time.Time

		err := mtime.UnmarshalBinary([]byte(data[:len(data)/2]))
		if err != nil {
			panic(fmt.Sprintf("Can't unmarshal stored mtime at path %v: %v", path, err))
		}

		if mtime.Equal(diskMtime) {
			err := mtime.UnmarshalBinary([]byte(data[len(data)/2:]))
			if err != nil {
				panic(fmt.Sprintf("Can't unmarshal stored mtime at path %v: %v", path, err))
			}

			diskMtime = mtime
		}
	}

	return diskMtime
}

func (r *virtualMtimeRepo) DeleteMtime(path string) {
	r.ns.Delete(path)
}
