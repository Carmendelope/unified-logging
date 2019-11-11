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

/*
RUN_INTEGRATION_TEST=true
IT_ELASTIC_ADDRESS=localhost:9200
*/

package expire

import (
	"context"
	"os"

	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-utils/pkg/test"

	"github.com/nalej/unified-logging/internal/pkg/handler"
	"github.com/nalej/unified-logging/internal/pkg/utils"
	"github.com/nalej/unified-logging/pkg/entities"
	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var _ = ginkgo.Describe("Expire", func() {
	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		elasticAddress = os.Getenv("IT_ELASTIC_ADDRESS")
	)

	if elasticAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	var provider *loggingstorage.ElasticSearchIT
	var listener *bufconn.Listener
	var server *grpc.Server

	var client grpc_unified_logging_go.SlaveClient

	ginkgo.BeforeSuite(func() {
		// Set prefix to be able to run tests concurrently
		prefix := "expire"

		// Create Elastic IT provider
		elasticProvider := loggingstorage.NewElasticSearch(elasticAddress)
		provider = &loggingstorage.ElasticSearchIT{elasticProvider, prefix}

		// Initialize template
		derr := provider.InitTemplate()
		gomega.Expect(derr).Should(gomega.Succeed())

		// Add some data
		derr = provider.AddTestData()
		gomega.Expect(derr).Should(gomega.Succeed())

		// Also check we actually have data to be sure we expire data
		filters := &entities.FilterFields{
			OrganizationId: provider.Prefix("org-id-1"),
		}
		sreq := &entities.SearchRequest{
			Filters: filters.ToFilters(),
		}
		gomega.Expect(provider.Search(context.Background(), sreq, -1)).Should(gomega.HaveLen(40))

		// Create listener and server
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		// Create and register manager and handler
		expireManager := NewManager(provider)
		h := handler.NewHandler(nil, expireManager)
		grpc_unified_logging_go.RegisterSlaveServer(server, h)

		// Launch test server
		test.LaunchServer(server, listener)

		// Create client
		client = grpc_unified_logging_go.NewSlaveClient(conn)
	})

	ginkgo.Context("Expire", func() {
		ginkgo.It("should not remove anything when expiration of non-existent application instance is requested", func() {
			req := &grpc_unified_logging_go.ExpirationRequest{
				OrganizationId: provider.Prefix("org-id-1"),
				AppInstanceId:  provider.Prefix("app-inst-id-foo"),
			}
			_, err := client.Expire(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			// Check we still have all data
			filters := &entities.FilterFields{
				OrganizationId: provider.Prefix("org-id-1"),
			}
			sreq := &entities.SearchRequest{
				Filters: filters.ToFilters(),
			}
			gomega.Expect(provider.Search(context.Background(), sreq, -1)).Should(gomega.HaveLen(40))
		})
		ginkgo.It("should be able to remove all entries for an application instance", func() {
			req := &grpc_unified_logging_go.ExpirationRequest{
				OrganizationId: provider.Prefix("org-id-1"),
				AppInstanceId:  provider.Prefix("app-inst-id-1"),
			}

			_, err := client.Expire(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			// Check we have expired data
			filters := &entities.FilterFields{
				OrganizationId: req.OrganizationId,
				AppInstanceId:  req.AppInstanceId,
			}
			sreq := &entities.SearchRequest{
				Filters: filters.ToFilters(),
			}
			gomega.Expect(provider.Search(context.Background(), sreq, -1)).Should(gomega.HaveLen(0))

			// Check we have the other data still
			filters = &entities.FilterFields{
				OrganizationId: provider.Prefix("org-id-1"),
			}
			sreq = &entities.SearchRequest{
				Filters: filters.ToFilters(),
			}
			gomega.Expect(provider.Search(context.Background(), sreq, -1)).Should(gomega.HaveLen(20))
		})
	})

	ginkgo.AfterSuite(func() {
		// Clear out elastic for next test
		err := provider.Clear()
		gomega.Expect(err).Should(gomega.Succeed())

		// Stop serverr and close connections
		server.Stop()
		listener.Close()
	})
})
