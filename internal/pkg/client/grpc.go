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

// Application cluster unified logging GRPC client

package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-app-cluster-api-go"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCLoggingClient struct {
	grpc_app_cluster_api_go.UnifiedLoggingClient
	conn *grpc.ClientConn
}

func NewGRPCLoggingClient(address string, params *LoggingClientParams) (LoggingClient, error) {
	var options []grpc.DialOption
	var hostname string

	log.Debug().Str("address", address).
		Bool("tls", params.UseTLS).
		Str("cert", params.CACertPath).
		Str("cert", params.ClientCertPath).
		Bool("skipServerCertValidation", params.SkipServerCertValidation).
		Msg("creating connection")

	if params.UseTLS {
		rootCAs := x509.NewCertPool()
		splitHostname := strings.Split(address, ":")
		if len(splitHostname) != 2 {
			hostname = splitHostname[0]
		} else {
			return nil, derrors.NewInvalidArgumentError("server address incorrectly set")
		}

		tlsConfig := &tls.Config{
			ServerName: hostname,
		}

		if params.CACertPath != "" {
			log.Debug().Str("serverCertPath", params.CACertPath).Msg("loading server certificate")
			serverCert, err := ioutil.ReadFile(params.CACertPath)
			if err != nil {
				return nil, derrors.NewInternalError("Error loading server certificate")
			}
			added := rootCAs.AppendCertsFromPEM(serverCert)
			if !added {
				return nil, derrors.NewInternalError("cannot add server certificate to the pool")
			}
			tlsConfig.RootCAs = rootCAs
		}

		if params.ClientCertPath != "" {
			log.Debug().Str("clientCertPath", params.ClientCertPath).Msg("loading client certificate")
			clientCert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/tls.crt", params.ClientCertPath), fmt.Sprintf("%s/tls.key", params.ClientCertPath))
			if err != nil {
				log.Error().Str("error", err.Error()).Msg("Error loading client certificate")
				return nil, derrors.NewInternalError("Error loading client certificate")
			}

			tlsConfig.Certificates = []tls.Certificate{clientCert}
			tlsConfig.BuildNameToCertificate()
		}

		log.Debug().Str("address", hostname).Bool("useTLS", params.UseTLS).Str("serverCertPath", params.CACertPath).Bool("skipServerCertValidation", params.SkipServerCertValidation).Msg("creating secure connection")

		if params.SkipServerCertValidation {
			log.Debug().Msg("skipping server cert validation")
			tlsConfig.InsecureSkipVerify = true
		}

		creds := credentials.NewTLS(tlsConfig)
		log.Debug().Interface("creds", creds.Info()).Msg("Secure credentials")
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, options...)
	if err != nil {
		return nil, err
	}

	client := grpc_app_cluster_api_go.NewUnifiedLoggingClient(conn)

	return &GRPCLoggingClient{client, conn}, nil
}

func (c *GRPCLoggingClient) Close() error {
	return c.conn.Close()
}

func addCert(pool *x509.CertPool, cert string) error {
	caCert, err := ioutil.ReadFile(cert)
	if err != nil {
		return err
	}

	added := pool.AppendCertsFromPEM(caCert)
	if !added {
		return fmt.Errorf("Failed to add certificate from %s", cert)
	}

	return nil
}
