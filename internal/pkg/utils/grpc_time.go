/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// GRPC timestamp utility functions
//
// The functions provided in github.com/golang/protobuf/ptypes have two
// drawbacks: they always validate, which in theory is a good thing,
// but we can't use them inline - we're going to assume that Elastic
// returns valid times; and for a nil protobuf time, they return
// a 1970-01-01 Golang time - which is not the Go convention
// (time.IsZero() checks for 0001-01-01), although the comment in the
// code states otherwise ("treat nil like the empty Timestamp")

package utils

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func GoTime(ts *timestamp.Timestamp) time.Time {
	if ts == nil {
		var t time.Time // Uninitialized time is zero
		return t
	}
	return time.Unix(ts.GetSeconds(), int64(ts.GetNanos()))
}

func GRPCTime(t time.Time) *timestamp.Timestamp {
	if t.IsZero() {
		return nil
	}
	return &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos: int32(t.Nanosecond()),
	}
}

func GRPCTimeAfter(ts1, ts2 *timestamp.Timestamp) bool {
	return GoTime(ts1).After(GoTime(ts2))
}
