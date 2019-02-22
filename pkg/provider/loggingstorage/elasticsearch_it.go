/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Overloaded ElasticSearch with some functionality used for integration tests

package loggingstorage

import (
	"context"
	"time"

        "github.com/nalej/derrors"
)

const templateJSON = `
{
  "index_patterns": [
    "*"
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

const templateName = "integration_test"

type ElasticITEntry struct {
	Timestamp time.Time `json:"@timestamp"`
	Stream string `json:"stream"`
	Message string `json:"message"`
	Kubernetes ElasticITEntryKubernetes `json:"kubernetes"`
}

type ElasticITEntryKubernetes struct {
	Namespace string `json:"namespace"`
	Labels ElasticITEntryKubernetesLabels `json:"labels"`
}

type ElasticITEntryKubernetesLabels struct {
	OrganizationID string `json:"nalej-organization"`
	AppInstanceID string `json:"nalej-app-instance-id"`
	ServiceGroupInstanceID string `json:"nalej-service-group-instance-id"`
}

type ElasticSearchIT struct {
	*ElasticSearch
}

func (es *ElasticSearchIT) InitTemplate() derrors.Error {
	client, derr := es.Connect()
	if derr != nil {
		return derr
	}

	_, err := client.IndexPutTemplate(templateName).BodyString(templateJSON).Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed adding template", err)
	}

	return nil

}

func (es *ElasticSearchIT) Add(entry *ElasticITEntry) derrors.Error {
	client, derr := es.Connect()
	if derr != nil {
		return derr
	}

	_, err := client.Index().Index(entry.Kubernetes.Namespace).Type("doc").
		BodyJson(entry).Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed adding entry", err)
	}

	return nil
}

func (es *ElasticSearchIT) Flush() derrors.Error {
	client, derr := es.Connect()
	if derr != nil {
		return derr
	}

	_, err := client.Flush().Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed flushing", err)
	}

	return nil
}

func (es *ElasticSearchIT) Clear() derrors.Error {
	client, derr := es.Connect()
	if derr != nil {
		return derr
	}

	_, err := client.DeleteIndex("_all").Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed deleting indices", err)
	}

	templates, err := client.IndexGetTemplate().Do(context.Background())
	if err != nil {
		return derrors.NewInternalError("failed listing templates", err)
	}

	for t, _ := range(templates) {
		_, err := client.IndexDeleteTemplate(t).Do(context.Background())
		if err != nil {
			return derrors.NewInternalError("failed listing templates", err)
		}
	}

	return nil
}
