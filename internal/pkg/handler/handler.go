/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Handler for both slave and coord, implementing Search and Expire
// Implements grpc-go-unified-logging-go.SlaveServer and
// grpc-go-unified-logging-go.CoordinatorServer

package handler

import (
	"context"

	"github.com/nalej/unified-logging/internal/pkg/managers"
	"github.com/nalej/grpc-common-go"
	grpc "github.com/nalej/grpc-unified-logging-go"
)

type Handler struct {
	searchManager	managers.Search
	expireManager	managers.Expire
}

func NewHandler(search managers.Search, expire managers.Expire) *Handler {
	return &Handler{
		searchManager: search,
		expireManager: expire,
	}
}

// Search for log entries matching a query.
func (h *Handler) Search(ctx context.Context, request *grpc.SearchRequest) (*grpc.LogResponse, error) {
	// Validate request
	// TBD

	// Execute request on manager
	// TBD

	return nil, nil
}

// Expire the logs of a given application.
func (h *Handler) Expire(ctx context.Context, requesst *grpc.ExpirationRequest) (*grpc_common_go.Success, error) {
	// Validate request
	// TBD

	// Execute request on manager
	// TBD

	return nil, nil
}
