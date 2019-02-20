/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package entities

import (
	"fmt"
        "github.com/onsi/ginkgo"
        "github.com/onsi/gomega"
)

const (
        OrganizationId = "2a95fe95-eade-4622-836f-e85d789024bf"
        AppInstanceId = "e9e38334-1da1-4f51-8f18-2bd8e2470123"
        ServiceGroupInstanceId = "44eb6008-e288-47ea-bc6c-44a7af56df51"
)

var _ = ginkgo.Describe("Fields", func() {
	ginkgo.Context("Filters", func() {
		ginkgo.It("should create filters from fields", func() {
			var fields = &FilterFields{
				OrganizationId: OrganizationId,
				AppInstanceId: AppInstanceId,
				ServiceGroupInstanceId: ServiceGroupInstanceId,
			}
			filter := fields.ToFilters()
			gomega.Expect(filter).Should(gomega.BeEquivalentTo(SearchFilter{
				"kubernetes.namespace": []string{fmt.Sprintf("%s-%s", OrganizationId, AppInstanceId)[:63]},
				"kubernetes.labels.nalej-organization": []string{OrganizationId},
				"kubernetes.labels.nalej-app-instance-id": []string{AppInstanceId},
				"kubernetes.labels.nalej-service-group-instance-id": []string{ServiceGroupInstanceId},
			}))
		})
	})
})
