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
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractSubWorkspaces

var ExtractSubWorkspaceMeta = core.SubTaskMeta{
	Name:             "extractSubWorkspaces",
	EntryPoint:       ExtractSubWorkspaces,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_workspaces",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractSubWorkspaces(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_SUB_WORKSPACE_TABLE, false)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var subWorkspaceRes struct {
				Workspace models.TapdSubWorkspace
			}
			err := errors.Convert(json.Unmarshal(row.Data, &subWorkspaceRes))
			if err != nil {
				return nil, err
			}

			ws := subWorkspaceRes.Workspace
			ws.ParentId = data.Options.WorkspaceId
			ws.ConnectionId = data.Options.ConnectionId
			return []interface{}{
				&ws,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
