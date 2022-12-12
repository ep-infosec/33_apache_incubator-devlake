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
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	mockMeta := mocks.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/jira")
	err := core.RegisterPlugin("jira", mockMeta)
	assert.Nil(t, err)
	bs := &core.BlueprintScopeV200{
		Entities: []string{"TICKET"},
		Id:       "",
		Name:     "",
	}
	bpScopes := make([]*core.BlueprintScopeV200, 0)
	bpScopes = append(bpScopes, bs)

	connection := &models.JiraConnection{
		RestConnection: helper.RestConnection{
			BaseConnection: helper.BaseConnection{
				Name: "Jira",
				Model: common.Model{
					ID: 1,
				},
			},
		},
	}

	basicRes = NewMockBasicRes()
	plan := make(core.PipelinePlan, len(bpScopes))
	plan, err = makeDataSourcePipelinePlanV200(nil, plan, bpScopes, connection)
	assert.Nil(t, err)
	basicRes = NewMockBasicRes()
	scopes, err := makeScopesV200(bpScopes, connection)
	assert.Nil(t, err)

	expectPlan := core.PipelinePlan{
		core.PipelineStage{
			{
				Plugin:     "jira",
				Subtasks:   []string{},
				SkipOnFail: false,
				Options: map[string]interface{}{
					"connectionId":         uint64(1),
					"boardId":              uint64(10),
					"projectId":            uint(20),
					"name":                 "a",
					"self":                 "self",
					"transformationRuleId": uint64(1),
					"transformationRules": map[string]interface{}{
						"name":         "jira transformation rule",
						"epicKeyField": "hey,man,wasup",
					},
					"type": "type",
				},
			},
		},
	}
	assert.Equal(t, expectPlan, plan)

	expectScopes := make([]core.Scope, 0)
	jiraBoard := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: "jira:JiraBoard:1:10",
		},
		Name:        "a",
		Description: "",
		Url:         "",
		CreatedDate: nil,
	}

	expectScopes = append(expectScopes, jiraBoard)
	assert.Equal(t, expectScopes, scopes)
}

// NewMockBasicRes FIXME ...
func NewMockBasicRes() *mocks.BasicRes {
	jiraBoard := &models.JiraBoard{
		ConnectionId:         1,
		BoardId:              10,
		Name:                 "a",
		ProjectId:            20,
		Self:                 "self",
		Type:                 "type",
		TransformationRuleId: 1,
	}

	transformationRule := &models.JiraTransformationRule{
		Model: common.Model{
			ID: 1,
		},
		Name:         "jira transformation rule",
		EpicKeyField: "hey,man,wasup",
	}
	mockRes := new(mocks.BasicRes)
	mockDal := new(mocks.Dal)

	mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.JiraBoard)
		*dst = *jiraBoard
	}).Return(nil).Once()

	mockDal.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*models.JiraTransformationRule)
		*dst = *transformationRule
	}).Return(nil).Once()

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetConfig", mock.Anything).Return("")

	return mockRes
}
