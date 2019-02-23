/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Validate unified logging requests

package handler

import (
        "github.com/nalej/derrors"

	grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/rs/zerolog/log"
)

const emptyOrganizationId = "organization_id cannot be empty"

// This is an interface with the methods that are indentical for search
// and expire requests, such that we can validate them in the same function
type LoggingRequest interface {
	String() string
	GetOrganizationId() string
	GetAppInstanceId() string
}

func validate(request LoggingRequest) derrors.Error {
	log.Debug().Str("request", request.String()).Msg("validating incoming request")

	if request.GetOrganizationId() == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}

	return nil
}

func validateSearch(request *grpc.SearchRequest) derrors.Error {
	return validate(request)
}

func validateExpire(request *grpc.ExpirationRequest) derrors.Error {
	return validate(request)
}
