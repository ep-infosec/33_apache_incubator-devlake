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
	goerror "errors"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"gorm.io/gorm"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_ISSUE_TABLE = "gitlab_api_issues"

var CollectApiIssuesMeta = core.SubTaskMeta{
	Name:             "collectApiIssues",
	EntryPoint:       CollectApiIssues,
	EnabledByDefault: true,
	Description:      "Collect issues data from Gitlab api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectApiIssues(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	if since == nil {
		var latestUpdated models.GitlabIssue
		clause := []dal.Clause{
			dal.Orderby("gitlab_updated_at DESC"),
		}
		err := db.First(&latestUpdated, clause...)
		if err != nil && !goerror.Is(err, gorm.ErrRecordNotFound) {
			return errors.Default.Wrap(err, "failed to get latest gitlab issue record")
		}
		if latestUpdated.GitlabId > 0 {
			since = &latestUpdated.GitlabUpdatedAt
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        incremental,

		UrlTemplate: "projects/{{ .Params.ProjectId }}/issues",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			if since != nil {
				query.Set("updated_after", since.String())
			}
			query.Set("sort", "asc")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},

		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var items []json.RawMessage
			err := helper.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
