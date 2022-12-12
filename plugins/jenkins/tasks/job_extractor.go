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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

// this struct should be moved to `gitub_api_common.go`

var ExtractApiJobsMeta = core.SubTaskMeta{
	Name:             "extractApiJobs",
	EntryPoint:       ExtractApiJobs,
	EnabledByDefault: true,
	Description:      "Extract raw jobs data into tool layer table jenkins_jobs",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ExtractApiJobs(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JenkinsTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			body := &models.Job{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 1+len(body.UpstreamProjects))

			job := &models.JenkinsJob{
				ConnectionId: data.Options.ConnectionId,
				FullName:     body.FullName,
				Name:         body.Name,
				Path:         data.Options.JobPath,
				Class:        body.Class,
				Color:        body.Color,
				Base:         body.Base,
				Url:          body.URL,
				Description:  body.Description,
				NoPKModel:    common.NoPKModel{},
			}
			if body.PrimaryView != nil {
				job.PrimaryView = body.PrimaryView.Name
			}
			for _, upstreamProject := range body.UpstreamProjects {
				upDownJob := models.JenkinsJobDag{
					ConnetionId:   data.Options.ConnectionId,
					UpstreamJob:   upstreamProject.Name,
					DownstreamJob: job.Name,
				}
				results = append(results, &upDownJob)
			}
			data.Options.JobFullName = body.FullName
			results = append(results, job)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
