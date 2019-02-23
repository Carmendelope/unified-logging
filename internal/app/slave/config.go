/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Wrapper for the configuration properties.

package slave

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
)

// Config struct for the API service.
type Config struct {
	// Port where the API service will listen requests.
	Port int
	// Address with host:port of the ElasticSearch server
	ElasticAddress string
}

// Validate the configuration.
func (conf *Config) Validate() derrors.Error {
	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be specified")
	}
	if conf.ElasticAddress == "" {
		return derrors.NewInvalidArgumentError("elasticAddress is required")
	}
	return nil
}

// Print the current API configuration to the log.
func (conf *Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("URL", conf.ElasticAddress).Msg("ElasticSearch")
}
