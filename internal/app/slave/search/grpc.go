/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// GRPC utility functions

package search

import (
	grpc "github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/unified-logging/pkg/entities"
)

func GRPCEntries(entries entities.LogEntries) []*grpc.LogEntry {
	result := make([]*grpc.LogEntry, len(entries))
	for i, e := range entries {
		result[i] = &grpc.LogEntry{
			Timestamp: e.Timestamp.Unix(),
			Msg:       e.Msg,
		}
	}
	return result
}
