/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Log-line entries for logging provider

package entities

type LogEntries []*LogEntry

type LogEntry struct {
	Timestamp int64
	Msg string
}
