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
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func NewTapdApiClient(taskCtx core.TaskContext, connection *models.TapdConnection) (*helper.ApiAsyncClient, errors.Error) {
	// create synchronize api client so we can calculate api rate limit dynamically
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", connection.Username, connection.Password)))
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %v", auth),
	}

	apiClient, err := helper.NewApiClient(taskCtx.GetContext(), connection.Endpoint, headers, 0, "", taskCtx)
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) errors.Error {
		if res.StatusCode == http.StatusUnprocessableEntity {
			return errors.HttpStatus(res.StatusCode).New("authentication failed, please check your AccessToken")
		}
		return nil
	})

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: connection.RateLimitPerHour,
	}
	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		apiClient,
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}

	return asyncApiClient, nil
}
