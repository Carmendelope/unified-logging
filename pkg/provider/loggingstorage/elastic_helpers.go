/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// ElasticSearch helper functions

package loggingstorage

import (
	"encoding/json"
	"time"

	"github.com/nalej/derrors"
	"github.com/nalej/unified-logging/pkg/entities"

        "github.com/rs/zerolog/log"

	"github.com/olivere/elastic"
)

func getLogEntries(searchResult *elastic.SearchResult) (entities.LogEntries, derrors.Error) {
	num := searchResult.Hits.TotalHits
	log.Debug().Int64("hits", num).Int("hits_len", len(searchResult.Hits.Hits)).Msg("matching log lines found")

	result := make(entities.LogEntries, len(searchResult.Hits.Hits))

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

func createFilterQuery(filters entities.SearchFilter, union bool) (*elastic.BoolQuery) {
	// Determine if we need one or all filters to match
        query := elastic.NewBoolQuery()
        if union {
		// We need just a single match out of all filters
                query = query.MinimumShouldMatch("1")
        } else {
		// We need all filters to match
		query = query.MinimumShouldMatch("100%")
	}

	// Build filter query
        for k, values := range filters {
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

	return query
}

func createTimeQuery(from, to time.Time) elastic.Query {
	query := elastic.NewRangeQuery(entities.TimestampField)
	if !from.IsZero() {
		query = query.From(from)
	}
	if !to.IsZero() {
		query = query.To(to)
	}

	return query
}

// Debug output for query string
func queryDebug(query elastic.Query) {
	if d := log.Debug(); d.Enabled() {
		// Only get query text if we actually want debug logging
		source, err := query.Source()
		if err != nil {
			d.Err(err).Msg("error getting query string")
		} else {
			d.Interface("query", source).Msg("executing search")
		}
	}
}
