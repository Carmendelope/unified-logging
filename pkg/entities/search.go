/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

// Search request structure for logging provider

package entities

import "fmt"

// SearchFilter is a mapping of keys to arrays of values - these will be used
// to match the logging metadata. If a log line matches a key to one or more
// of the values, the line is returned. Depending on IsUnionFilter below,
// all keys should match, or at least one key should match
type SearchFilter map[Field][]string

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
	From int64
	// to is the ending date in Unix time format.
	To int64
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
