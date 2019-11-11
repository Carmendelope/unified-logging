/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Handler for both slave and coord, implementing Search and Expire
// Implements grpc-go-unified-logging-go.SlaveServer and
// grpc-go-unified-logging-go.CoordinatorServer

package handler

import (
	"context"

	"github.com/nalej/grpc-common-go"
	grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/unified-logging/internal/pkg/managers"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	searchManager managers.Search
	expireManager managers.Expire
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
	err := validateSearch(request)
	if err != nil {
		log.Info().Str("err", err.DebugReport()).Err(err).Msg("invalid request")
		return nil, err
	}

	// Execute request on manager
	res, err := h.searchManager.Search(ctx, request)
	if err != nil {
		log.Info().Str("err", err.DebugReport()).Err(err).Msg("error executing search")
		return nil, err
	}

	return res, nil
}

// Expire the logs of a given application.
func (h *Handler) Expire(ctx context.Context, request *grpc.ExpirationRequest) (*grpc_common_go.Success, error) {
	// Validate request
	err := validateExpire(request)
	if err != nil {
		log.Info().Str("err", err.DebugReport()).Err(err).Msg("invalid request")
		return nil, err
	}

	// Execute request on manager
	res, err := h.expireManager.Expire(ctx, request)
	if err != nil {
		log.Info().Str("err", err.DebugReport()).Err(err).Msg("error executing search")
		return nil, err
	}

	return res, nil
}
