/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Unified logging component for application clusters.
// Executes search queries against local Elastic

package main

import (
	"github.com/nalej/unified-logging/cmd/unified-logging-slave/commands"
	"github.com/nalej/golang-template/version"
)

var MainVersion string
var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
