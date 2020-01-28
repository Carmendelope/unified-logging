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

// Expire manager for unified logging slave

package expire

import (
	"context"
	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"github.com/nalej/derrors"

	"github.com/nalej/unified-logging/pkg/entities"
	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

	"github.com/nalej/grpc-common-go"
	grpc "github.com/nalej/grpc-unified-logging-go"
)

// expire logs timer

const LoopSleep = time.Minute * 60 * 24

const defaultDaysExpired = 7
const indexPrefix = "filebeat-6.6.0-"

type Manager struct {
	Provider loggingstorage.Provider
}

func NewManager(provider loggingstorage.Provider) *Manager {
	return &Manager{
		Provider: provider,
	}
}

func (m *Manager) Expire(ctx context.Context, request *grpc.ExpirationRequest) (*grpc_common_go.Success, derrors.Error) {
	// We have a verified request - translate to entities.SearchRequest and execute
	fields := entities.FilterFields{
		OrganizationId: request.GetOrganizationId(),
		AppInstanceId:  request.GetAppInstanceId(),
	}

	search := &entities.SearchRequest{
		Filters:       fields.ToFilters(),
		IsUnionFilter: false,
	}

	err := m.Provider.Expire(ctx, search)
	if err != nil {
		return nil, err
	}

	return &grpc_common_go.Success{}, nil
}

// check if the index must be deleted
func (m *Manager) checkRemoveIndex(index string) (bool, derrors.Error) {
	// index example "filebeat-6.6.0-2020.01.17"
	indexDate := strings.Replace(index, indexPrefix, "", -1)

	date, err := time.Parse("2006.01.02", indexDate)
	if err != nil {
		return false, derrors.NewInternalError("error checking the index")
	}
	limitDate := time.Now().AddDate(0,0,-1 * (defaultDaysExpired +1)) // I need to sum one day because I am comparing now with time and the index date without it
	if limitDate.After(date) {
		return true, nil
	}
	return false, nil
}

func (m *Manager) deleteIndex() {
	log.Debug().Msg("Delete Index")

	indexList, err := m.Provider.GetIndexList(context.Background())
	if err != nil {
		log.Warn().Str("err", err.DebugReport()).Msg("error cleaning index")
		return
	}
	log.Debug().Interface("indexList", indexList).Msg("listing the indexes")

	for _, index := range indexList {
		remove, err := m.checkRemoveIndex(index)
		if err != nil {
			log.Debug().Str("index", index).Msg("error checking the index")
		} else {
			log.Debug().Str("index", index).Bool("remove", remove).Msg("checking the index")
			if remove{
				err = m.Provider.RemoveIndex(context.Background(), index)
				if err != nil {
					log.Warn().Str("index", index).Str("err", err.DebugReport()).Msg("error cleaning index")
				}
			}
		}

	}
}

func (m *Manager) DeleteIndexLoop() {
	log.Debug().Msg("Delete Index Loop Begins")
	ticker := time.NewTicker(LoopSleep)
	for {
		select {
		case <-ticker.C:
			m.deleteIndex()
		}
	}
}
