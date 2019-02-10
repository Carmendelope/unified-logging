/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Expire manager for unified logging slave

package expire

import (
	"context"

        "github.com/nalej/derrors"

	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

        grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-common-go"
)

type Manager struct {
        Provider loggingstorage.Provider
}

func NewManager(provider loggingstorage.Provider) *Manager {
	return &Manager{
		Provider: provider,
	}
}

func (m *Manager) Expire(ctx context.Context, request *grpc.ExpirationRequest) (*grpc_common_go.Success, derrors.Error) {
        // We have a verified request - translate to entities.SearchRequest and execute
        filters := make(entities.SearchFilter)
        // TODO: fill filters

        search := &entities.SearchRequest{
                Filters: filters,
                IsUnionFilter: false,
        }

        err := m.Provider.Expire(ctx, search)
        if err != nil {
                return nil, err
        }

        return *grpc_common_go.Success{}, nil
}
