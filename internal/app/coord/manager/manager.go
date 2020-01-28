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

// Manager for unified logging coordinator

package manager

import (
	"context"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-connectivity-manager-go"
	"github.com/nalej/unified-logging/pkg/entities"
	"github.com/rs/zerolog/log"
	"sort"
	"time"

	"github.com/nalej/grpc-app-cluster-api-go"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
	"github.com/nalej/grpc-unified-logging-go"
)

const defaultLogExpiration = 7

type Manager struct {
	ApplicationsClient grpc_application_go.ApplicationsClient
	ClustersClient     grpc_infrastructure_go.ClustersClient
	OrgClient          grpc_organization_manager_go.OrganizationsClient
	Executor           *LoggingExecutor

	appClusterPrefix string
	appClusterPort   int
}

func NewManager(apps grpc_application_go.ApplicationsClient, clusters grpc_infrastructure_go.ClustersClient, executor *LoggingExecutor, prefix string, port int) *Manager {
	return &Manager{
		ApplicationsClient: apps,
		ClustersClient:     clusters,
		Executor:           executor,
		appClusterPrefix:   prefix,
		appClusterPort:     port,
	}
}

type ClusterInfo struct {
	host string
	id   string
}

func (m *Manager) GetHosts(ctx context.Context, fields *entities.FilterFields) ([]ClusterInfo, derrors.Error) {
	// For now we just return all hosts for an organization
	// TODO: filter out hosts for appinstanceid, servicegroupinstanceid, servicegroupid, serviceId, serviceinstanceid

	org := &grpc_organization_go.OrganizationId{
		OrganizationId: fields.OrganizationId,
	}
	clusters, err := m.ClustersClient.ListClusters(ctx, org)
	if err != nil {
		return nil, derrors.NewInternalError("error getting cluster list", err)
	}

	prefix := m.appClusterPrefix
	if prefix != "" {
		prefix = prefix + "."
	}

	clusterList := clusters.GetClusters()
	hosts := make([]ClusterInfo, 0)
	for _, cluster := range clusterList {
		if cluster.ClusterStatus != grpc_connectivity_manager_go.ClusterStatus_OFFLINE && cluster.ClusterStatus != grpc_connectivity_manager_go.ClusterStatus_OFFLINE_CORDON {
			host := fmt.Sprintf("%s%s:%d", prefix, cluster.GetHostname(), m.appClusterPort)
			hosts = append(hosts, ClusterInfo{host, cluster.ClusterId})
		}
	}

	return hosts, nil
}

// Search method that sends a Search message to all the clusters (logging-slave)
// TODO: the slaves returns a ReponseList. The ccoordinator has to convert this into an array log entries, order all the messages by timestamp and group again by identifiers.
// we should change the slaves so that they return an array of logs
func (m *Manager) Search(ctx context.Context, request *grpc_unified_logging_go.SearchRequest) (*grpc_unified_logging_go.LogResponseList, derrors.Error) {

	// We have a verified request
	fields := &entities.FilterFields{
		OrganizationId:         request.GetOrganizationId(),
		AppDescriptorId:        request.GetAppDescriptorId(),
		AppInstanceId:          request.GetAppInstanceId(),
		ServiceGroupInstanceId: request.GetServiceGroupInstanceId(),
		ServiceGroupId:         request.ServiceGroupId,
		ServiceId:              request.ServiceId,
		ServiceInstanceId:      request.ServiceInstanceId,
	}

	hosts, err := m.GetHosts(ctx, fields)
	if err != nil {
		return nil, err
	}

	// TODO: call to slave in different threads
	out := make([]*grpc_unified_logging_go.LogResponseList, len(hosts))

	execFunc := func(ctx context.Context, client grpc_app_cluster_api_go.UnifiedLoggingClient, i int) (int, error) {
		res, err := client.Search(ctx, request)
		if err != nil {
			return 0, err
		}
		out[i] = res
		return len(out[i].Responses), nil
	}

	total, errorIds, err := m.Executor.ExecRequests(ctx, hosts, execFunc)
	// TODO: Do we return some logs when we have an error, or none?
	if err != nil {
		return nil, err
	}

	return m.mergeAllResponses(out, total, request, errorIds), nil
}

