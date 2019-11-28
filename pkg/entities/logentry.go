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

// Log-line entries for logging provider

package entities

import (
	"fmt"
	"github.com/nalej/grpc-unified-logging-go"
	"time"
)

type LogEntries []*LogEntry

type KubernetesLabelsEntry struct {
	OrganizationId            string `json:"nalej-organization"`
	AppDescriptorId           string `json:"nalej-app-descriptor"`
	AppDescriptorName         string `json:"nalej-app-descriptor-name"`
	AppInstanceId             string `json:"nalej-app-instance-id"`
	AppInstanceName           string `json:"nalej-app-name"`
	AppServiceGroupId         string `json:"nalej-service-group-id"`
	AppServiceGroupName       string `json:"nalej-service-group-name"`
	AppServiceGroupInstanceId string `json:"nalej-service-group-instance-id"`
	AppServiceId              string `json:"nalej-service-id"`
	AppServiceName            string `json:"nalej-service-name"`
	AppServiceInstanceId      string `json:"nalej-service-instance-id"`
}

type KubernetesEntry struct {
	Namespace string                `json:"namespace"`
	Labels    KubernetesLabelsEntry `json:"labels"`
}

type LogEntry struct {
	Timestamp  time.Time       `json:"@timestamp"`
	Msg        string          `json:"message"`
	Kubernetes KubernetesEntry `json:"kubernetes"`
}

func  getLogEntryPK(entry LogEntry) string {
	return fmt.Sprintf("%s#%s",
		entry.Kubernetes.Labels.AppInstanceId,
		entry.Kubernetes.Labels.AppServiceInstanceId,
	)
}

// mergeLogEntries group all log entries by identifiers (organizationId, appDescriptorId, AppInstanceId, etc.)
func MergeLogEntries(organizationID string, from int64, to int64, entries LogEntries) *grpc_unified_logging_go.LogResponseList {

	// responses is an array of responses (all messages group by serviceInstanceID)
	responses := make([]*grpc_unified_logging_go.LogResponse, 0)
	nextIndex := 0
	// aux stores the index where the messages of a serviceInstanceID are stored
	// is indexed by Instance+InstanceID
	// The reason for implementing it in this way is because we will have the logReponses stored in an array (as we have to return them)
	mapIndex := make(map[string]int, 0)
	for _, entry := range entries {
		pk := getLogEntryPK(*entry)

		// index is the index where the responses of this entry is stored
		index, exists := mapIndex[pk]
		if !exists {
			mapIndex[pk] = nextIndex
			index = nextIndex

			responses = append(responses, &grpc_unified_logging_go.LogResponse{
				AppDescriptorId:        entry.Kubernetes.Labels.AppDescriptorId,
				AppDescriptorName:      entry.Kubernetes.Labels.AppDescriptorName,
				AppInstanceId:          entry.Kubernetes.Labels.AppInstanceId,
				AppInstanceName:        entry.Kubernetes.Labels.AppInstanceName,
				ServiceGroupId:         entry.Kubernetes.Labels.AppServiceGroupId,
				ServiceGroupName:       entry.Kubernetes.Labels.AppServiceGroupName,
				ServiceGroupInstanceId: entry.Kubernetes.Labels.AppServiceGroupInstanceId,
				ServiceId:              entry.Kubernetes.Labels.AppServiceId,
				ServiceName:            entry.Kubernetes.Labels.AppServiceName,
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