/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Search manager for unified logging slave

package search

import (
        "github.com/nalej/derrors"

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

func (m *Manager) Search(*grpc.SearchRequest) (*grpc.LogResponse, derrors.Error) {
	return nil, nil
}
