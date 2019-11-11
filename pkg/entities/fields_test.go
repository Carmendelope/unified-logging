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

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	OrganizationId         = "2a95fe95-eade-4622-836f-e85d789024bf"
	AppInstanceId          = "e9e38334-1da1-4f51-8f18-2bd8e2470123"
	ServiceGroupInstanceId = "44eb6008-e288-47ea-bc6c-44a7af56df51"
)

var _ = ginkgo.Describe("Fields", func() {
	ginkgo.Context("Filters", func() {
		ginkgo.It("should create filters from fields", func() {
			var fields = &FilterFields{
				OrganizationId:         OrganizationId,
				AppInstanceId:          AppInstanceId,
				ServiceGroupInstanceId: ServiceGroupInstanceId,
			}
			filter := fields.ToFilters()
			gomega.Expect(filter).Should(gomega.BeEquivalentTo(SearchFilter{
				OrganizationIdField:         []string{OrganizationId},
				AppInstanceIdField:          []string{AppInstanceId},
				ServiceGroupInstanceIdField: []string{ServiceGroupInstanceId},
			}))
		})
	})
})
