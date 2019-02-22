/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_ELASTIC_ADDRESS=localhost:9200
*/

package search

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

// OrganizationID -> ApplicationInstanceID -> ServiceGroupInstanceId
var instances = map[string]map[string][]string{
	"org-id-1": map[string][]string{
		"app-inst-id-1": []string{
			"sg-inst-id-1",
			"sg-inst-id-2",
		},
		"app-inst-id-2": []string{
			"sg-inst-id-3",
			"sg-inst-id-4",
		},
	},
}

var startTime = time.Unix(1550789643, 0).UTC()

// 10 lines for each org/app/sg combo, with 10 seconds between lines, starting at startTime
func generateEntries() []*loggingstorage.ElasticITEntry {
	entries := make([]*loggingstorage.ElasticITEntry, 0)

	currentLine := 0

	for org, apps := range(instances) {
		for app, sgs := range(apps) {
			for _, sg := range(sgs) {
				t := startTime
				for i := 0; i < 10; i++ {
					entry := &loggingstorage.ElasticITEntry{
						Timestamp: t,
						Stream: "stdout",
						Message: fmt.Sprintf("Log line %d", currentLine),
						Kubernetes: loggingstorage.ElasticITEntryKubernetes{
							Namespace: fmt.Sprintf("%s-%s", org, app), // Hope it's not longer than 64
							Labels: loggingstorage.ElasticITEntryKubernetesLabels{
								OrganizationID: org,
								AppInstanceID: app,
								ServiceGroupInstanceID: sg,
							},
						},
					}

					entries = append(entries, entry)
					t = t.Add(time.Second * 10)
					currentLine++
				}
			}
		}
	}

	return entries
}


var _ = ginkgo.Describe("Search", func() {
	defer ginkgo.GinkgoRecover()

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
		// Create Elastic IT provider
		elasticProvider := loggingstorage.NewElasticSearch(elasticAddress)
		provider = &loggingstorage.ElasticSearchIT{elasticProvider}

		// Clear out existing indices
		derr := provider.Clear()
		gomega.Expect(derr).Should(gomega.Succeed())

		// Initialize template
		derr = provider.InitTemplate()
		gomega.Expect(derr).Should(gomega.Succeed())

		// Add some data
		for _, e := range(generateEntries()) {
			err := provider.Add(e)
			gomega.Expect(err).Should(gomega.Succeed())
		}
		derr = provider.Flush()
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
		_ = client
	})

	ginkgo.Context("Search", func() {
		ginkgo.It("should be able to retrieve logs for an application instance", func() {
			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: "org-id-1",
				AppInstanceId: "app-inst-id-2",
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal("org-id-1"),
				"AppInstanceId": gomega.Equal("app-inst-id-2"),
				"From": gomega.BeNil(),
				"To": gomega.BeNil(),
			}))

			msgs := make([]string, len(res.Entries))
			for i, e := range(res.Entries) {
				msgs[i] = e.Msg
			}

			expected := []string{}
			for i := 20; i <= 39; i++ {
				expected = append(expected, fmt.Sprintf("Log line %d", i))
			}

			gomega.Expect(msgs).Should(gomega.ConsistOf(expected))
		})
		ginkgo.It("should be able to retrieve logs for a service group instance", func() {
			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: "org-id-1",
				AppInstanceId: "app-inst-id-1",
				ServiceGroupInstanceId: "sg-inst-id-2",
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal("org-id-1"),
				"AppInstanceId": gomega.Equal("app-inst-id-1"),
				"From": gomega.BeNil(),
				"To": gomega.BeNil(),
			}))

			msgs := make([]string, len(res.Entries))
			for i, e := range(res.Entries) {
				msgs[i] = e.Msg
			}

			expected := []string{}
			for i := 10; i <= 19; i++ {
				expected = append(expected, fmt.Sprintf("Log line %d", i))
			}

			gomega.Expect(msgs).Should(gomega.ConsistOf(expected))

		})
		ginkgo.It("should return an empty result when searching for non-existing application instance", func() {
			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: "org-id-1",
				AppInstanceId: "app-inst-id-foo",
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal("org-id-1"),
				"AppInstanceId": gomega.Equal("app-inst-id-foo"),
				"From": gomega.BeNil(),
				"To": gomega.BeNil(),
			}))

			gomega.Expect(res.Entries).Should(gomega.BeEmpty())

		})
		ginkgo.It("should be able to retrieve logs from a certain point in time", func() {

		})
		ginkgo.It("should be able to retrieve logs to a certain point in time", func() {

		})
		ginkgo.It("should be able to retrieve logs between two points in time", func() {

		})
		ginkgo.It("should return an empty results for points in time with no log entries", func() {

		})
		ginkgo.It("should be able to retrieve logs matching a certain message", func() {
			req := &grpc_unified_logging_go.SearchRequest{
				OrganizationId: "org-id-1",
				AppInstanceId: "app-inst-id-1",
				MsgQueryFilter: "15",
			}
			res, err := client.Search(context.Background(), req)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(*res).Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"OrganizationId": gomega.Equal("org-id-1"),
				"AppInstanceId": gomega.Equal("app-inst-id-1"),
				"From": gomega.BeNil(),
				"To": gomega.BeNil(),
			}))

			gomega.Expect(res.Entries).Should(gomega.HaveLen(1))
			gomega.Expect(res.Entries[0].Msg).Should(gomega.ContainSubstring("15"))
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
