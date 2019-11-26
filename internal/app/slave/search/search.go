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
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/unified-logging/pkg/entities"
	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

	"github.com/nalej/grpc-unified-logging-go"
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
		K8sIdQueryFilter: m.convertK8Ids(request.K8SIdQueryFilter),
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
		from = result[0].Timestamp.Unix()
		to = result[len(result)-1].Timestamp.Unix()

	}

	// Create GRPC response
	list :=  m.mergeLogEntries(request.OrganizationId, from, to, result)//, nil

	return list, nil
}

func (m *Manager) convertK8Ids (labels map[string]*grpc_unified_logging_go.IdList) map[string][]string {
	ids := make (map[string][]string, 0)
	for key, value := range labels {
		ids[key] = value.Ids
	}
	return ids
}


func (m *Manager) getLogEntryPK(entry entities.LogEntry) string {
	return fmt.Sprintf("%s#%s",
		entry.Kubernetes.Labels.AppInstanceId,
		entry.Kubernetes.Labels.AppServiceInstanceId,
	)
}

// mergeLogEntries group all log entries by identifiers (organizationId, appDescriptorId, AppInstanceId, etc.)
func (m *Manager) mergeLogEntries(organizationID string, from int64, to int64, entries entities.LogEntries) *grpc_unified_logging_go.LogResponseList {

	// responses is an array of responses (all messages group by serviceInstanceID)
	responses := make([]*grpc_unified_logging_go.LogResponse, 0)
	nextIndex := 0
	// aux stores the index where the messages of a serviceInstanceID are stored
	// is indexed by Instance+InstanceID
	// The reason for implementing it in this way is because we will have the logReponses stored in an array (as we have to return them)
	mapIndex := make(map[string]int, 0)
	for _, entry := range entries {
		pk := m.getLogEntryPK(*entry)

		// index is the index where the responses of this entry is stored
		index, exists := mapIndex[pk]
		if !exists {
			mapIndex[pk] = nextIndex
			index = nextIndex

			responses = append(responses, &grpc_unified_logging_go.LogResponse{
				AppDescriptorId:        entry.Kubernetes.Labels.AppDescriptorId,
				AppInstanceId:          entry.Kubernetes.Labels.AppInstanceId,
				ServiceGroupId:         entry.Kubernetes.Labels.AppServiceGroupId,
				ServiceGroupInstanceId: entry.Kubernetes.Labels.AppServiceGroupInstanceId,
				ServiceId:              entry.Kubernetes.Labels.AppServiceId,
				ServiceInstanceId:      entry.Kubernetes.Labels.AppServiceInstanceId,
				Entries:                []*grpc_unified_logging_go.LogEntry{},
			})
			// we point to the next position of the array
			nextIndex++
		}
		// add the message
		responses[index].Entries = append(responses[index].Entries, &grpc_unified_logging_go.LogEntry{
			Timestamp: entry.Timestamp.Unix(),
			Msg:       entry.Msg,
		})

	}

	return &grpc_unified_logging_go.LogResponseList{
		OrganizationId: organizationID,
		From:           from,
		To:             to,
		Responses:      responses,
	}
}