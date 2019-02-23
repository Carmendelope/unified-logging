/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// GRPC utility functions

package search

import (
	"time"

	"github.com/nalej/unified-logging/pkg/entities"

        grpc "github.com/nalej/grpc-unified-logging-go"
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

func GRPCEntries(entries entities.LogEntries) []*grpc.LogEntry {
	result := make([]*grpc.LogEntry, len(entries))
	for i, e := range(entries) {
		result[i] = &grpc.LogEntry{
			Timestamp: GRPCTime(e.Timestamp),
			Msg: e.Msg,
		}
	}
	return result
}
