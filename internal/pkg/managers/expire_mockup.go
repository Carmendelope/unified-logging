/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package managers

import (
	"context"

	"github.com/nalej/derrors"

	grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-common-go"
)

type MockupExpireManager struct {
}

func NewMockupExpireManager() *MockupExpireManager {
	return &MockupExpireManager{}
}

func (m *MockupExpireManager) Expire(ctx context.Context, request *grpc.ExpirationRequest) (*grpc_common_go.Success, derrors.Error) {
	return &grpc_common_go.Success{}, nil
}
