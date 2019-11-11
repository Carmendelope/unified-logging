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

// Search request structure for logging provider

package entities

import (
	"fmt"
	"time"
)

// SearchFilter is a mapping of keys to arrays of values - these will be used
// to match the logging metadata. If a log line matches a key to one or more
// of the values, the line is returned. Depending on IsUnionFilter below,
// all keys should match, or at least one key should match
type SearchFilter map[Field][]string

// SortOrder defines the ordering of the search result based on timestamp
// Defaults to ascending (oldest first)
type SortOrder int

const (
	Ascending  SortOrder = 0
	Descending SortOrder = 1
)

func (s SortOrder) ToAscending() bool {
	if s == Ascending {
		return true
	}
	return false
}

// SearchRequest is the structure that is used to describe the search query for the logging storage provider
type SearchRequest struct {
	// Filters is a map with all the fields that you want to include, the filter is a exact filter.
	// More than one filter will result in a query that's the intersection
	// of all the filters (AND)
	Filters SearchFilter
	// Indicates to treat mutiple filters as a union (OR) instead of intersection
	IsUnionFilter bool
	// MsgFilter is a string that filters the log entries by message text. It allows wildcards.
	MsgFilter string
	// from is the beginning date in Unix time format.
	From time.Time
	// to is the ending date in Unix time format.
	To time.Time
	// Order specifies the timestamp sort ordering
	// Defaults to ascending
	Order SortOrder
}

// IsValid check if the search request is well-formed.
func (e *SearchRequest) IsValid() bool {
	// Even an empty request is valid - will just return everything
	return true
}

// String returns the string representation of the search request
func (e *SearchRequest) String() string {
	return fmt.Sprintf("%#v", e)
}
