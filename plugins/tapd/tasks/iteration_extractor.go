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

package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractIterations

var ExtractIterationMeta = core.SubTaskMeta{
	Name:             "extractIterations",
	EntryPoint:       ExtractIterations,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractIterations(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ITERATION_TABLE, false)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var iterBody struct {
				Iteration models.TapdIteration
			}
			err := errors.Convert(json.Unmarshal(row.Data, &iterBody))
			if err != nil {
				return nil, err
			}
			iter := iterBody.Iteration

			iter.ConnectionId = data.Options.ConnectionId
			workspaceIter := &models.TapdWorkspaceIteration{
				ConnectionId: data.Options.ConnectionId,
				WorkspaceId:  iter.WorkspaceId,
				IterationId:  iter.Id,
			}
			return []interface{}{
				&iter, workspaceIter,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
