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
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/feishu/apimodels"
	"github.com/apache/incubator-devlake/plugins/feishu/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const AUTH_ENDPOINT = "https://open.feishu.cn"
const ENDPOINT = "https://open.feishu.cn/open-apis/vc/v1"

func NewFeishuApiClient(taskCtx core.TaskContext, connection *models.FeishuConnection) (*helper.ApiAsyncClient, errors.Error) {

	authApiClient, err := helper.NewApiClient(taskCtx.GetContext(), AUTH_ENDPOINT, nil, 0, connection.Proxy, taskCtx)
	if err != nil {
		return nil, err
	}

	// request for access token
	tokenReqBody := &apimodels.ApiAccessTokenRequest{
		AppId:     connection.AppId,
		AppSecret: connection.SecretKey,
	}
	tokenRes, err := authApiClient.Post("open-apis/auth/v3/tenant_access_token/internal", nil, tokenReqBody, nil)
	if err != nil {
		return nil, err
	}
	tokenResBody := &apimodels.ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return nil, err
	}
	if tokenResBody.AppAccessToken == "" && tokenResBody.TenantAccessToken == "" {
		return nil, errors.Default.New("failed to request access token")
	}
	// real request apiClient
	apiClient, err := helper.NewApiClient(taskCtx.GetContext(), ENDPOINT, nil, 0, connection.Proxy, taskCtx)
	if err != nil {
		return nil, err
	}
	// set token
	apiClient.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", tokenResBody.TenantAccessToken),
	})

	// create async api client
	asyncApiCLient, err := helper.CreateAsyncApiClient(taskCtx, apiClient, &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: connection.RateLimitPerHour,
	})
	if err != nil {
		return nil, err
	}

	return asyncApiCLient, nil
}
