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

package entities

// copy of deployment-manager variables needed in this project
const (
	// Annotation application instance
	NALEJ_ANNOTATION_APP_INSTANCE_ID = "nalej-app-instance-id"
	// Annotation application instance
	NALEJ_ANNOTATION_APP_DESCRIPTOR = "nalej-app-descriptor"
	// Annotation for metadata to identify the group service
	NALEJ_ANNOTATION_SERVICE_GROUP_INSTANCE_ID = "nalej-service-group-instance-id"
	// Annotation for the organization
	NALEJ_ANNOTATION_ORGANIZATION_ID = "nalej-organization"
	// Annotation for metadata to identify the group
	NALEJ_ANNOTATION_SERVICE_GROUP_ID = "nalej-service-group-id"
	// Annotation for metadata to identify the service
	NALEJ_ANNOTATION_SERVICE_ID = "nalej-service-id"
	// Annotation for metadata to identify the service instance
	NALEJ_ANNOTATION_SERVICE_INSTANCE_ID = "nalej-service-instance-id"
)

const LimitPerSearch = 1