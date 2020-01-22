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

// Logging storage provider interface

package loggingstorage

import (
	"context"

	"github.com/nalej/derrors"
	"github.com/nalej/unified-logging/pkg/entities"
)

// Provider is the interface of the Logging provider.
type Provider interface {
	Search(ctx context.Context, request *entities.SearchRequest, limit int) (entities.LogEntries, derrors.Error)
	Expire(ctx context.Context, request *entities.SearchRequest) derrors.Error
	RemoveIndex(ctx context.Context, index string) derrors.Error
}
