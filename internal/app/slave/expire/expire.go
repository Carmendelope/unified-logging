/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Expire manager for unified logging slave

package expire

import (
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

func (m *Manager) Expire(*grpc.ExpirationRequest) (*grpc_common_go.Success, derrors.Error) {
	return nil, nil
}
