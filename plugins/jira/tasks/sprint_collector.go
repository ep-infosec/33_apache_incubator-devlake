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
	"github.com/apache/incubator-devlake/errors"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_SPRINT_TABLE = "jira_api_sprints"

var _ core.SubTaskEntryPoint = CollectSprints

var CollectSprintsMeta = core.SubTaskMeta{
	Name:             "collectSprints",
	EntryPoint:       CollectSprints,
	EnabledByDefault: true,
	Description:      "collect Jira sprints",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectSprints(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect sprints")
	jql := "ORDER BY created ASC"
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_SPRINT_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    50,
		UrlTemplate: "agile/1.0/board/{{ .Params.BoardId }}/sprint",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("jql", jql)
			query.Set("startAt", fmt.Sprintf("%v", reqData.Pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Values []json.RawMessage `json:"values"`
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Values, nil
		},
		AfterResponse: ignoreHTTPStatus400,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
