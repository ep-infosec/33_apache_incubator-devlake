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
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
	"regexp"
)

var ExtractRemotelinksMeta = core.SubTaskMeta{
	Name:             "extractRemotelinks",
	EntryPoint:       ExtractRemotelinks,
	EnabledByDefault: true,
	Description:      "extract Jira remote links",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractRemotelinks(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("extract remote links")
	var commitShaRegex *regexp.Regexp
	var err error
	if data.Options.TransformationRules != nil && data.Options.TransformationRules.RemotelinkCommitShaPattern != "" {
		pattern := data.Options.TransformationRules.RemotelinkCommitShaPattern
		commitShaRegex, err = regexp.Compile(pattern)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile pattern failed")
		}
	}

	// select all remotelinks belongs to the board, cursor is important for low memory footprint
	clauses := []dal.Clause{
		dal.From(&models.JiraRemotelink{}),
		dal.Select("*"),
		dal.Join("left join _tool_jira_board_issues on _tool_jira_board_issues.issue_id = _tool_jira_remotelinks.issue_id"),
		dal.Where("_tool_jira_board_issues.board_id = ? AND _tool_jira_board_issues.connection_id = ?", boardId, connectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return errors.Convert(err)
	}
	defer cursor.Close()

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_REMOTELINK_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var result []interface{}
			var raw apiv2models.RemoteLink
			err := errors.Convert(json.Unmarshal(row.Data, &raw))
			if err != nil {
				return nil, err
			}
			var input apiv2models.Input
			err = errors.Convert(json.Unmarshal(row.Input, &input))
			if err != nil {
				return nil, err
			}
			remotelink := &models.JiraRemotelink{
				ConnectionId: connectionId,
				RemotelinkId: raw.ID,
				IssueId:      input.IssueId,
				Self:         raw.Self,
				Title:        raw.Object.Title,
				Url:          raw.Object.URL,
				IssueUpdated: &input.UpdateTime,
			}
			result = append(result, remotelink)
			if commitShaRegex != nil {
				groups := commitShaRegex.FindStringSubmatch(remotelink.Url)
				if len(groups) > 1 {
					issueCommit := &models.JiraIssueCommit{
						ConnectionId: connectionId,
						IssueId:      remotelink.IssueId,
						CommitSha:    groups[1],
						CommitUrl:    remotelink.Url,
					}
					result = append(result, issueCommit)
				}
			}
			return result, nil
		},
	})
	if err != nil {
		return errors.Convert(err)
	}

	return extractor.Execute()
}
