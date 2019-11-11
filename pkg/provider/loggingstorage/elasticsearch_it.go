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

// Overloaded ElasticSearch with some functionality used for integration tests

package loggingstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/nalej/derrors"

	"github.com/olivere/elastic"
)

const templateJSON = `
{
  "index_patterns": [
    "%s_*"
  ],
  "mappings": {
    "doc": {
      "properties": {
        "@timestamp": {
          "type": "date"
        },
        "message": {
          "type": "text"
        },
        "stream": {
          "type": "text"
        },
        "kubernetes": {
          "properties": {
            "namespace": {
              "type": "keyword"
            },
            "labels": {
              "properties": {
                "nalej-organization": {
                  "type": "keyword"
                },
                "nalej-app-instance-id": {
                  "type": "keyword"
                },
                "nalej-service-group-instance-id": {
                  "type": "keyword"
                }
              }
            }
          }
        }
      }
    }
  }
}
`

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

const templateName = "integration_test"

type ElasticITEntry struct {
	Timestamp  time.Time                `json:"@timestamp"`
	Stream     string                   `json:"stream"`
	Message    string                   `json:"message"`
	Kubernetes ElasticITEntryKubernetes `json:"kubernetes"`
}

type ElasticITEntryKubernetes struct {
	Namespace string                         `json:"namespace"`
	Labels    ElasticITEntryKubernetesLabels `json:"labels"`
}

type ElasticITEntryKubernetesLabels struct {
	OrganizationID         string `json:"nalej-organization"`
	AppInstanceID          string `json:"nalej-app-instance-id"`
	ServiceGroupInstanceID string `json:"nalej-service-group-instance-id"`
}

type ElasticSearchIT struct {
	*ElasticSearch

	// Prefix everything with this such that we can run test concurrently
	PrefixStr string
}

func (es *ElasticSearchIT) Prefix(s string) string {
	return fmt.Sprintf("%s_%s", es.PrefixStr, s)
}

func (es *ElasticSearchIT) InitTemplate() derrors.Error {
	client, derr := es.Connect()
	if derr != nil {
		return derr
	}

	template := fmt.Sprintf(templateJSON, es.PrefixStr)
	_, err := client.IndexPutTemplate(es.Prefix(templateName)).BodyString(template).Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed adding template", err)
	}

	return nil

}

func (es *ElasticSearchIT) Add(entries []*ElasticITEntry) derrors.Error {
	client, derr := es.Connect()
	if derr != nil {
		return derr
	}

	toAdd := client.Bulk().Index(es.Prefix(entries[0].Kubernetes.Namespace)).Type("doc")
	for _, entry := range entries {
		toAdd = toAdd.Add(elastic.NewBulkIndexRequest().Doc(entry))
	}

	_, err := toAdd.Refresh("wait_for").Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed adding entry", err)
	}

	return nil
}

func (es *ElasticSearchIT) Clear() derrors.Error {
	client, derr := es.Connect()
	if derr != nil {
		return derr
	}

	_, err := client.DeleteIndex(es.Prefix("*")).Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed deleting indices", err)
	}

	templates, err := client.IndexGetTemplate(es.Prefix("*")).Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed listing templates", err)
	}

	for t, _ := range templates {
		_, err := client.IndexDeleteTemplate(t).Do(context.Background())
		if err != nil {
			return derrors.NewInternalError("failed listing templates", err)
		}
	}

	return nil
}

func (es *ElasticSearchIT) AddTestData() derrors.Error {
	// Add some data
	for _, e := range es.generateEntries() {
		err := es.Add(e)
		if err != nil {
			return err
		}
	}

	return nil
}

// 10 lines for each org/app/sg combo, with 10 seconds between lines, starting at startTime
func (es *ElasticSearchIT) generateEntries() [][]*ElasticITEntry {
	entriesList := make([][]*ElasticITEntry, 0)

	for org, apps := range instances {
		for app, sgs := range apps {
			entries := make([]*ElasticITEntry, 0)
			for _, sg := range sgs {
				t := startTime
				for i := 0; i < 10; i++ {
					entry := &ElasticITEntry{
						Timestamp: t,
						Stream:    "stdout",
						Message:   fmt.Sprintf("Log line %s %s %s %d", org, app, sg, i),
						Kubernetes: ElasticITEntryKubernetes{
							Namespace: fmt.Sprintf("%s-%s", es.Prefix(org), es.Prefix(app)), // Hope it's not longer than 64
							Labels: ElasticITEntryKubernetesLabels{
								OrganizationID:         es.Prefix(org),
								AppInstanceID:          es.Prefix(app),
								ServiceGroupInstanceID: es.Prefix(sg),
							},
						},
					}

					entries = append(entries, entry)
					t = t.Add(time.Second * 10)
				}
			}
			entriesList = append(entriesList, entries)
		}
	}

	return entriesList
}
