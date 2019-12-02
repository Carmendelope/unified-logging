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

// ElasticSearch helper functions

package loggingstorage

import (
	"encoding/json"
	"github.com/nalej/derrors"
	"github.com/nalej/unified-logging/pkg/entities"
	"time"

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

func createFilterQuery(filters entities.SearchFilter) *elastic.BoolQuery {
	// Determine if we need one or all filters to match
	query := elastic.NewBoolQuery()

	// Build filter query
	for k, values := range filters {
		if len(values) == 0 {
			continue
		}
		for _, v := range values {
			query = query.Must(elastic.NewTermQuery(k.String(), v))
		}
	}

	return query
}

func createTimeQuery(from, to int64) elastic.Query {
	query := elastic.NewRangeQuery(entities.TimestampField.String())
	if from != 0  {
		query = query.From(time.Unix(from, 0))
	}
	if to != 0 {
		query = query.To(time.Unix(to,0))
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
