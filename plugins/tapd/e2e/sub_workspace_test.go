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
	"testing"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

func TestTapdWorkspaceDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			CompanyId:    99,
			WorkspaceId:  991,
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_sub_workspaces.csv",
		"_raw_tapd_api_sub_workspaces")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdSubWorkspace{})
	dataflowTester.Subtask(tasks.ExtractSubWorkspaceMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdSubWorkspace{},
		"./snapshot_tables/_tool_tapd_sub_workspaces.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"id",
			"name",
			"pretty_name",
			"category",
			"status",
			"description",
			"begin_date",
			"end_date",
			"external_on",
			"parent_id",
			"creator",
		),
	)

	dataflowTester.FlushTabler(&ticket.Board{})
	dataflowTester.Subtask(tasks.ConvertSubWorkspaceMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Board{},
		"./snapshot_tables/boards.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"name",
			"description",
			"url",
			"created_date",
		),
	)

}
