/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Fields we can filter on / search for

package entities

import (
	"github.com/nalej/deployment-manager/pkg/utils"
	"github.com/nalej/deployment-manager/pkg/common"
)

type Field string

// https://www.elastic.co/guide/en/beats/filebeat/current/exported-fields-kubernetes-processor.html
// https://www.elastic.co/guide/en/beats/filebeat/current/add-kubernetes-metadata.html
// https://www.elastic.co/blog/shipping-kubernetes-logs-to-elasticsearch-with-filebeat

// Should we get these from deployment-manager/pkg/utils?
const (
        TimestampField = "@timestamp"
	NamespaceField Field = "kubernetes.namespace"
	OrganizationIdField = "kubernetes.labels." + "nalej-organization" // TODO
	AppInstanceIdField = "kubernetes.labels." + utils.NALEJ_ANNOTATION_INSTANCE_ID
	ServiceGroupInstanceIdField = "kubernetes.labels." + "nalej-service-group" //TODO
        MessageField = "message"
)

func (f Field) String() string {
	return string(f)
}

type FilterFields struct {
	OrganizationId string
	AppInstanceId string
	ServiceGroupInstanceId string
}

func (f *FilterFields) ToFilters() SearchFilter {
	filters := make(SearchFilter)

	// Is this needed? Does this speed up or slow down?
	if f.OrganizationId != "" && f.AppInstanceId != "" {
		filters[NamespaceField] = []string{common.GetNamespace(f.OrganizationId, f.AppInstanceId)}
	}

	if f.OrganizationId != "" {
		filters[OrganizationIdField] = []string{f.OrganizationId}
	}

	if f.AppInstanceId != "" {
		filters[AppInstanceIdField] = []string{f.AppInstanceId}
	}

	if f.ServiceGroupInstanceId != "" {
		filters[ServiceGroupInstanceIdField] = []string{f.ServiceGroupInstanceId}
	}

	return filters
}
