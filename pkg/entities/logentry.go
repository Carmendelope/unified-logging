/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Log-line entries for logging provider

package entities

import (
	"time"
)

type LogEntries []*LogEntry

type LogEntry struct {
	Timestamp time.Time `json:"@timestamp"`
	Msg string `json:"message"`
}
