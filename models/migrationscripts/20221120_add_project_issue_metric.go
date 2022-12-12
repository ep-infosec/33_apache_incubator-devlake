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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

type addProjectIssueMetric struct{}

func (u *addProjectIssueMetric) Up(baseRes core.BasicRes) errors.Error {
	db := baseRes.GetDal()
	err := migrationhelper.AutoMigrateTables(
		baseRes,
		&archived.ProjectIssueMetric{},
	)
	if err != nil {
		return err
	}
	return db.DropColumns(`issues`, "deployment_id")
}

func (*addProjectIssueMetric) Version() uint64 {
	return 20221120000001
}

func (*addProjectIssueMetric) Name() string {
	return "add project issue metric tables"
}
