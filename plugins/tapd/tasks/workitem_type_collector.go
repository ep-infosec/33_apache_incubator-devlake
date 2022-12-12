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

const RAW_WORKITEM_TYPE_TABLE = "tapd_api_workitem_types"

var _ core.SubTaskEntryPoint = CollectWorkitemTypes

func CollectWorkitemTypes(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_WORKITEM_TYPE_TABLE, false)
	logger := taskCtx.GetLogger()
	logger.Info("collect workitem_type")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "workitem_types",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				WorkitemTypes []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.WorkitemTypes, err
		},
	})
	if err != nil {
		logger.Error(err, "collect workitem_type error")
		return err
	}
	return collector.Execute()
}

var CollectWorkitemTypesMeta = core.SubTaskMeta{
	Name:             "collectWorkitemTypes",
	EntryPoint:       CollectWorkitemTypes,
	EnabledByDefault: true,
	Description:      "collect Tapd WorkitemTypes",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}
