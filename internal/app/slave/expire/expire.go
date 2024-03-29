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
	"github.com/nalej/unified-logging/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"regexp"
	"time"

	"github.com/nalej/derrors"

	"github.com/nalej/unified-logging/pkg/entities"
	"github.com/nalej/unified-logging/pkg/provider/loggingstorage"

	"github.com/nalej/grpc-common-go"
	grpc "github.com/nalej/grpc-unified-logging-go"
)

// expire logs timer
const LoopSleep = time.Minute * 60 * 24

// DefaultLogEntryTTL number of days the logs will be alive in the system, then they will be deleted
const DefaultLogEntryTTL = 7

// indexPattern is a regular expression to find the date in the index name "YYYY.MM.DD"
const indexPattern = "\\d{4}.\\d{2}.\\d{2}$"

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

	// the index name is like "filebeat-6.6.0-2020.01.17", we need to find the date of the end
	re := regexp.MustCompile(indexPattern)
	ind := re.FindStringSubmatch(index)
	if len(ind) <= 0 {
		return false, derrors.NewInternalError("error parsing the index").WithParams(index)
	}

	indexDate := ind[0]
	date, err := time.Parse("2006.01.02", indexDate)
	if err != nil {
		return false, derrors.NewInternalError("error checking the index")
	}
	limitDate := time.Now().AddDate(0, 0, -1*(DefaultLogEntryTTL+1)) // I need to sum one day because I am comparing now with time and the index date without it
	if limitDate.After(date) {
		return true, nil
	}
	return false, nil
}

// deleteIndex gets all the indexes and removes the old ones
func (m *Manager) deleteIndex() {
	log.Debug().Msg("Delete Index")

	listCtx, listCancel := utils.GetContext()
	defer listCancel()
	indexList, err := m.Provider.GetIndexList(listCtx)
	if err != nil {
		log.Warn().Str("err", err.DebugReport()).Msg("error cleaning index")
		return
	}

	for _, index := range indexList {
		remove, err := m.checkRemoveIndex(index)
		if err != nil {
			log.Warn().Str("index", index).Msg("error checking the index")
		} else {
			log.Debug().Str("index", index).Bool("remove", remove).Msg("checking the index")
			if remove {
				ctx, cancel := utils.GetContext()
				defer cancel()
				err = m.Provider.RemoveIndex(ctx, index)
				if err != nil {
					log.Warn().Str("index", index).Str("err", err.DebugReport()).Msg("error cleaning index")
				}
			}
		}

	}
}

// DeleteIndexLoop Loop to remove old indexes
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
