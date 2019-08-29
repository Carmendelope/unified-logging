/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// GRPC utility functions

package search

import (
	grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/unified-logging/pkg/entities"
)

func GRPCEntries(entries entities.LogEntries) []*grpc.LogEntry {
	result := make([]*grpc.LogEntry, len(entries))
	for i, e := range(entries) {
		result[i] = &grpc.LogEntry{
			Timestamp: conversions.GRPCTime(e.Timestamp),
			Msg: e.Msg,
		}
	}
	return result
}
