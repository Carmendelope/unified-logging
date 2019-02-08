/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package slave

import (
	"github.com/nalej/derrors"

	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"
	"github.com/nalej/unified-logging/internal/pkg/handler"
	"github.com/nalej/unified-logging/internal/pkg/managers"

	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/nalej/grpc-unified-logging-go"

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
		conf,
		tools.NewGenericGRPCServer(uint32(conf.Port)),
	}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() {
	// Create ElasticSearch provider
	elasticProvider := loggingstorage.NewElasticSearch(s.Configuration.ElasticAddress)
	_ = elasticProvider
	// Create managers and handler
        // searchManager := search.NewManager(elasticProvider)
        // expireManager := expire.NewManager(elasticProvider)
	var searchManager managers.Search = nil
	var expireManager managers.Expire = nil
	handler := handler.NewHandler(searchManager, expireManager)

	// Register handler
	grpc_unified_logging_go.RegisterSlaveServer(s.Server.Server, handler)

	s.Server.Run()
}
