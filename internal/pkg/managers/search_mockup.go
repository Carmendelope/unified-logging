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

package managers

import (
	"context"

	"github.com/nalej/derrors"

	grpc "github.com/nalej/grpc-unified-logging-go"
)

type MockupSearchManager struct {
}

func NewMockupSearchManager() *MockupSearchManager {
	return &MockupSearchManager{}
}

func (m *MockupSearchManager) Search(ctx context.Context, request *grpc.SearchRequest) (*grpc.LogResponse, derrors.Error) {
	response := &grpc.LogResponse{
		OrganizationId: request.GetOrganizationId(),
		AppInstanceId:  request.GetAppInstanceId(),
		From:           request.GetFrom(),
		To:             request.GetTo(),
		Entries:        []*grpc.LogEntry{},
	}
	return response, nil
}
