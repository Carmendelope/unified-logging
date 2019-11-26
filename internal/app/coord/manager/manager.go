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
	"time"

	"github.com/nalej/grpc-app-cluster-api-go"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-unified-logging-go"
)

type Manager struct {
	ApplicationsClient grpc_application_go.ApplicationsClient
	ClustersClient     grpc_infrastructure_go.ClustersClient

	Executor *LoggingExecutor

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

func (m *Manager) GetHosts(ctx context.Context, fields *entities.FilterFields) ([]string, derrors.Error) {
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
	hosts := make([]string, 0)
	for _, cluster := range clusterList {
		if cluster.ClusterStatus != grpc_connectivity_manager_go.ClusterStatus_OFFLINE && cluster.ClusterStatus != grpc_connectivity_manager_go.ClusterStatus_OFFLINE_CORDON {
			hosts = append(hosts, fmt.Sprintf("%s%s:%d", prefix, cluster.GetHostname(), m.appClusterPort))
		}
	}

	return hosts, nil
}

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

	// if we don't have date range, we ask for the last hour
	if request.From == 0 && request.To == 0 {
		oneHourAgo := time.Now().Add(time.Hour * (-1))
		request.From = oneHourAgo.Unix()
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

	total, err := m.Executor.ExecRequests(ctx, hosts, execFunc)
	// TODO: Do we return some logs when we have an error, or none?
	if err != nil {
		return nil, err
	}

	return m.mergeAllResponses(out, total, request), nil
}

func (m *Manager) mergeAllResponses(lists []*grpc_unified_logging_go.LogResponseList, total int, request *grpc_unified_logging_go.SearchRequest) *grpc_unified_logging_go.LogResponseList {

	// store all the responses in an array
	// update the time-range. (Min From and MAx to)
	result := make([]*grpc_unified_logging_go.LogResponse, total)
	count := 0
	maxTo := request.To
	minFrom := request.From
	for _, responses := range lists {
		if responses == nil || responses.Responses == nil {
			continue
		}
		count += copy(result[count:], responses.Responses)
		if responses.From < minFrom {
			minFrom = responses.From
		}
		if responses.To > maxTo {
			maxTo = responses.To
		}
	}

	return &grpc_unified_logging_go.LogResponseList{
		OrganizationId: request.OrganizationId,
		From: minFrom,
		To: maxTo,
		Responses: result,
	}

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

	_, err = m.Executor.ExecRequests(ctx, hosts, execFunc)
	// Even with error we'll have expired something maybe - what do we do here?
	if err != nil {
		return nil, err
	}

	return &grpc_common_go.Success{}, nil
}

