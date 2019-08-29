/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
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

	log.Debug().Str("address", address).
		Bool("tls", params.UseTLS).
		Str("cert", params.CACertPath).
		Str("cert", params.ClientCertPath).
		Bool("skipServerCertValidation", params.SkipServerCertValidation).
		Msg("creating connection")

	if params.UseTLS {
		rootCAs := x509.NewCertPool()
		hostname := strings.Split(address, ":")[0]
		if len(hostname) != 2 {
		} else {
			return nil, derrors.NewInvalidArgumentError("server address incorrectly set")
		}

		tlsConfig := &tls.Config{
			ServerName:   hostname,
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
