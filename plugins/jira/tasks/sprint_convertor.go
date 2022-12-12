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
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var ConvertSprintsMeta = core.SubTaskMeta{
	Name:             "convertSprints",
	EntryPoint:       ConvertSprints,
	EnabledByDefault: true,
	Description:      "convert Jira sprints",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertSprints(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert sprints")
	clauses := []dal.Clause{
		dal.Select("tjs.*"),
		dal.From("_tool_jira_sprints tjs"),
		dal.Join(`LEFT JOIN _tool_jira_board_sprints tjbs
              ON tjbs.sprint_id = tjs.sprint_id
                 AND tjbs.connection_id = tjs.connection_id`),
		dal.Where("tjs.connection_id = ? AND tjbs.board_id = ?", connectionId, boardId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	var converter *helper.DataConverter
	domainBoardId := didgen.NewDomainIdGenerator(&models.JiraBoard{}).Generate(connectionId, boardId)
	sprintIdGen := didgen.NewDomainIdGenerator(&models.JiraSprint{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.JiraBoard{})
	converter, err = helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_SPRINT_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraSprint{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			var result []interface{}
			jiraSprint := inputRow.(*models.JiraSprint)
			sprint := &ticket.Sprint{
				DomainEntity:    domainlayer.DomainEntity{Id: sprintIdGen.Generate(connectionId, jiraSprint.SprintId)},
				Url:             jiraSprint.Self,
				Status:          strings.ToUpper(jiraSprint.State),
				Name:            jiraSprint.Name,
				StartedDate:     jiraSprint.StartDate,
				EndedDate:       jiraSprint.EndDate,
				CompletedDate:   jiraSprint.CompleteDate,
				OriginalBoardID: boardIdGen.Generate(connectionId, jiraSprint.OriginBoardID),
			}
			result = append(result, sprint)
			boardSprint := &ticket.BoardSprint{
				BoardId:  domainBoardId,
				SprintId: sprint.Id,
			}
			result = append(result, boardSprint)
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
