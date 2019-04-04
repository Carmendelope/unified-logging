/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package coord

import (
	"fmt"
	"net"

	"github.com/nalej/derrors"

	"github.com/nalej/unified-logging/internal/pkg/handler"
	"github.com/nalej/unified-logging/internal/pkg/client"

	"github.com/nalej/unified-logging/internal/app/coord/manager"

	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-infrastructure-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/rs/zerolog/log"
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
		UseTLS: s.Configuration.UseTLS,
		Insecure: s.Configuration.Insecure,
		CACert: s.Configuration.CACert,
	}
	executor := manager.NewLoggingExecutor(client.NewGRPCLoggingClient, params)

	// Create managers and handler
	clientManager := manager.NewManager(appsClient, clustersClient, executor, s.Configuration.AppClusterPrefix, s.Configuration.AppClusterPort)
	handler := handler.NewHandler(clientManager, clientManager)

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
