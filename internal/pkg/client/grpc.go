/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Application cluster unified logging GRPC client

package client

import (
	"github.com/nalej/grpc-app-cluster-api-go"

        "google.golang.org/grpc"
)

type GRPCLoggingClient struct {
	grpc_app_cluster_api_go.UnifiedLoggingClient
	conn *grpc.ClientConn
}

func NewGRPCLoggingClient(address string) (LoggingClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := grpc_app_cluster_api_go.NewUnifiedLoggingClient(conn)

	return &GRPCLoggingClient{client, conn}, nil
}

func (c *GRPCLoggingClient) Close() error {
	return c.conn.Close()
}
