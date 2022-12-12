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
	"github.com/apache/incubator-devlake/errors"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertIssueCommentsMeta = core.SubTaskMeta{
	Name:             "convertIssueComments",
	EntryPoint:       ConvertIssueComments,
	EnabledByDefault: true,
	Description:      "ConvertIssueComments data from Gitee api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertIssueComments(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMENTS_TABLE)
	repoId := data.Repo.GiteeId

	cursor, err := db.Cursor(
		dal.From(&models.GiteeIssueComment{}),
		dal.Join("left join _tool_gitee_issues "+
			"on _tool_gitee_issues.gitee_id = _tool_gitee_issue_comments.issue_id"),
		dal.Where("repo_id = ? and _tool_gitee_issues.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.GiteeIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GiteeAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GiteeIssueComment{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			giteeIssueComment := inputRow.(*models.GiteeIssueComment)
			domainIssueComment := &ticket.IssueComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(data.Options.ConnectionId, giteeIssueComment.GiteeId),
				},
				IssueId:     issueIdGen.Generate(data.Options.ConnectionId, giteeIssueComment.IssueId),
				Body:        giteeIssueComment.Body,
				AccountId:   accountIdGen.Generate(data.Options.ConnectionId, giteeIssueComment.AuthorUserId),
				CreatedDate: giteeIssueComment.GiteeCreatedAt,
			}
			return []interface{}{
				domainIssueComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
