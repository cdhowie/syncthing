// Copyright (C) 2015 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package db

import (
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type VirtualMtimeRepo struct {
	ns *NamespacedKV
}

func NewVirtualMtimeRepo(ldb *leveldb.DB, folder string) *VirtualMtimeRepo {
	prefix := string(KeyTypeVirtualMtime) + folder

	return &VirtualMtimeRepo{
		ns: NewNamespacedKV(ldb, prefix),
	}
}

func (r *VirtualMtimeRepo) UpdateMtime(path string, diskMtime, actualMtime time.Time) {
	if debug {
		l.Debugf("virtual mtime: storing values for path:%v disk:%v actual:%v", path, diskMtime, actualMtime)
	}

	diskBytes, _ := diskMtime.MarshalBinary()
	actualBytes, _ := actualMtime.MarshalBinary()

	data := append(diskBytes, actualBytes...)

	r.ns.PutBytes(path, data)
}

func (r *VirtualMtimeRepo) GetMtime(path string, diskMtime time.Time) time.Time {
	var debugResult string

	if data, exists := r.ns.Bytes(path); exists {
		var mtime time.Time

		if err := mtime.UnmarshalBinary(data[:len(data)/2]); err != nil {
			panic(fmt.Sprintf("Can't unmarshal stored mtime at path %v: %v", path, err))
		}

		if mtime.Equal(diskMtime) {
			if err := mtime.UnmarshalBinary(data[len(data)/2:]); err != nil {
				panic(fmt.Sprintf("Can't unmarshal stored mtime at path %v: %v", path, err))
			}

			debugResult = "got it"
			diskMtime = mtime
		} else if debug {
			debugResult = fmt.Sprintf("record exists, but mismatch inDisk:%v dbDisk:%v", diskMtime, mtime)
		}
	} else {
		debugResult = "record does not exist"
	}

	if debug {
		l.Debugf("virtual mtime: value get result:%v path:%v", debugResult, path)
	}

	return diskMtime
}

func (r *VirtualMtimeRepo) DeleteMtime(path string) {
	r.ns.Delete(path)
}

func (r *VirtualMtimeRepo) Drop() {
	r.ns.Reset()
}
