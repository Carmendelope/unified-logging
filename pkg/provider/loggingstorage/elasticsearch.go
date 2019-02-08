/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// ElasticSearch logging storage provider implementation

package loggingstorage

import (
	"github.com/nalej/derrors"
	"github.com/nalej/unified-logging/pkg/entities"
)

type ElasticSearch struct {
	address string
}

func NewElasticSearch(address string) *ElasticSearch {
	return &ElasticSearch{
		address: address,
	}
}

func (es *ElasticSearch) Search(request *entities.SearchRequest, limit int) (entities.LogEntries, derrors.Error) {
	// Log

	return nil, nil
}

func (es *ElasticSearch) Expire(request *entities.SearchRequest, retention string) derrors.Error {
	// Log

	return nil
}
