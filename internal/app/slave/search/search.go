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

// Search manager for unified logging slave

package search

import (
	"context"
	"time"

	"github.com/nalej/derrors"

	"github.com/nalej/grpc-utils/pkg/conversions"
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
		OrganizationId:         request.GetOrganizationId(),
		AppInstanceId:          request.GetAppInstanceId(),
		ServiceGroupInstanceId: request.GetServiceGroupInstanceId(),
	}

	search := &entities.SearchRequest{
		Filters:       fields.ToFilters(),
		IsUnionFilter: false,
		MsgFilter:     request.GetMsgQueryFilter(),
		From:          conversions.GoTime(request.GetFrom()),
		To:            conversions.GoTime(request.GetTo()),
		Order:         entities.SortOrder(request.GetOrder()),
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
		AppInstanceId:  request.GetAppInstanceId(),
		From:           conversions.GRPCTime(from),
		To:             conversions.GRPCTime(to),
		Entries:        GRPCEntries(result),
	}
	return response, nil
}
