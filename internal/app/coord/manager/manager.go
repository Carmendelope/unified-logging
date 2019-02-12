/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Manager for unified logging coordinator

package manager

import (
	"context"
	"fmt"
	"sort"

        "github.com/nalej/derrors"

        "github.com/nalej/unified-logging/pkg/entities"

        "github.com/nalej/grpc-unified-logging-go"
        "github.com/nalej/grpc-application-go"
        "github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-app-cluster-api-go"
	"github.com/nalej/grpc-common-go"
)

type Manager struct {
        ApplicationsClient grpc_application_go.ApplicationsClient
        ClustersClient grpc_infrastructure_go.ClustersClient

	Executor *LoggingExecutor
}

func NewManager(apps grpc_application_go.ApplicationsClient, clusters grpc_infrastructure_go.ClustersClient, executor *LoggingExecutor) *Manager {
	return &Manager{
		ApplicationsClient: apps,
		ClustersClient: clusters,
		Executor: executor,
	}
}

func (m *Manager) GetHosts(ctx context.Context, fields *entities.FilterFields) ([]string, derrors.Error) {
	// For now we just return al hosts for an organization
	// TODO: filter out hosts for appinstanceid and servicegroupinstanceid

	org := &grpc_organization_go.OrganizationId{
		OrganizationId: fields.OrganizationId,
	}
	clusters, err := m.ClustersClient.ListClusters(ctx, org)
	if err != nil {
		return nil, derrors.NewInternalError("error getting cluster list", err)
	}

	clusterList := clusters.GetClusters()
	hosts := make([]string, len(clusterList))
	for i, cluster := range(clusterList) {
		hosts[i] = fmt.Sprintf("appcluster.%s:443", cluster.GetHostname())
	}

	return hosts, nil
}

func (m *Manager) Search(ctx context.Context, request *grpc_unified_logging_go.SearchRequest) (*grpc_unified_logging_go.LogResponse, derrors.Error) {
	// We have a verified request
        fields := &entities.FilterFields{
                OrganizationId: request.GetOrganizationId(),
                AppInstanceId: request.GetAppInstanceId(),
                ServiceGroupInstanceId: request.GetServiceGroupInstanceId(),
        }

	hosts, err := m.GetHosts(ctx, fields)
	if err != nil {
		return nil, err
	}

	out := make([][]*grpc_unified_logging_go.LogEntry, len(hosts))
	execFunc := func(ctx context.Context, client grpc_app_cluster_api_go.UnifiedLoggingClient, i int) (int, error) {
		res, err := client.Search(ctx, request)
		if err != nil {
			return 0, err
		}
		out[i] = res.GetEntries()
		return len(out[i]), nil
	}

	total, err := m.Executor.ExecRequests(ctx, hosts, execFunc)
	// TODO: Do we return some logs when we have an error, or none?
	if err != nil {
		return nil, err
	}

        // Create GRPC response
        response := &grpc_unified_logging_go.LogResponse{
                OrganizationId: request.GetOrganizationId(),
                AppInstanceId: request.GetAppInstanceId(),
                From: request.GetFrom(),
                To: request.GetTo(),
                Entries: MergeAndSort(out, total),
        }

	return response, nil
}

func (m *Manager) Expire(ctx context.Context, request *grpc_unified_logging_go.ExpirationRequest) (*grpc_common_go.Success, derrors.Error) {
	// We have a verified request
        fields := &entities.FilterFields{
                OrganizationId: request.GetOrganizationId(),
                AppInstanceId: request.GetAppInstanceId(),
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

func MergeAndSort(in [][]*grpc_unified_logging_go.LogEntry, total int) []*grpc_unified_logging_go.LogEntry {
	// Merge requests
	result := make([]*grpc_unified_logging_go.LogEntry, total)
	var count int = 0
	for _, slice := range(in) {
		if slice == nil {
			continue
		}
		count += copy(result[count:], slice)
	}

	// Sort in descending order (so less function is opposite)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.GetSeconds() > result[j].Timestamp.GetSeconds()
	})

	return result
}
