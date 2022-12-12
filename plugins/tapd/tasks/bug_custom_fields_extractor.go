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

var _ core.SubTaskEntryPoint = ExtractBugCustomFields

var ExtractBugCustomFieldsMeta = core.SubTaskMeta{
	Name:             "extractBugCustomFields",
	EntryPoint:       ExtractBugCustomFields,
	EnabledByDefault: true,
	Description:      "Extract raw company data into tool layer table _tool_tapd_bug_custom_fields",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractBugCustomFields(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_CUSTOM_FIELDS_TABLE, false)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var bugCustomFields struct {
				CustomFieldConfig models.TapdBugCustomFields
			}
			err := errors.Convert(json.Unmarshal(row.Data, &bugCustomFields))
			if err != nil {
				return nil, err
			}

			toolL := bugCustomFields.CustomFieldConfig

			toolL.ConnectionId = data.Options.ConnectionId
			return []interface{}{
				&toolL,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
