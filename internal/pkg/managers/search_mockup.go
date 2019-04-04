/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package managers

import (
	"context"

	"github.com/nalej/derrors"

	grpc "github.com/nalej/grpc-unified-logging-go"
)

type MockupSearchManager struct {
}

func NewMockupSearchManager() *MockupSearchManager {
	return &MockupSearchManager{}
}

func (m *MockupSearchManager) Search(ctx context.Context, request *grpc.SearchRequest) (*grpc.LogResponse, derrors.Error) {
	response := &grpc.LogResponse{
		OrganizationId: request.GetOrganizationId(),
		AppInstanceId: request.GetAppInstanceId(),
		From: request.GetFrom(),
		To: request.GetTo(),
		Entries: []*grpc.LogEntry{},
	}
	return response, nil
}
