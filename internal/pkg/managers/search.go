/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package managers

import (
	"context"

	"github.com/nalej/derrors"

	grpc "github.com/nalej/grpc-unified-logging-go"
)

// Interface for Search Manager
type Search interface {
	Search(context.Context, *grpc.SearchRequest) (*grpc.LogResponse, derrors.Error)
}
