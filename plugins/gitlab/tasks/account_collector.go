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
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_USER_TABLE = "gitlab_api_users"

var CollectAccountsMeta = core.SubTaskMeta{
	Name:             "collectAccounts",
	EntryPoint:       CollectAccounts,
	EnabledByDefault: true,
	Description:      "collect gitlab users",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func CollectAccounts(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect gitlab users")

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "/projects/{{ .Params.ProjectId }}/members/all",
		//PageSize:           100,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			// query.Set("sort", "asc")
			// query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			// query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},

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
		logger.Error(err, "collect user error")
		return err
	}

	return collector.Execute()
}
