/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Manager for unified logging coordinator

package manager

import (
	"context"

        "github.com/nalej/derrors"

	"github.com/nalej/unified-logging/internal/pkg/client"
	"github.com/nalej/grpc-app-cluster-api-go"
        "github.com/rs/zerolog/log"
)

type ExecFunc func(context.Context, grpc_app_cluster_api_go.UnifiedLoggingClient, int) (int, error)

type LoggingExecutor struct {
	clientFactory client.LoggingClientFactory
	params *client.LoggingClientParams
}

func NewLoggingExecutor(factory client.LoggingClientFactory, params *client.LoggingClientParams) *LoggingExecutor {
	return &LoggingExecutor{
		clientFactory: factory,
		params: params,
	}
}

func (le *LoggingExecutor) ExecRequests(ctx context.Context, hosts []string, f ExecFunc) (int, derrors.Error) {
	// TODO: concurrent execution of request
	var total int = 0

	for i, host := range(hosts) {
		log.Debug().Str("host", host).Msg("executing on host")
		client, err := le.clientFactory(host, le.params)
		if err != nil {
			log.Warn().Str("host", host).Err(err).Msg("failed creating connection")
			continue
		}

		count, err := f(ctx, client, i)
		if err != nil {
			log.Warn().Str("host", host).Err(err).Msg("failed executing command")
			// Continue on to next host - after trying to close connection
		}
		total += count
		log.Debug().Int("count", count).Int("total", total).Msg("rows returned")

		cerr := client.Close()
		if cerr != nil {
			log.Warn().Str("host", host).Err(cerr).Msg("failed closing connection")
			// continue anyway
		}
	}

	// TODO: collect errors

	return total, nil
}
