/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Search manager for unified logging slave

package search

import (
	"context"

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
	}

	result, err := m.Provider.Search(ctx, search, 0 /* No limit */)
	if err != nil {
		return nil, err
	}

	// Create GRPC response
	response := &grpc.LogResponse{
		OrganizationId: request.GetOrganizationId(),
		AppInstanceId: request.GetAppInstanceId(),
		From: request.GetFrom(),
		To: request.GetTo(),
		Entries: GRPCEntries(result),
	}
	return response, nil
}
