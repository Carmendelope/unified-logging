/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Search manager for unified logging slave

package search

import (
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

func (m *Manager) Search(request *grpc.SearchRequest) (*grpc.LogResponse, derrors.Error) {
	// We have a verified request - translate to entities.SearchRequest and execute
	filters := make(entities.SearchFilter)
	// TODO: fill filters

	search := &entities.SearchRequest{
		Filters: filters,
		IsUnionFilter: false,
		MsgFilter: request.GetMsgQueryFilter(),
		From: UnixTime(request.GetFrom()),
		To: UnixTime(request.GetTo()),
	}

	result, err := m.Provider.Search(search, 0 /* No limit */)
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
