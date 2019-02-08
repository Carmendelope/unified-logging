/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Fields we can filter on / search for

package entities

type Field string

// https://www.elastic.co/guide/en/beats/filebeat/current/exported-fields-kubernetes-processor.html
// https://www.elastic.co/guide/en/beats/filebeat/current/add-kubernetes-metadata.html
// https://www.elastic.co/blog/shipping-kubernetes-logs-to-elasticsearch-with-filebeat

// Should we get these from deployment-manager/pkg/utils?
const (
	NamespaceField Field = "kubernetes.namespace"
	AppInstanceIdField = "kubernetes.labels.nalej-instance"
	ServiceGroupInstanceId = "kubernetes.labels.nalej-service"
)
