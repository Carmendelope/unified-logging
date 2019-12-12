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
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/unified-logging/pkg/entities"
	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"
)

type Manager struct {
	Provider loggingstorage.Provider
}

func NewManager(provider loggingstorage.Provider) *Manager {
	return &Manager{
		Provider: provider,
	}
}

func (m *Manager) Search(ctx context.Context, request *grpc_unified_logging_go.SearchRequest) (*grpc_unified_logging_go.LogResponseList, derrors.Error) {

	// We have a verified request - translate to entities.SearchRequest and execute
	fields := entities.FilterFields{
		OrganizationId:         request.GetOrganizationId(),
		AppDescriptorId:        request.GetAppDescriptorId(),
		AppInstanceId:          request.GetAppInstanceId(),
		ServiceGroupId:         request.ServiceGroupId,
		ServiceGroupInstanceId: request.GetServiceGroupInstanceId(),
		ServiceId:              request.ServiceId,
		ServiceInstanceId:      request.ServiceInstanceId,
	}

	search := &entities.SearchRequest{
		Filters:       fields.ToFilters(),
		IsUnionFilter: true,
		MsgFilter:     request.GetMsgQueryFilter(),
		From:          request.From,
		To:            request.To,
		NFirst:        request.NFirst,
	}

	result, err := m.Provider.Search(ctx, search, -1 /* No limit */)
	if err != nil {
		return nil, err
	}

	// Assuming the entries are sorted, we can get the timestamp of
	// the first and last entry to get the whole range
	from := request.From
	to := request.To
	if len(result) > 0 {
		from = result[0].Timestamp.UnixNano()
		to = result[len(result)-1].Timestamp.UnixNano()

	}

	// Create GRPC response
	list := entities.MergeLogEntries(request.OrganizationId, from, to, result, []string{""})

	return list, nil
}
