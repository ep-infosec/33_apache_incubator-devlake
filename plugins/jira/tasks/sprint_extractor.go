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
	"github.com/apache/incubator-devlake/plugins/jira/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractSprints

var ExtractSprintsMeta = core.SubTaskMeta{
	Name:             "extractSprints",
	EntryPoint:       ExtractSprints,
	EnabledByDefault: true,
	Description:      "extract Jira sprints",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractSprints(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_SPRINT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var sprint apiv2models.Sprint
			err := errors.Convert(json.Unmarshal(row.Data, &sprint))
			if err != nil {
				return nil, err
			}
			boardSprint := models.JiraBoardSprint{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
				SprintId:     sprint.ID,
			}
			return []interface{}{sprint.ToToolLayer(data.Options.ConnectionId), &boardSprint}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
