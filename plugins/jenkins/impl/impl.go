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

package impl

import (
	"fmt"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jenkins/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/jenkins/tasks"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Jenkins)(nil)
var _ core.PluginInit = (*Jenkins)(nil)
var _ core.PluginTask = (*Jenkins)(nil)
var _ core.PluginApi = (*Jenkins)(nil)
var _ core.PluginModel = (*Jenkins)(nil)
var _ core.PluginMigration = (*Jenkins)(nil)
var _ core.CloseablePluginTask = (*Jenkins)(nil)
var _ core.PluginSource = (*Jenkins)(nil)

type Jenkins struct{}

func (plugin Jenkins) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) errors.Error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Jenkins) Connection() interface{} {
	return &models.JenkinsConnection{}
}

func (plugin Jenkins) Scope() interface{} {
	return &models.JenkinsJob{}
}

func (plugin Jenkins) TransformationRule() interface{} {
	return &models.JenkinsTransformationRule{}
}

func (plugin Jenkins) GetTablesInfo() []core.Tabler {
	return []core.Tabler{
		&models.JenkinsBuild{},
		&models.JenkinsBuildCommit{},
		&models.JenkinsConnection{},
		&models.JenkinsJob{},
		&models.JenkinsJobDag{},
		&models.JenkinsPipeline{},
		&models.JenkinsResponse{},
		&models.JenkinsStage{},
		&models.JenkinsTask{},
	}
}

func (plugin Jenkins) Description() string {
	return "To collect and enrich data from Jenkins"
}

func (plugin Jenkins) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectApiJobsMeta,
		tasks.ExtractApiJobsMeta,
		tasks.ConvertJobsMeta,
		tasks.CollectApiBuildsMeta,
		tasks.ExtractApiBuildsMeta,
		tasks.CollectApiStagesMeta,
		tasks.ExtractApiStagesMeta,
		tasks.EnrichApiBuildWithStagesMeta,
		tasks.ConvertBuildsToCICDMeta,
		tasks.ConvertStagesMeta,
		tasks.ConvertBuildReposMeta,
	}
}
func (plugin Jenkins) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	logger := taskCtx.GetLogger()
	logger.Debug("%v", options)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}

	connection := &models.JenkinsConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}

	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse("2006-01-02T15:04:05Z", op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}

	if !strings.HasSuffix(op.JobPath, "/") {
		op.JobPath = fmt.Sprintf("%s/", op.JobPath)
	}

	if op.JenkinsTransformationRule == nil && op.TransformationRuleId != 0 {
		var transformationRule models.JenkinsTransformationRule
		err = taskCtx.GetDal().First(&transformationRule, dal.Where("id = ?", op.TransformationRuleId))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "fail to get transformationRule")
		}
		op.JenkinsTransformationRule = &transformationRule
	}
	taskData := &tasks.JenkinsTaskData{
		Options:    op,
		ApiClient:  apiClient,
		Connection: connection,
	}
	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logger.Debug("collect data created from %s", createdDateAfter)
	}
	return taskData, nil
}

func (plugin Jenkins) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/jenkins"
}

func (plugin Jenkins) MigrationScripts() []core.MigrationScript {
	return migrationscripts.All()
}

func (plugin Jenkins) MakePipelinePlan(connectionId uint64, scope []*core.BlueprintScopeV100) (core.PipelinePlan, errors.Error) {
	return api.MakePipelinePlanV100(plugin.SubTaskMetas(), connectionId, scope)
}

func (plugin Jenkins) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*core.BlueprintScopeV200) (pp core.PipelinePlan, sc []core.Scope, err errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(plugin.SubTaskMetas(), connectionId, scopes)
}

func (plugin Jenkins) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
			"GET":    api.GetConnection,
		},
		"connections/:connectionId/scopes/:fullName": {
			"GET":   api.GetScope,
			"PATCH": api.UpdateScope,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
		"transformation_rules": {
			"POST": api.CreateTransformationRule,
			"GET":  api.GetTransformationRuleList,
		},
		"transformation_rules/:id": {
			"PATCH": api.UpdateTransformationRule,
			"GET":   api.GetTransformationRule,
		},
		"connections/:connectionId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

func (plugin Jenkins) Close(taskCtx core.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.JenkinsTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
