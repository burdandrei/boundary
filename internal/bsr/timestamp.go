// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bsr

import (
	"encoding/binary"
	"time"
)

const (
	secondSize     = 8
	nanosecondSize = 4

	timestampSize = secondSize + nanosecondSize
)

// Timestamp is a time.Time that can be marshaled/unmarshaled to/from a bsr data file.
// A Timestamp in a bsr data file is represented as:
//
//	uint64 seconds     8 bytes
//	uint32 nanoseconds 4 bytes
//
// Where seconds is the number of seconds since unix epoch (Jan 1, 1970 00:00:00)
// and nanoseconds are the number of nanoseconds since the last second.
// This means the BSR cannot have times earlier than unix epoch.
type Timestamp time.Time

// NewTimestamp creates a Timestamp.
func NewTimestamp(t time.Time) *Timestamp {
	tt := Timestamp(t)
	return &tt
}

func (t *Timestamp) marshal() []byte {
	tt := time.Time(*t)
	seconds := uint64(tt.Unix())
	nanoseconds := uint32(tt.Nanosecond())

	d := make([]byte, 0, timestampSize)
	d = binary.BigEndian.AppendUint64(d, seconds)
	d = binary.BigEndian.AppendUint32(d, nanoseconds)
	return d
}

// AsTime returns a time.Time for a Timestamp.
func (t *Timestamp) AsTime() time.Time {
	return time.Time(*t)
}
