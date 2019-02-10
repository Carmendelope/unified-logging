/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// ElasticSearch logging storage provider implementation

package loggingstorage

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/nalej/derrors"
	"github.com/nalej/unified-logging/pkg/entities"

        "github.com/rs/zerolog/log"

	"github.com/olivere/elastic"
)

type ElasticSearch struct {
	address string
}

func NewElasticSearch(address string) *ElasticSearch {
	return &ElasticSearch{
		address: address,
	}
}

func (es *ElasticSearch) Search(ctx context.Context, request *entities.SearchRequest, limit int) (entities.LogEntries, derrors.Error) {
	log.Debug().Str("address", es.address).Msg("elastic search")

        client, err := elastic.NewSimpleClient(
                elastic.SetURL(fmt.Sprintf("http://%s", es.address)))
        if err != nil {
                return nil, derrors.NewUnavailableError("elastic search connection has failed", err)
        }

	// Determine if we need one or all filters to match
        query := elastic.NewBoolQuery()
        if request.IsUnionFilter {
		// We need just a single match out of all filters
                query = query.MinimumShouldMatch("1")
        } else {
		// We need all filters to match
		query = query.MinimumShouldMatch("100%")
	}

	// Build filter query
        for k, values := range request.Filters {
		if len(values) == 0 {
			continue
		}

		subQuery := elastic.NewBoolQuery()
		for _, v := range values {
			subQuery = subQuery.Should(elastic.NewTermQuery(k.String(), v))
		}
		subQuery = subQuery.MinimumNumberShouldMatch(1)

		query = query.Should(subQuery)
        }

	// Add required filter for actual log line
        if request.MsgFilter != "" {
                query = query.Must(elastic.NewQueryStringQuery(request.MsgFilter).
                        DefaultField(entities.MessageField).AllowLeadingWildcard(true))
        }

	// Add time constraints
        if !request.From.IsZero() || !request.To.IsZero() {
                subQuery := elastic.NewRangeQuery(entities.TimestampField)
                if !request.From.IsZero() {
                        subQuery = subQuery.From(request.From)
                }
                if !request.To.IsZero() {
                        subQuery = subQuery.To(request.To)
                }
                query = query.Must(subQuery)
        }

	// Execute
        searchResult, err := client.Search().Index("_all").Query(query).
                Sort(entities.TimestampField, false).
                Size(limit).
                Do(ctx)

        if err != nil {
                return nil, derrors.NewInternalError("elastic search query has failed", err)
        }

	// Create result
	_ = searchResult
	return es.extractResult(searchResult)
}

func (es *ElasticSearch) extractResult(
        searchResult *elastic.SearchResult) (entities.LogEntries, derrors.Error) {
	num := searchResult.Hits.TotalHits
	log.Debug().Int64("hits", num).Msg("matching log lines found")

	result := make(entities.LogEntries, num)

	for k, hit := range searchResult.Hits.Hits {
		var entry entities.LogEntry
		err := json.Unmarshal(*hit.Source, &entry)
		if err != nil {
			return nil, derrors.NewInternalError("elastic document deserialization error", err)
		}
		result[k] = &entry
	}

	return result, nil
}

func (es *ElasticSearch) Expire(ctx context.Context, request *entities.SearchRequest, retention string) derrors.Error {
	// Log

	return nil
}
