/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Application cluster unified logging client interface

package client

import (
	"github.com/nalej/grpc-app-cluster-api-go"
)

type LoggingClient interface {
	grpc_app_cluster_api_go.UnifiedLoggingClient
	Close() error
}

type LoggingClientFactory func(address string) (LoggingClient, error)
