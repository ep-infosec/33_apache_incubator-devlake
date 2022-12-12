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
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func ConvertStoryCommit(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_COMMIT_TABLE, false)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.TapdStoryCommit{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	issueIdGen := didgen.NewDomainIdGenerator(&models.TapdStory{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdStoryCommit{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolL := inputRow.(*models.TapdStoryCommit)
			domainL := &crossdomain.IssueCommit{
				IssueId:   issueIdGen.Generate(data.Options.ConnectionId, toolL.StoryId),
				CommitSha: toolL.CommitId,
			}

			return []interface{}{
				domainL,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertStoryCommitMeta = core.SubTaskMeta{
	Name:             "convertStoryCommit",
	EntryPoint:       ConvertStoryCommit,
	EnabledByDefault: true,
	Description:      "convert Tapd StoryCommit",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}
