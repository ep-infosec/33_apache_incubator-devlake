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

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
)

const RAW_REPOSITORIES_TABLE = "azure_api_repositories"

var CollectApiRepoMeta = core.SubTaskMeta{
	Name:        "collectApiRepo",
	EntryPoint:  CollectApiRepositories,
	Required:    true,
	Description: "Collect repositories data from Azure api",
	DomainTypes: []string{core.DOMAIN_TYPE_CODE},
}

func CollectApiRepositories(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*AzureTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AzureApiParams{
				ConnectionId: data.Options.ConnectionId,
				Project:      data.Options.Project,
			},
			Table: RAW_REPOSITORIES_TABLE,
		},
		ApiClient: data.ApiClient,

		UrlTemplate: "{{ .Params.Project }}/_apis/git/repositories?api-version=7.1-preview.1",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Repos []json.RawMessage `json:"value"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Repos, err
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
