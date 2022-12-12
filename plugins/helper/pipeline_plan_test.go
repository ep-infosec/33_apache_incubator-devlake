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

package helper

import (
	"testing"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/stretchr/testify/assert"
)

func TestMakePipelinePlanSubtasks(t *testing.T) {

	subtasks1, err := MakePipelinePlanSubtasks(
		[]core.SubTaskMeta{
			{
				Name:             "collectApiIssues",
				EnabledByDefault: true,
				DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
			},
			{
				Name:             "extractApiIssues",
				EnabledByDefault: true,
				DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
			},
			{
				Name:             "collectApiPullRequests",
				EnabledByDefault: true,
				DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
			},
		},
		[]string{core.DOMAIN_TYPE_TICKET},
	)
	assert.Nil(t, err)
	assert.Equal(
		t,
		subtasks1,
		[]string{"collectApiIssues", "extractApiIssues"},
	)

	subtasks2, err := MakePipelinePlanSubtasks(
		[]core.SubTaskMeta{
			{
				Name:             "collectApiRepo",
				EnabledByDefault: true,
				DomainTypes:      []string{core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CODE},
			},
			{
				Name:             "collectApiIssues",
				EnabledByDefault: true,
				DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
			},
			{
				Name:             "collectApiPullRequests",
				EnabledByDefault: true,
				DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
			},
		},
		[]string{core.DOMAIN_TYPE_TICKET},
	)
	assert.Nil(t, err)
	assert.Equal(
		t,
		subtasks2,
		[]string{"collectApiRepo", "collectApiIssues"},
	)
}
