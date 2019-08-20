/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Application cluster unified logging GRPC client

package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nalej/grpc-app-cluster-api-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/rs/zerolog/log"
)

type GRPCLoggingClient struct {
	grpc_app_cluster_api_go.UnifiedLoggingClient
	conn *grpc.ClientConn
}

func NewGRPCLoggingClient(address string, params *LoggingClientParams) (LoggingClient, error) {
	var options []grpc.DialOption

	log.Debug().Str("address", address).
		Bool("tls", params.UseTLS).
		Str("cert", params.CACert).
		Bool("skipServerCertValidation", params.SkipServerCertValidation).
		Msg("creating connection")

	if params.UseTLS {
		rootCAs := x509.NewCertPool()
		if params.CACert != "" {
			err := addCert(rootCAs, params.CACert)
			if err != nil {
				return nil, err
			}
		}

		tlsConfig := &tls.Config{
			RootCAs: rootCAs,
			ServerName: strings.Split(address, ":")[0],
			InsecureSkipVerify: params.SkipServerCertValidation,
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
