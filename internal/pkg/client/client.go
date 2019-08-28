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

type LoggingClientParams struct {
	UseTLS bool
	SkipServerCertValidation bool
	CACertPath string
	ClientCertPath string
}

type LoggingClientFactory func(address string, params *LoggingClientParams) (LoggingClient, error)
