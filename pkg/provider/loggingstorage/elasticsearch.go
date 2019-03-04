/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// ElasticSearch logging storage provider implementation

package loggingstorage

import (
	"context"
	"fmt"

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

func (es *ElasticSearch) Connect() (*elastic.Client, derrors.Error) {
	// TODO: Create long-lived client if needed
	client, err := elastic.NewSimpleClient(elastic.SetURL(fmt.Sprintf("http://%s", es.address)))
	if err != nil {
                return nil, derrors.NewUnavailableError("elastic search connection has failed", err)
	}
	return client, nil
}

func (es *ElasticSearch) Search(ctx context.Context, request *entities.SearchRequest, limit int) (entities.LogEntries, derrors.Error) {
	log.Debug().Str("address", es.address).Msg("elastic search")

        client, derr := es.Connect()
        if derr != nil {
                return nil, derr
        }

	query := createFilterQuery(request.Filters, request.IsUnionFilter)

	// Add required filter for actual log line
        if request.MsgFilter != "" {
                query = query.Must(elastic.NewQueryStringQuery(request.MsgFilter).
                        DefaultField(entities.MessageField.String()).AllowLeadingWildcard(true))
        }

	// Add time constraints
        if !request.From.IsZero() || !request.To.IsZero() {
                query = query.Must(createTimeQuery(request.From, request.To))
        }

	// Output query string for debugging
	queryDebug(query)

	// If no limit, we set to the default maximum window
	// TODO: use scroll API and pagination to retrieve results
	if limit < 0 {
		limit = 10000
	}

	// Execute
        searchResult, err := client.Search().Query(query).
		Sort(entities.TimestampField.String(), request.Order.ToAscending()).
		Size(limit).
		Do(ctx)
        if err != nil {
                return nil, derrors.NewInternalError("elastic search query has failed", err)
        }

	// Create result
	return getLogEntries(searchResult)
}

func (es *ElasticSearch) Expire(ctx context.Context, request *entities.SearchRequest) derrors.Error {
	log.Debug().Str("address", es.address).Msg("elastic expire")

        client, derr := es.Connect()
        if derr != nil {
                return derr
        }

	query := createFilterQuery(request.Filters, request.IsUnionFilter)

	// TODO: Delete a specific time range

	// Output query string for debugging
	queryDebug(query)

	// Execute
        res, err := client.DeleteByQuery().
                Query(query).Index("_all").
                Do(ctx)
        if err != nil {
                return derrors.NewInternalError("elastic expire query failed", err)
        }
        log.Debug().Int64("deleted", res.Deleted).Msg("expired entries")

	// Flush deleted docs
        _, err = elastic.NewIndicesFlushService(client).Do(ctx)
        if err != nil {
                return derrors.NewInternalError("elastic flush query failed", err)
        }

	return nil
}