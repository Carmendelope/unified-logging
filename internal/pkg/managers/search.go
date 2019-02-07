/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package managers

import (
        grpc "github.com/nalej/grpc-unified-logging-go"
)

// Interface for Search Manager
type Search interface {
	Search(*grpc.SearchRequest) (*grpc.LogResponse, error)
}
