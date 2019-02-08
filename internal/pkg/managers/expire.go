/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package managers

import (
        "github.com/nalej/derrors"

        grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/grpc-common-go"
)

// Interface for Expire Manager
type Expire interface {
	Expire(*grpc.ExpirationRequest) (*grpc_common_go.Success, derrors.Error)
}
