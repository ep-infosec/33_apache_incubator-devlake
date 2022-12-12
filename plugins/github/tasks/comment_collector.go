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
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

const RAW_COMMENTS_TABLE = "github_api_comments"

func CollectApiComments(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	since := data.Since
	incremental := false
	var err errors.Error
	if since == nil {
		since, incremental, err = calculateSince(data, db)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_COMMENTS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues/comments",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
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
		return errors.Default.Wrap(err, "error collecting github comments")
	}

	return collector.Execute()
}

var CollectApiCommentsMeta = core.SubTaskMeta{
	Name:             "collectApiComments",
	EntryPoint:       CollectApiComments,
	EnabledByDefault: true,
	Description:      "Collect comments data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE_REVIEW, core.DOMAIN_TYPE_TICKET},
}

func calculateSince(data *GithubTaskData, db dal.Dal) (*time.Time, bool, errors.Error) {
	var since *time.Time
	var latestUpdatedIssueComt models.GithubIssueComment
	var latestUpdatedPrComt models.GithubPrComment

	incremental := false
	err := db.All(
		&latestUpdatedIssueComt,
		dal.Join("left join _tool_github_issues on _tool_github_issues.github_id = _tool_github_issue_comments.issue_id"),
		dal.Where(
			"_tool_github_issues.repo_id = ? AND _tool_github_issues.connection_id = ?", data.Repo.GithubId, data.Repo.ConnectionId,
		),
		dal.Orderby("github_updated_at DESC"),
		dal.Limit(1),
	)
	if err != nil {
		return nil, false, errors.Default.Wrap(err, "failed to get latest github issue record")
	}

	err = db.All(
		&latestUpdatedPrComt,
		dal.Join("left join _tool_github_pull_requests on _tool_github_pull_requests.github_id = _tool_github_pull_request_comments.pull_request_id"),
		dal.Where("_tool_github_pull_requests.repo_id = ? AND _tool_github_pull_requests.connection_id = ?", data.Repo.GithubId, data.Repo.ConnectionId),
		dal.Orderby("github_updated_at DESC"),
		dal.Limit(1),
	)
	if err != nil {
		return nil, false, errors.Default.Wrap(err, "failed to get latest github issue record")
	}
	if latestUpdatedIssueComt.GithubId > 0 && latestUpdatedPrComt.GithubId > 0 {
		if latestUpdatedIssueComt.GithubUpdatedAt.Before(latestUpdatedPrComt.GithubUpdatedAt) {
			since = &latestUpdatedPrComt.GithubUpdatedAt
		} else {
			since = &latestUpdatedIssueComt.GithubUpdatedAt
		}
		incremental = true
	} else if latestUpdatedIssueComt.GithubId > 0 {
		since = &latestUpdatedIssueComt.GithubUpdatedAt
		incremental = true
	} else if latestUpdatedPrComt.GithubId > 0 {
		since = &latestUpdatedPrComt.GithubUpdatedAt
		incremental = true
	}
	return since, incremental, nil
}
