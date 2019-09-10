/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Fields we can filter on / search for

package entities

type Field string

// https://www.elastic.co/guide/en/beats/filebeat/current/exported-fields-kubernetes-processor.html
// https://www.elastic.co/guide/en/beats/filebeat/current/add-kubernetes-metadata.html
// https://www.elastic.co/blog/shipping-kubernetes-logs-to-elasticsearch-with-filebeat

const (
	TimestampField Field = "@timestamp"
	NamespaceField Field = "kubernetes.namespace"
	OrganizationIdField Field = "kubernetes.labels." + NALEJ_ANNOTATION_ORGANIZATION_ID
	AppInstanceIdField Field = "kubernetes.labels." + NALEJ_ANNOTATION_APP_INSTANCE_ID
	ServiceGroupInstanceIdField Field = "kubernetes.labels." + NALEJ_ANNOTATION_SERVICE_GROUP_INSTANCE_ID
	MessageField Field = "message"
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
