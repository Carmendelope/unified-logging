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

package handler

import (
	"time"

	"github.com/nalej/grpc-unified-logging-go"
)

const (
	OrganizationId         = "2a95fe95-eade-4622-836f-e85d789024bf"
	AppInstanceId          = "e9e38334-1da1-4f51-8f18-2bd8e2470123"
	ServiceGroupInstanceId = "44eb6008-e288-47ea-bc6c-44a7af56df51"
	MsgQueryFilter         = "random message"
)

var (
	From = 0
	To   = time.Now().UnixNano()
)

var ValidSearchRequest = &grpc_unified_logging_go.SearchRequest{
	OrganizationId:         OrganizationId,
	AppInstanceId:          AppInstanceId,
	ServiceGroupInstanceId: ServiceGroupInstanceId,
	MsgQueryFilter:         MsgQueryFilter,
	From:                   0,
	To:                     To,
}

var ValidExpirationRequest = &grpc_unified_logging_go.ExpirationRequest{
	OrganizationId: OrganizationId,
	AppInstanceId:  AppInstanceId,
}

/*
var _ = ginkgo.Describe("Handler", func() {
	// const numServices = 2

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener

	// clients
	var coordClient grpc_unified_logging_go.CoordinatorClient
	var slaveClient grpc_unified_logging_go.SlaveClient

	// Managers
	var searchManager managers.Search
	var expireManager managers.Expire

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Create managers
		searchManager = managers.NewMockupSearchManager()
		expireManager = managers.NewMockupExpireManager()

		handler := NewHandler(searchManager, expireManager)
		grpc_unified_logging_go.RegisterCoordinatorServer(server, handler)
		grpc_unified_logging_go.RegisterSlaveServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		coordClient = grpc_unified_logging_go.NewCoordinatorClient(conn)
		slaveClient = grpc_unified_logging_go.NewSlaveClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.Context("Coordinator client", func() {
		ginkgo.Context("Search handler", func() {
			ginkgo.It("should reject request without organization", func() {
				req := &grpc_unified_logging_go.SearchRequest{}
				res, err := coordClient.Search(context.Background(), req)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(res).Should(gomega.BeNil())
			})
			ginkgo.It("should return requested values in result", func() {
				res, err := coordClient.Search(context.Background(), ValidSearchRequest)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(res.GetOrganizationId()).Should(gomega.Equal(OrganizationId))
				gomega.Expect(res.GetAppInstanceId()).Should(gomega.Equal(AppInstanceId))
				//gomega.Expect(*res.GetFrom()).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				//	"Seconds": gomega.Equal(int64(0)),
				//	"Nanos":   gomega.Equal(int32(0)),
				//}))
				//gomega.Expect(*res.GetTo()).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				//	"Seconds": gomega.Equal(To.Seconds),
				//	"Nanos":   gomega.Equal(To.Nanos),
				//}))
				gomega.Expect(res.GetEntries()).Should(gomega.BeEmpty())
			})
		})

		ginkgo.Context("Expire handler", func() {
			ginkgo.It("should reject request without organization", func() {
				req := &grpc_unified_logging_go.ExpirationRequest{}
				res, err := coordClient.Expire(context.Background(), req)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(res).Should(gomega.BeNil())
			})
			ginkgo.It("should return ok on valid request", func() {
				res, err := coordClient.Expire(context.Background(), ValidExpirationRequest)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(res).Should(gomega.Equal(&grpc_common_go.Success{}))
			})
		})
	})

	ginkgo.Context("Slave client", func() {
		ginkgo.Context("Search handler", func() {
			ginkgo.It("should reject request without organization", func() {
				req := &grpc_unified_logging_go.SearchRequest{}
				res, err := slaveClient.Search(context.Background(), req)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(res).Should(gomega.BeNil())
			})
			ginkgo.It("should return requested values in result", func() {
				res, err := slaveClient.Search(context.Background(), ValidSearchRequest)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(res.GetOrganizationId()).Should(gomega.Equal(OrganizationId))
				gomega.Expect(res.GetAppInstanceId()).Should(gomega.Equal(AppInstanceId))
				//gomega.Expect(*res.GetFrom()).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				//	"Seconds": gomega.Equal(int64(0)),
				//	"Nanos":   gomega.Equal(int32(0)),
				//}))
				//gomega.Expect(*res.GetTo()).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				//	"Seconds": gomega.Equal(To.Seconds),
				//	"Nanos":   gomega.Equal(To.Nanos),
				//}))
				gomega.Expect(res.GetEntries()).Should(gomega.BeEmpty())
			})
		})

		ginkgo.Context("Expire handler", func() {
			ginkgo.It("should reject request without organization", func() {
				req := &grpc_unified_logging_go.ExpirationRequest{}
				res, err := slaveClient.Expire(context.Background(), req)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(res).Should(gomega.BeNil())
			})
			ginkgo.It("should return ok on valid request", func() {
				res, err := slaveClient.Expire(context.Background(), ValidExpirationRequest)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(res).Should(gomega.Equal(&grpc_common_go.Success{}))
			})
		})
	})
})
*/
