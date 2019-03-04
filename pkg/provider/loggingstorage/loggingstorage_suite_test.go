/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package loggingstorage

import (
        "github.com/onsi/ginkgo"
        "github.com/onsi/gomega"
        "testing"
)

func TestHandlerPackage(t *testing.T) {
        gomega.RegisterFailHandler(ginkgo.Fail)
        ginkgo.RunSpecs(t, "Loggingstorage package suite")
}