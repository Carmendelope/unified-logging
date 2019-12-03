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

package loggingstorage

const (
	OrganizationId         = "77b5425b-4276-45b8-85f4-c01f74bbc376"
	AppInstanceId          = "e5a51a0b-63ea-4736-8c1c-be3d423f28f0"
	AppInstanceId2         = "e9e38334-1da1-4f51-8f18-2bd8e2470123"
	ServiceGroupInstanceId = "413654be-3c62-48cd-beb5-86d09462a1dc"
)

/*
var jsonResult = []byte(`
{"took":4,"timed_out":false,"_shards":{"total":10,"successful":10,"skipped":0,"failed":0},"hits":{"total":15987,"max_score":1.0,"hits":[{"_index":"filebeat-6.6.0-2019.02.20","_type":"doc","_id":"fe8bDGkBNkPZggIxE_ej","_score":1.0,"_source":{"@timestamp":"2019-02-20T16:06:25.276Z","host":{"name":"filebeat-ltk7p"},"message":"Logline 1","stream":"stderr","kubernetes":{"namespace":"77b5425b-4276-45b8-85f4-c01f74bbc376-e5a51a0b-63ea-4736-8c1c-be","labels":{"nalej-service-group-id":"f15b707c-280f-4670-b502-903e59b6dcdd","nalej-stage-id":"f2a0d800-a65f-4067-93be-23ddbc2814df","nalej-organization":"77b5425b-4276-45b8-85f4-c01f74bbc376","app":"simple-mysql","nalej-service-id":"38d0328c-4d2b-4915-8aba-7bad6b374639","nalej-app-descriptor":"5500140c-72ff-41e3-8384-2ec7f7a6a601","nalej-service-instance-id":"f84a48f5-9d53-4fb1-916c-b3010a3496db","pod-template-hash":"d7674595c","component":"simple-app","nalej-service-group-instance-id":"413654be-3c62-48cd-beb5-86d09462a1dc","nalej-app-instance-id":"e5a51a0b-63ea-4736-8c1c-be3d423f28f0"}},"beat":{"hostname":"filebeat-ltk7p","version":"6.6.0","name":"filebeat-ltk7p"}}},{"_index":"filebeat-6.6.0-2019.02.20","_type":"doc","_id":"ge8bDGkBNkPZggIxE_ej","_score":1.0,"_source":{"@timestamp":"2019-02-20T16:06:25.278Z","beat":{"hostname":"filebeat-ltk7p","version":"6.6.0","name":"filebeat-ltk7p"},"host":{"name":"filebeat-ltk7p"},"message":"Logline 2","stream":"stderr","kubernetes":{"namespace":"77b5425b-4276-45b8-85f4-c01f74bbc376-e5a51a0b-63ea-4736-8c1c-be","labels":{"nalej-service-group-instance-id":"413654be-3c62-48cd-beb5-86d09462a1dc","nalej-organization":"77b5425b-4276-45b8-85f4-c01f74bbc376","component":"simple-app","nalej-app-descriptor":"5500140c-72ff-41e3-8384-2ec7f7a6a601","app":"simple-mysql","nalej-service-id":"38d0328c-4d2b-4915-8aba-7bad6b374639","nalej-service-group-id":"f15b707c-280f-4670-b502-903e59b6dcdd","nalej-stage-id":"f2a0d800-a65f-4067-93be-23ddbc2814df","nalej-service-instance-id":"f84a48f5-9d53-4fb1-916c-b3010a3496db","pod-template-hash":"d7674595c","nalej-app-instance-id":"e5a51a0b-63ea-4736-8c1c-be3d423f28f0"}}}}]}}
`)

var logEntries = entities.LogEntries{
	&entities.LogEntry{
		Timestamp: time.Unix(1550678785, 276000000).UTC(),
		Msg:       "Logline 1",
	},
	&entities.LogEntry{
		Timestamp: time.Unix(1550678785, 278000000).UTC(),
		Msg:       "Logline 2",
	},
}

var _ = ginkgo.Describe("elastic_helpers", func() {

	var emptyResult, validResult, badResult *elastic.SearchResult

	var filter, multifilter entities.SearchFilter

	var from, to, zeroTime time.Time

	ginkgo.BeforeSuite(func() {
		emptyResult = &elastic.SearchResult{
			Hits: &elastic.SearchHits{
				Hits: []*elastic.SearchHit{},
			},
		}

		validResult = &elastic.SearchResult{}
		err := json.Unmarshal(jsonResult, validResult)
		gomega.Expect(err).Should(gomega.Succeed())

		badResult = &elastic.SearchResult{}
		err = json.Unmarshal(jsonResult, badResult)
		gomega.Expect(err).Should(gomega.Succeed())
		emptyJson := json.RawMessage(``) // Malformed JSON on purpose
		badResult.Hits.Hits[1].Source = &emptyJson

		filter = entities.SearchFilter{
			entities.NamespaceField:              []string{fmt.Sprintf("%s-%s", OrganizationId, AppInstanceId)[:63]},
			entities.OrganizationIdField:         []string{OrganizationId},
			entities.AppInstanceIdField:          []string{AppInstanceId},
			entities.ServiceGroupInstanceIdField: []string{ServiceGroupInstanceId},
		}
		multifilter = entities.SearchFilter{
			entities.NamespaceField:      []string{fmt.Sprintf("%s-%s", OrganizationId, AppInstanceId)[:63]},
			entities.OrganizationIdField: []string{OrganizationId},
			entities.AppInstanceIdField:  []string{AppInstanceId, AppInstanceId2},
		}

		from = time.Unix(1550678785, 278000000).UTC()
		to = time.Unix(1550698785, 278000000).UTC()
	})

	ginkgo.Context("getLogEntries", func() {
		ginkgo.It("should convert elastic results into LogEntries", func() {
			res, err := getLogEntries(validResult)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(res).Should(gomega.HaveLen(2))
			gomega.Expect(res).Should(gomega.BeEquivalentTo(logEntries))
		})
		ginkgo.It("should error on malformed resulst", func() {
			_, err := getLogEntries(badResult)
			gomega.Expect(err).Should(gomega.HaveOccurred())
		})
		ginkgo.It("should handle empty result set", func() {
			res, err := getLogEntries(emptyResult)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(res).Should(gomega.BeEmpty())
		})
	})

	ginkgo.Context("createFilterQuery", func() {
		ginkgo.It("should create a union query", func() {
			q, err := createFilterQuery(filter).MinimumShouldMatch("1").Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonQ, err := json.Marshal(q)
			gomega.Expect(err).Should(gomega.Succeed())

			expected, err := elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.NamespaceField.String(), fmt.Sprintf("%s-%s", OrganizationId, AppInstanceId)[:63]),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.OrganizationIdField.String(), OrganizationId),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.AppInstanceIdField.String(), AppInstanceId),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.ServiceGroupInstanceIdField.String(), ServiceGroupInstanceId),
				),
			).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonE, err := json.Marshal(expected)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(jsonQ).Should(MatchUnorderedJSON(jsonE))
		})
		ginkgo.It("should create an intersection query", func() {
			q, err := createFilterQuery(filter).MinimumShouldMatch("100%").Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonQ, err := json.Marshal(q)
			gomega.Expect(err).Should(gomega.Succeed())

			expected, err := elastic.NewBoolQuery().MinimumShouldMatch("100%").Should(
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.NamespaceField.String(), fmt.Sprintf("%s-%s", OrganizationId, AppInstanceId)[:63]),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.OrganizationIdField.String(), OrganizationId),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.AppInstanceIdField.String(), AppInstanceId),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.ServiceGroupInstanceIdField.String(), ServiceGroupInstanceId),
				),
			).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonE, err := json.Marshal(expected)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(jsonQ).Should(MatchUnorderedJSON(jsonE))
		})
		ginkgo.It("should create a query with multiple filter values", func() {
			q, err := createFilterQuery(multifilter,).MinimumShouldMatch("1").Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonQ, err := json.Marshal(q)
			gomega.Expect(err).Should(gomega.Succeed())

			expected, err := elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.NamespaceField.String(), fmt.Sprintf("%s-%s", OrganizationId, AppInstanceId)[:63]),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.OrganizationIdField.String(), OrganizationId),
				),
				elastic.NewBoolQuery().MinimumNumberShouldMatch(1).Should(
					elastic.NewTermQuery(entities.AppInstanceIdField.String(), AppInstanceId),
					elastic.NewTermQuery(entities.AppInstanceIdField.String(), AppInstanceId2),
				),
			).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonE, err := json.Marshal(expected)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(jsonQ).Should(MatchUnorderedJSON(jsonE))
		})
	})

	ginkgo.Context("createTimeQuery", func() {
		ginkgo.It("should create a query with from time", func() {
			q, err := createTimeQuery(from, zeroTime).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonQ, err := json.Marshal(q)
			gomega.Expect(err).Should(gomega.Succeed())

			expected, err := elastic.NewRangeQuery(entities.TimestampField.String()).From(from).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonE, err := json.Marshal(expected)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(jsonQ).Should(MatchUnorderedJSON(jsonE))

		})
		ginkgo.It("should create a query with to time", func() {
			q, err := createTimeQuery(zeroTime, to).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonQ, err := json.Marshal(q)
			gomega.Expect(err).Should(gomega.Succeed())

			expected, err := elastic.NewRangeQuery(entities.TimestampField.String()).To(to).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonE, err := json.Marshal(expected)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(jsonQ).Should(MatchUnorderedJSON(jsonE))
		})
		ginkgo.It("should create a query with from and to time", func() {
			q, err := createTimeQuery(from, to).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonQ, err := json.Marshal(q)
			gomega.Expect(err).Should(gomega.Succeed())

			expected, err := elastic.NewRangeQuery(entities.TimestampField.String()).From(from).To(to).Source()
			gomega.Expect(err).Should(gomega.Succeed())
			jsonE, err := json.Marshal(expected)
			gomega.Expect(err).Should(gomega.Succeed())

			gomega.Expect(jsonQ).Should(MatchUnorderedJSON(jsonE))
		})
	})
})
*/
