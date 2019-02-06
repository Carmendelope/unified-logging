/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Unified logging component for management cluster.
// Executes and aggregates search queries against application clusters

package main

import (
	"github.com/nalej/unified-logging/cmd/unified-logging-coord/commands"
	"github.com/nalej/golang-template/version"
)

var MainVersion string
var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
