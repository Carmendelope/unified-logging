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

// Manager for unified logging coordinator

package manager

import (
	"context"
	"github.com/nalej/derrors"

	"github.com/nalej/grpc-app-cluster-api-go"
	"github.com/nalej/unified-logging/internal/pkg/client"
	"github.com/rs/zerolog/log"
)

type ExecFunc func(context.Context, grpc_app_cluster_api_go.UnifiedLoggingClient, int) (int, error)

type LoggingExecutor struct {
	clientFactory client.LoggingClientFactory
	params        *client.LoggingClientParams
}

func NewLoggingExecutor(factory client.LoggingClientFactory, params *client.LoggingClientParams) *LoggingExecutor {
	return &LoggingExecutor{
		clientFactory: factory,
		params:        params,
	}
}

func (le *LoggingExecutor) ExecRequests(ctx context.Context, hosts []ClusterInfo, f ExecFunc) (int, []string, derrors.Error) {
	// TODO: concurrent execution of request
	var total int = 0
	errorIds := make ([]string, 0)

	for i, host := range hosts {
		log.Debug().Str("host", host.host).Msg("executing on host")
		client, err := le.clientFactory(host.host, le.params)
		if err != nil {
			log.Warn().Str("host", host.host).Err(err).Msg("failed creating connection")
			errorIds = append(errorIds, host.id)
			continue
		}

		count, err := f(ctx, client, i)
		if err != nil {
			log.Warn().Str("host", host.host).Err(err).Msg("failed executing command")
			errorIds = append(errorIds, host.id)
			// Continue on to next host - after trying to close connection
		}
		total += count
		log.Debug().Int("count", count).Int("total", total).Msg("rows returned")

		cerr := client.Close()
		if cerr != nil {
			log.Warn().Str("host", host.host).Err(cerr).Msg("failed closing connection")
			// continue anyway
		}
	}

	// TODO: collect errors

	return total, errorIds, nil
}
