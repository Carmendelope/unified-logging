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

// Fields we can filter on / search for

package entities

type Field string

// https://www.elastic.co/guide/en/beats/filebeat/current/exported-fields-kubernetes-processor.html
// https://www.elastic.co/guide/en/beats/filebeat/current/add-kubernetes-metadata.html
// https://www.elastic.co/blog/shipping-kubernetes-logs-to-elasticsearch-with-filebeat

const (
	TimestampField              Field = "@timestamp"
	NamespaceField              Field = "kubernetes.namespace"
	OrganizationIdField         Field = "kubernetes.labels." + NALEJ_ANNOTATION_ORGANIZATION_ID
	AppInstanceIdField          Field = "kubernetes.labels." + NALEJ_ANNOTATION_APP_INSTANCE_ID
	ServiceGroupInstanceIdField Field = "kubernetes.labels." + NALEJ_ANNOTATION_SERVICE_GROUP_INSTANCE_ID
	MessageField                Field = "message"
	ServiceGroupIdField         Field = "kubernetes.labels." + NALEJ_ANNOTATION_SERVICE_GROUP_ID
	ServiceIdField              Field = "kubernetes.labels." + NALEJ_ANNOTATION_SERVICE_ID
	ServiceInstanceIdField      Field = "kubernetes.labels." + NALEJ_ANNOTATION_SERVICE_INSTANCE_ID
)

func (f Field) String() string {
	return string(f)
}

type FilterFields struct {
	OrganizationId         string
	AppInstanceId          string
	ServiceGroupId         string
	ServiceGroupInstanceId string
	ServiceId              string
	ServiceInstanceId      string
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

	if f.ServiceGroupId != "" {
		filters[ServiceGroupIdField] = []string{f.ServiceGroupId}
	}

	if f.ServiceId != "" {
		filters[ServiceIdField] = []string{f.ServiceId}
	}

	if f.ServiceInstanceId != "" {
		filters[ServiceInstanceIdField] = []string{f.ServiceInstanceId}
	}

	return filters
}