func (m *Manager) mergeAllResponses(lists []*grpc_unified_logging_go.LogResponseList, total int, request *grpc_unified_logging_go.SearchRequest, errorIds []string) *grpc_unified_logging_go.LogResponseList {
	// we need to get only the last limitPerSearch entry logs.
	// 1) convert LogResponseList in []LogEntry
	// 2) order by timestamp
	// 3) get last or first limitPerSearch Log entries
	// 4) and convert into LogResponseList again

	// 1)
	logEntries := make([]*entities.LogEntry, 0)
	for _, logResponseList := range lists {
		// if one of the slaves returns an error, logResponseList can be nil
		if logResponseList == nil {
			continue
		}
		for _, logResponse := range logResponseList.Responses {
			for _, entry := range logResponse.Entries {
				logEntries = append(logEntries, &entities.LogEntry{
					Timestamp: time.Unix(0, entry.Timestamp),
					Msg:       entry.Msg,
					Kubernetes: entities.KubernetesEntry{
						Labels: entities.KubernetesLabelsEntry{
							OrganizationId:            logResponseList.OrganizationId,
							AppDescriptorId:           logResponse.AppDescriptorId,
							AppDescriptorName:         logResponse.AppDescriptorName,
							AppInstanceId:             logResponse.AppInstanceId,
							AppInstanceName:           logResponse.AppInstanceName,
							AppServiceGroupId:         logResponse.ServiceGroupId,
							AppServiceGroupName:       logResponse.ServiceGroupName,
							AppServiceGroupInstanceId: logResponse.ServiceGroupInstanceId,
							AppServiceId:              logResponse.ServiceId,
							AppServiceName:            logResponse.ServiceName,
							AppServiceInstanceId:      logResponse.ServiceInstanceId,
						},
					},
				})
			}
		}
	}
	// 2)
	sort.SliceStable(logEntries, func(i, j int) bool {
		return logEntries[i].Timestamp.Before(logEntries[j].Timestamp)
	})

	// 3)
	if len(logEntries) > entities.LimitPerSearch {
		if request.NFirst {
			logEntries = logEntries[0:entities.LimitPerSearch]
		} else {
			logEntries = logEntries[len(logEntries)-entities.LimitPerSearch-1 : entities.LimitPerSearch]
		}
	}

	// 4)
	var from, to int64
	from = request.From
	to = request.To

	if len(logEntries) > 0 {
		from = logEntries[len(logEntries)-1].Timestamp.UnixNano()
		to = logEntries[0].Timestamp.UnixNano()

		if from > to {
			aux := to
			to = from
			from = aux
		}
	}

	list := entities.MergeLogEntries(request.OrganizationId, from, to, logEntries, errorIds)
	return list

}

func (m *Manager) Expire(ctx context.Context, request *grpc_unified_logging_go.ExpirationRequest) (*grpc_common_go.Success, derrors.Error) {
	// We have a verified request
	fields := &entities.FilterFields{
		OrganizationId: request.GetOrganizationId(),
		AppInstanceId:  request.GetAppInstanceId(),
	}

	hosts, err := m.GetHosts(ctx, fields)
	if err != nil {
		return nil, err
	}

	execFunc := func(ctx context.Context, client grpc_app_cluster_api_go.UnifiedLoggingClient, i int) (int, error) {
		_, err := client.Expire(ctx, request)
		return 0, err
	}
	_, errorIds, err := m.Executor.ExecRequests(ctx, hosts, execFunc)
	// Even with error we'll have expired something maybe - what do we do here?
	if err != nil {
		return nil, err
	}

	log.Debug().Interface("errors", errorIds).Msg("errors in search")

	return &grpc_common_go.Success{}, nil
}
