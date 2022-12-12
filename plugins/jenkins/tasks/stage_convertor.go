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
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"

	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type JenkinsBuildWithRepoStage struct {
	// collected fields
	ConnectionId        uint64 `gorm:"primaryKey"`
	ID                  string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name                string `json:"name" gorm:"type:varchar(255)"`
	ExecNode            string `json:"execNode" gorm:"type:varchar(255)"`
	CommitSha           string `gorm:"type:varchar(255)"`
	Result              string // Result
	Status              string `json:"status" gorm:"type:varchar(255)"`
	StartTimeMillis     int64  `json:"startTimeMillis"`
	DurationMillis      int    `json:"durationMillis"`
	PauseDurationMillis int    `json:"pauseDurationMillis"`
	Type                string `gorm:"index;type:varchar(255)"`
	BuildName           string `gorm:"primaryKey;type:varchar(255)"`
	Branch              string `gorm:"type:varchar(255)"`
	RepoUrl             string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

var ConvertStagesMeta = core.SubTaskMeta{
	Name:             "convertStages",
	EntryPoint:       ConvertStages,
	EnabledByDefault: true,
	Description:      "convert jenkins_stages",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConvertStages(taskCtx core.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)
	deploymentPattern := data.Options.DeploymentPattern
	productionPattern := data.Options.ProductionPattern
	regexEnricher := helper.NewRegexEnricher()
	err = regexEnricher.AddRegexp(deploymentPattern, productionPattern)
	if err != nil {
		return err
	}
	clauses := []dal.Clause{
		dal.Select(`tjb.connection_id, tjs.build_name, tjs.id, tjs._raw_data_remark, tjs.name,
			tjs._raw_data_id, tjs._raw_data_table, tjs._raw_data_params,
			tjs.status, tjs.start_time_millis, tjs.duration_millis, 
			tjs.pause_duration_millis, tjs.type, 
			tjb.triggered_by, tjb.building`),
		dal.From("_tool_jenkins_stages tjs"),
		dal.Join("left join _tool_jenkins_builds tjb on tjs.build_name = tjb.full_name"),
		dal.Where("tjb.connection_id = ? and tjb.job_path = ? and tjb.job_name = ? ",
			data.Options.ConnectionId, data.Options.JobPath, data.Options.JobName),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	stageIdGen := didgen.NewDomainIdGenerator(&models.JenkinsStage{})
	buildIdGen := didgen.NewDomainIdGenerator(&models.JenkinsBuild{})
	jobIdGen := didgen.NewDomainIdGenerator(&models.JenkinsJob{})

	convertor, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(JenkinsBuildWithRepoStage{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_STAGE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			body := inputRow.(*JenkinsBuildWithRepoStage)
			if body.Name == "" {
				return nil, err
			}
			durationSec := int64(body.DurationMillis / 1000)
			jenkinsTaskResult := ""
			jenkinsTaskStatus := devops.DONE
			var jenkinsTaskFinishedDate *time.Time
			results := make([]interface{}, 0)
			if body.Status == "SUCCESS" {
				jenkinsTaskResult = devops.SUCCESS
			} else if body.Result == "FAILED" {
				jenkinsTaskResult = devops.FAILURE
			} else if body.Result == "ABORTED" {
				jenkinsTaskResult = devops.ABORT
			} else {
				jenkinsTaskResult = ""
				jenkinsTaskStatus = devops.IN_PROGRESS
			}

			startedDate := time.Unix(body.StartTimeMillis/1000, 0)
			finishedDate := startedDate.Add(time.Duration(durationSec * int64(time.Second)))
			jenkinsTaskFinishedDate = &finishedDate
			jenkinsTask := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: stageIdGen.Generate(body.ConnectionId, body.BuildName, body.ID),
				},
				Name:         body.Name,
				PipelineId:   buildIdGen.Generate(body.ConnectionId, body.BuildName),
				Result:       jenkinsTaskResult,
				Status:       jenkinsTaskStatus,
				DurationSec:  uint64(body.DurationMillis / 1000),
				StartedDate:  time.Unix(durationSec, 0),
				FinishedDate: jenkinsTaskFinishedDate,
				CicdScopeId:  jobIdGen.Generate(body.ConnectionId, data.Options.JobFullName),
			}
			jenkinsTask.Type = regexEnricher.GetEnrichResult(deploymentPattern, jenkinsTask.Name, devops.DEPLOYMENT)
			jenkinsTask.Environment = regexEnricher.GetEnrichResult(productionPattern, jenkinsTask.Name, devops.PRODUCTION)
			jenkinsTask.RawDataOrigin = body.RawDataOrigin

			results = append(results, jenkinsTask)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
