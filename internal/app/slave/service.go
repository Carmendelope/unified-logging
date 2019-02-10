/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package slave

import (
	"fmt"
	"net"

	"github.com/nalej/derrors"

	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

	"github.com/nalej/unified-logging/internal/pkg/handler"

	"github.com/nalej/unified-logging/internal/app/slave/search"
	"github.com/nalej/unified-logging/internal/app/slave/expire"

	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/nalej/grpc-unified-logging-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/rs/zerolog/log"
)

// Service with configuration and gRPC server
type Service struct {
	Configuration *Config
	Server *tools.GenericGRPCServer
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
	// Create ElasticSearch provider
	elasticProvider := loggingstorage.NewElasticSearch(s.Configuration.ElasticAddress)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		return derrors.NewUnavailableError("failed to listen", err)
	}

	// Create managers and handler
        searchManager := search.NewManager(elasticProvider)
        expireManager := expire.NewManager(elasticProvider)
	handler := handler.NewHandler(searchManager, expireManager)

	// Create server and register handler
	server := grpc.NewServer()
	grpc_unified_logging_go.RegisterSlaveServer(server, handler)

	reflection.Register(server)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := server.Serve(lis); err != nil {
		return derrors.NewUnavailableError("failed to serve", err)
	}

	return nil
}
