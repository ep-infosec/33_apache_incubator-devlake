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
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jenkins/models/migrationscripts/archived"
)

type addInitTables struct{}

func (*addInitTables) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(
		"_raw_jenkins_api_jobs",
		"_raw_jenkins_api_builds",
		&archived.JenkinsJob{},
		&archived.JenkinsBuild{},
	)
	if err != nil {
		return err
	}

	err = migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.JenkinsJob{},
		&archived.JenkinsBuild{},
		&archived.JenkinsConnection{},
	)
	if err != nil {
		return err
	}

	v := config.GetConfig()

	encKey := v.GetString("ENCODE_KEY")
	endPoint := v.GetString("JENKINS_ENDPOINT")
	useName := v.GetString("JENKINS_USERNAME")
	passWord := v.GetString("JENKINS_PASSWORD")
	if encKey == "" || endPoint == "" || useName == "" || passWord == "" {
		return nil
	}
	conn := &archived.JenkinsConnection{}
	conn.Name = "init jenkins connection"
	conn.ID = 1
	conn.Endpoint = endPoint
	conn.Proxy = v.GetString("JENKINS_PROXY")
	conn.RateLimitPerHour = v.GetInt("JENKINS_API_REQUESTS_PER_HOUR")
	conn.Username = useName
	conn.Password, err = core.Encrypt(encKey, passWord)
	if err != nil {
		return err
	}
	err = db.CreateIfNotExist(conn)

	if err != nil {
		return err
	}

	return nil
}

func (*addInitTables) Version() uint64 {
	return 20220714201237
}

func (*addInitTables) Name() string {
	return "Jenkins init schemas"
}
