/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"gorm.io/gorm"
	"time"
)

// ApiCollectorStateManager save collector state in framework table
type ApiCollectorStateManager struct {
	RawDataSubTaskArgs
	*ApiCollector
	LatestState      models.CollectorLatestState
	CreatedDateAfter *time.Time
	ExecuteStart     time.Time
}

// NewApiCollectorWithState create a new ApiCollectorStateManager
func NewApiCollectorWithState(args RawDataSubTaskArgs, createdDateAfter *time.Time) (*ApiCollectorStateManager, errors.Error) {
	db := args.Ctx.GetDal()

	rawDataSubTask, err := NewRawDataSubTask(args)
	if err != nil {
		return nil, errors.Default.Wrap(err, "Couldn't resolve raw subtask args")
	}
	latestState := models.CollectorLatestState{}
	err = db.First(&latestState, dal.Where(`raw_data_table = ? AND raw_data_params = ?`, rawDataSubTask.table, rawDataSubTask.params))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			latestState = models.CollectorLatestState{
				RawDataTable:  rawDataSubTask.table,
				RawDataParams: rawDataSubTask.params,
			}
		} else {
			return nil, errors.Default.Wrap(err, "failed to load JiraLatestCollectorMeta")
		}
	}
	return &ApiCollectorStateManager{
		RawDataSubTaskArgs: args,
		LatestState:        latestState,
		CreatedDateAfter:   createdDateAfter,
		ExecuteStart:       time.Now(),
	}, nil
}

// CanIncrementCollect return if the old data can support collect incrementally.
// only when latest collection is success &&
// (m.LatestState.CreatedDateAfter == nil means all data have been collected ||
// CreatedDateAfter at this time exists and later than in the LatestState)
// if CreatedDateAfter at this time not exists, collect incrementally only when "m.LatestState.CreatedDateAfter == nil"
func (m ApiCollectorStateManager) CanIncrementCollect() bool {
	return m.LatestState.LatestSuccessStart != nil &&
		(m.LatestState.CreatedDateAfter == nil || m.CreatedDateAfter != nil && m.CreatedDateAfter.After(*m.LatestState.CreatedDateAfter))
}

// InitCollector init the embedded collector
func (m *ApiCollectorStateManager) InitCollector(args ApiCollectorArgs) (err errors.Error) {
	args.RawDataSubTaskArgs = m.RawDataSubTaskArgs
	m.ApiCollector, err = NewApiCollector(args)
	return err
}

// Execute the embedded collector and record execute state
func (m ApiCollectorStateManager) Execute() errors.Error {
	err := m.ApiCollector.Execute()
	if err != nil {
		return err
	}

	db := m.Ctx.GetDal()
	m.LatestState.LatestSuccessStart = &m.ExecuteStart
	m.LatestState.CreatedDateAfter = m.CreatedDateAfter
	err = db.CreateOrUpdate(&m.LatestState)
	return err
}
