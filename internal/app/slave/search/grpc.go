/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// GRPC utility functions

package search

import (
	"github.com/nalej/unified-logging/pkg/entities"

        grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func UnixTime(ts *timestamp.Timestamp) int64 {
	if ts == nil {
		return 0
	}
	return ts.GetSeconds()
}

func GRPCTime(unix int64) *timestamp.Timestamp {
	if unix == 0 {
		return nil
	}
	return &timestamp.Timestamp{
		Seconds: unix,
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
