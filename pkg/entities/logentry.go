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
	"time"
)

type 	LogEntries []*LogEntry

type KubernetesLabelsEntry struct {
	OrganizationId            string `json:"nalej-organization"`
	AppDescriptorId           string `json:"nalej-app-descriptor"`
	AppInstanceId             string `json:"nalej-app-instance-id"`
	AppServiceGroupId         string `json:"nalej-service-group-id"`
	AppServiceGroupInstanceId string `json:"nalej-service-group-instance-id"`
	AppServiceId              string `json:"nalej-service-id"`
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
