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

package e2e

import (
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/dora/impl"
	"github.com/apache/incubator-devlake/plugins/dora/tasks"
)

func TestEnrichEnvDataFlow(t *testing.T) {
	var plugin impl.Dora
	dataflowTester := e2ehelper.NewDataFlowTester(t, "dora", plugin)

	taskData := &tasks.DoraTaskData{
		Options: &tasks.DoraOptions{
			ProjectName: "project1",
			TransformationRules: tasks.TransformationRules{
				ProductionPattern: "(?i)deploy",
				StagingPattern:    "(?i)stag",
				TestingPattern:    "(?i)test",
			},
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoTabler("./raw_tables/project_mapping.csv", &crossdomain.ProjectMapping{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/cicd_tasks_for_enrich_env.csv", &devops.CICDTask{})

	// verify enrich with repoId
	dataflowTester.Subtask(tasks.EnrichTaskEnvMeta, taskData)
	dataflowTester.VerifyTable(
		&devops.CICDTask{},
		"./snapshot_tables/enrich_cicd_tasks.csv",
		[]string{
			"id",
			"name",
			"environment",
			"type",
			"cicd_scope_id",
		},
	)
}
