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

// Wrapper for the configuration properties.

package coord

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
)

// Config struct for the API service.
type Config struct {
	// Port where the API service will listen requests.
	Port int
	// Address with host:port of the ElasticSearch server
	SystemModelAddress string
	// Prefix for application cluster hostnames
	AppClusterPrefix string
	// Port used by app-cluster-api
	AppClusterPort int
	// Use TLS
	UseTLS bool
	// Don't validate TLS certificates
	SkipServerCertValidation bool
	// Alternative certificate path to use for validation
	CACertPath string
	// client certificate path to use for validation
	ClientCertPath string
}

// Validate the configuration.
func (conf *Config) Validate() derrors.Error {
	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be specified")
	}
	if conf.SystemModelAddress == "" {
		return derrors.NewInvalidArgumentError("systemModelAddress is required")
	}
	if conf.AppClusterPort <= 0 {
		return derrors.NewInvalidArgumentError("appClusterPort is required")
	}
	if conf.CACertPath == "" {
		return derrors.NewInvalidArgumentError("caCertPath is required")
	}
	if conf.ClientCertPath == "" {
		return derrors.NewInvalidArgumentError("clientCertPath is required")
	}
	return nil
}

// Print the current API configuration to the log.
func (conf *Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("URL", conf.SystemModelAddress).Msg("systemModelAddress")
	log.Info().Str("prefix", conf.AppClusterPrefix).Msg("appClusterPrefix")
	log.Info().Int("port", conf.AppClusterPort).Msg("appClusterPort")
	log.Info().Bool("tls", conf.UseTLS).Bool("skipServerCertValidation", conf.SkipServerCertValidation).Str("cert", conf.CACertPath).Str("cert", conf.ClientCertPath).Msg("TLS parameters")
}
