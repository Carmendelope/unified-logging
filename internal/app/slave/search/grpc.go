/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// GRPC utility functions

package search

import (
	"github.com/nalej/unified-logging/internal/pkg/utils"
	"github.com/nalej/unified-logging/pkg/entities"

	grpc "github.com/nalej/grpc-unified-logging-go"
)

func GRPCEntries(entries entities.LogEntries) []*grpc.LogEntry {
	result := make([]*grpc.LogEntry, len(entries))
	for i, e := range(entries) {
		result[i] = &grpc.LogEntry{
			Timestamp: utils.GRPCTime(e.Timestamp),
			Msg: e.Msg,
		}
	}
	return result
}
