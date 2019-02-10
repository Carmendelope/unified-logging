/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Logging storage provider interface

package loggingstorage

import (
	"context"

	"github.com/nalej/derrors"
	"github.com/nalej/unified-logging/pkg/entities"
)

// Provider is the interface of the Logging provider.
type Provider interface {
	Search(ctx context.Context, request *entities.SearchRequest, limit int) (entities.LogEntries, derrors.Error)
	Expire(ctx context.Context, request *entities.SearchRequest, retention string) derrors.Error
}
