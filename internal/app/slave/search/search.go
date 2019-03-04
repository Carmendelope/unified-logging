/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Search manager for unified logging slave

package search

import (
	"context"
	"time"

        "github.com/nalej/derrors"

	"github.com/nalej/unified-logging/pkg/entities"
	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

        grpc "github.com/nalej/grpc-unified-logging-go"
)

type Manager struct {
        Provider loggingstorage.Provider
}

func NewManager(provider loggingstorage.Provider) *Manager {
	return &Manager{
		Provider: provider,
	}
}

func (m *Manager) Search(ctx context.Context, request *grpc.SearchRequest) (*grpc.LogResponse, derrors.Error) {
	// We have a verified request - translate to entities.SearchRequest and execute
	fields := entities.FilterFields{
		OrganizationId: request.GetOrganizationId(),
		AppInstanceId: request.GetAppInstanceId(),
		ServiceGroupInstanceId: request.GetServiceGroupInstanceId(),
	}

	search := &entities.SearchRequest{
		Filters: fields.ToFilters(),
		IsUnionFilter: false,
		MsgFilter: request.GetMsgQueryFilter(),
		From: GoTime(request.GetFrom()),
		To: GoTime(request.GetTo()),
		Order: entities.SortOrder(request.GetOrder()),
	}

	result, err := m.Provider.Search(ctx, search, -1 /* No limit */)
	if err != nil {
		return nil, err
	}

	// Assuming the entries are sorted, we can get the timestamp of
	// the first and last entry to get the whole range
	var from, to time.Time
	if len(result) > 0 {
		from = result[0].Timestamp
		to = result[len(result)-1].Timestamp

		// Make from/to determination independent of sort order
		if from.After(to) {
			tmp := from
			from = to
			to = tmp
		}
	}

	// Create GRPC response
	response := &grpc.LogResponse{
		OrganizationId: request.GetOrganizationId(),
		AppInstanceId: request.GetAppInstanceId(),
		From: GRPCTime(from),
		To: GRPCTime(to),
		Entries: GRPCEntries(result),
	}
	return response, nil
}
