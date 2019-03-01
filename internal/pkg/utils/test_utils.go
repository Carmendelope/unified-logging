/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package utils

import (
	"os"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// RunIntegrationTests checks whether integration tests should be executed.
func RunIntegrationTests() bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	return runIntegration == "true"
}

func MatchTimestamp(t *timestamp.Timestamp) types.GomegaMatcher {
	return gstruct.PointTo(
		gstruct.MatchFields(gstruct.IgnoreExtras,
			gstruct.Fields{
				"Seconds": gomega.Equal(t.Seconds),
				"Nanos": gomega.Equal(t.Nanos),
			},
		),
	)
}

