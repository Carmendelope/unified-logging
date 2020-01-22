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

package coord

import (
	"fmt"
	"net"

	"github.com/nalej/derrors"

	"github.com/nalej/unified-logging/internal/pkg/client"
	"github.com/nalej/unified-logging/internal/pkg/handler"

	"github.com/nalej/unified-logging/internal/app/coord/manager"

	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-unified-logging-go"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Service with configuration and gRPC server
type Service struct {
	Configuration *Config
}

func NewService(conf *Config) (*Service, derrors.Error) {
	err := conf.Validate()
	if err != nil {
		log.Error().Msg("Invalid configuration")
		return nil, err
	}
	conf.Print()

	return &Service{
		Configuration: conf,
	}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() derrors.Error {
	// Create system model connection
	smConn, err := grpc.Dial(s.Configuration.SystemModelAddress, grpc.WithInsecure())
	if err != nil {
		return derrors.NewUnavailableError("cannot create connection with the system model", err)
	}

	// Create clients
	appsClient := grpc_application_go.NewApplicationsClient(smConn)
	clustersClient := grpc_infrastructure_go.NewClustersClient(smConn)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		return derrors.NewUnavailableError("failed to listen", err)
	}

	// Executor for application cluster requests
	params := &client.LoggingClientParams{
		UseTLS:                   s.Configuration.UseTLS,
		SkipServerCertValidation: s.Configuration.SkipServerCertValidation,
		CACertPath:               s.Configuration.CACertPath,
		ClientCertPath:           s.Configuration.ClientCertPath,
	}
	executor := manager.NewLoggingExecutor(client.NewGRPCLoggingClient, params)

	// Create managers and handler
	clientManager := manager.NewManager(appsClient, clustersClient, executor, s.Configuration.AppClusterPrefix, s.Configuration.AppClusterPort)
	handler := handler.NewHandler(clientManager, clientManager)

	go clientManager.Clean()

	// Create server and register handler
	server := grpc.NewServer()
	grpc_unified_logging_go.RegisterCoordinatorServer(server, handler)

	reflection.Register(server)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := server.Serve(lis); err != nil {
		return derrors.NewUnavailableError("failed to serve", err)
	}

	return nil
}
