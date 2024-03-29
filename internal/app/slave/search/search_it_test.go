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

package search

/*
import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-utils/pkg/test"

	"github.com/nalej/unified-logging/internal/pkg/handler"
	"github.com/nalej/unified-logging/internal/pkg/utils"
	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var _ = ginkgo.Describe("Search", func() {
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

	var from, to, start, end, toEarly time.Time

	ginkgo.BeforeSuite(func() {
		// Set prefix to be able to run tests concurrently
		prefix := "search"

		// Create Elastic IT provider
		elasticProvider := loggingstorage.NewElasticSearch(elasticAddress)
		provider = &loggingstorage.ElasticSearchIT{elasticProvider, prefix}

		// Initialize template
		derr := provider.InitTemplate()
		gomega.Expect(derr).Should(gomega.Succeed())

		// Add some data
		derr = provider.AddTestData()
		gomega.Expect(derr).Should(gomega.Succeed())

		// Create listener and server
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		// Create and register manager and handler
		searchManager := NewManager(provider)
		h := handler.NewHandler(searchManager, nil)
		grpc_unified_logging_go.RegisterSlaveServer(server, h)

		// Launch test server
		test.LaunchServer(server, listener)

		// Create client
		client = grpc_unified_logging_go.NewSlaveClient(conn)

		// Time bounds
		startTime := time.Unix(1550789643, 0).UTC() // From loggingstorage.elasticsearch_it.go
		from = startTime.Add(time.Second * 30)
		start = startTime

		to = startTime.Add(time.Second * 80)
		end = startTime.Add(time.Second * 90)

		toEarly = time.Unix(946684800, 0).UTC() // 1/1/2000
	})

	ginkgo.Context("Search", func() {
		ginkgo.It("should be able to retrieve logs for an application instance", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-2")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: org,
				AppInstanceId:  app,
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"From":           gomega.Equal(start),
				"To":             gomega.Equal(end),
			}))

			msgs := make([]string, len(res.Entries))
			for i, e := range res.Entries {
				msgs[i] = e.Msg
			}

			expected := []string{}
			for i := 0; i < 10; i++ {
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-3 %d", i))
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-4 %d", i))
			}

			gomega.Expect(msgs).Should(gomega.ConsistOf(expected))
		})
		ginkgo.It("should be able to retrieve logs for a service group instance", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-1")
			sg := provider.Prefix("sg-inst-id-2")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId:         org,
				AppInstanceId:          app,
				ServiceGroupInstanceId: sg,
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"From":           gomega.Equal(start),
				"To":             gomega.Equal(end),
			}))

			msgs := make([]string, len(res.Entries))
			for i, e := range res.Entries {
				msgs[i] = e.Msg
			}

			expected := []string{}
			for i := 0; i < 10; i++ {
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-1 sg-inst-id-2 %d", i))
			}

			gomega.Expect(msgs).Should(gomega.ConsistOf(expected))
		})
		ginkgo.It("should return an empty result when searching for non-existing application instance", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-foo")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: org,
				AppInstanceId:  app,
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"From":           gomega.BeNil(),
				"To":             gomega.BeNil(),
			}))

			gomega.Expect(res.Entries).Should(gomega.BeEmpty())
		})
		ginkgo.It("should be able to retrieve logs from a certain point in time", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-2")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: org,
				AppInstanceId:  app,
				From:           from.Unix(),
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"From":           gomega.Equal(from),
				"To":             gomega.Equal(end),
			}))

			msgs := make([]string, len(res.Entries))
			for i, e := range res.Entries {
				msgs[i] = e.Msg
			}

			expected := []string{}
			for i := 3; i < 10; i++ {
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-3 %d", i))
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-4 %d", i))
			}

			gomega.Expect(msgs).Should(gomega.ConsistOf(expected))
		})
		ginkgo.It("should be able to retrieve logs to a certain point in time", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-2")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: org,
				AppInstanceId:  app,
				To:             to.Unix(),
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"From":           gomega.Equal(start),
				"To":             gomega.Equal(to),
			}))

			msgs := make([]string, len(res.Entries))
			for i, e := range res.Entries {
				msgs[i] = e.Msg
			}

			expected := []string{}
			for i := 0; i <= 8; i++ {
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-3 %d", i))
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-4 %d", i))
			}

			gomega.Expect(msgs).Should(gomega.ConsistOf(expected))
		})
		ginkgo.It("should be able to retrieve logs between two points in time", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-2")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: org,
				AppInstanceId:  app,
				From:           from.Unix(),
				To:             to.Unix(),
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"From":           gomega.Equal(from),
				"To":             gomega.Equal(to),
			}))

			msgs := make([]string, len(res.Entries))
			for i, e := range res.Entries {
				msgs[i] = e.Msg
			}

			expected := []string{}
			for i := 3; i <= 8; i++ {
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-3 %d", i))
				expected = append(expected, fmt.Sprintf("Log line org-id-1 app-inst-id-2 sg-inst-id-4 %d", i))
			}

			gomega.Expect(msgs).Should(gomega.ConsistOf(expected))

		})
		ginkgo.It("should return an empty result for points in time with no log entries", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-2")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: org,
				AppInstanceId:  app,
				To:             toEarly.Unix(),
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"From":           gomega.BeNil(),
				"To":             gomega.BeNil(),
				"Entries":        gomega.HaveLen(0),
			}))

		})
		ginkgo.It("should be able to retrieve logs matching a certain message", func() {
			org := provider.Prefix("org-id-1")
			app := provider.Prefix("app-inst-id-1")

			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: org,
				AppInstanceId:  app,
				MsgQueryFilter: " 5",
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal(org),
				"AppInstanceId":  gomega.Equal(app),
				"Entries":        gomega.HaveLen(2),
			}))

			gomega.Expect(res.Entries[0].Msg).Should(gomega.ContainSubstring(" 5"))
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
*/
