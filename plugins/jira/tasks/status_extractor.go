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
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var ExtractStatusMeta = core.SubTaskMeta{
	Name:             "extractStatus",
	EntryPoint:       ExtractStatus,
	EnabledByDefault: true,
	Description:      "extract Jira status",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractStatus(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract Status, connection_id=%d, board_id=%d", connectionId, boardId)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_STATUS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var apiStatus apiv2models.Status
			err := errors.Convert(json.Unmarshal(row.Data, &apiStatus))
			if err != nil {
				return nil, err
			}
			if apiStatus.Scope != nil {
				// FIXME: skip scope status
				return nil, nil
			}
			var jiraStatus = &models.JiraStatus{
				ConnectionId:   connectionId,
				ID:             apiStatus.ID,
				Name:           apiStatus.Name,
				Self:           apiStatus.Self,
				StatusCategory: apiStatus.StatusCategory.Key,
			}
			var result []interface{}
			result = append(result, jiraStatus)
			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
