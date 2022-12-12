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

package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type apiProject struct {
	models.GitlabProject
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type req struct {
	Data []*models.GitlabProject `json:"data"`
}

// PutScope create or update gitlab project
// @Summary create or update gitlab project
// @Description Create or update gitlab project
// @Tags plugins/gitlab
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body req true "json"
// @Success 200  {object} []models.GitlabProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes [PUT]
func PutScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	var projects req
	err := errors.Convert(mapstructure.Decode(input.Body, &projects))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Gitlab project error")
	}
	keeper := make(map[int]struct{})
	now := time.Now()
	for _, project := range projects.Data {
		if _, ok := keeper[project.GitlabId]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[project.GitlabId] = struct{}{}
		}
		project.ConnectionId = connectionId
		project.CreatedDate = now
		project.UpdatedDate = &now
		err = verifyProject(project)
		if err != nil {
			return nil, err
		}
	}
	err = BasicRes.GetDal().CreateOrUpdate(projects.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving GitlabProject")
	}
	return &core.ApiResourceOutput{Body: projects.Data, Status: http.StatusOK}, nil
}

// UpdateScope patch to gitlab project
// @Summary patch to gitlab project
// @Description patch to gitlab project
// @Tags plugins/gitlab
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param projectId path int false "project ID"
// @Param scope body models.GitlabProject true "json"
// @Success 200  {object} models.GitlabProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes/{projectId} [PATCH]
func UpdateScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, projectId := extractParam(input.Params)
	if connectionId*projectId == 0 {
		return nil, errors.BadInput.New("invalid connectionId or projectId")
	}
	var project models.GitlabProject
	err := BasicRes.GetDal().First(&project, dal.Where("connection_id = ? AND gitlab_id = ?", connectionId, projectId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting GitlabProject error")
	}
	err = helper.DecodeMapStruct(input.Body, &project)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch gitlab project error")
	}
	err = verifyProject(&project)
	if err != nil {
		return nil, err
	}
	err = BasicRes.GetDal().Update(project)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving GitlabProject")
	}
	return &core.ApiResourceOutput{Body: project, Status: http.StatusOK}, nil
}

// GetScopeList get Gitlab projects
// @Summary get Gitlab projects
// @Description get Gitlab projects
// @Tags plugins/gitlab
// @Param connectionId path int false "connection ID"
// @Success 200  {object} []apiProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var projects []models.GitlabProject
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	limit, offset := helper.GetLimitOffset(input.Query, "pageSize", "page")
	err := BasicRes.GetDal().All(&projects, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	var ruleIds []uint64
	for _, proj := range projects {
		if proj.TransformationRuleId > 0 {
			ruleIds = append(ruleIds, proj.TransformationRuleId)
		}
	}
	var rules []models.GitlabTransformationRule
	if len(ruleIds) > 0 {
		err = BasicRes.GetDal().All(&rules, dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}
	names := make(map[uint64]string)
	for _, rule := range rules {
		names[rule.ID] = rule.Name
	}
	var apiProjects []apiProject
	for _, proj := range projects {
		apiProjects = append(apiProjects, apiProject{proj, names[proj.TransformationRuleId]})
	}
	return &core.ApiResourceOutput{Body: apiProjects, Status: http.StatusOK}, nil
}

// GetScope get one Gitlab project
// @Summary get one Gitlab project
// @Description get one Gitlab project
// @Tags plugins/gitlab
// @Param connectionId path int false "connection ID"
// @Param projectId path int false "project ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} apiProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes/{projectId} [GET]
func GetScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var project models.GitlabProject
	connectionId, projectId := extractParam(input.Params)
	if connectionId*projectId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := BasicRes.GetDal().First(&project, dal.Where("connection_id = ? AND gitlab_id = ?", connectionId, projectId))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.New("record not found")
	}
	if err != nil {
		return nil, err
	}
	var rule models.GitlabTransformationRule
	if project.TransformationRuleId > 0 {
		err = BasicRes.GetDal().First(&rule, dal.Where("id = ?", project.TransformationRuleId))
		if err != nil {
			return nil, err
		}
	}
	return &core.ApiResourceOutput{Body: apiProject{project, rule.Name}, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, uint64) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	projectId, _ := strconv.ParseUint(params["projectId"], 10, 64)
	return connectionId, projectId
}

func verifyProject(project *models.GitlabProject) errors.Error {
	if project.ConnectionId == 0 {
		return errors.BadInput.New("invalid connectionId")
	}
	if project.GitlabId <= 0 {
		return errors.BadInput.New("invalid projectId")
	}
	return nil
}
