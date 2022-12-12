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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/github/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func CreateApiClient(taskCtx core.TaskContext, connection *models.GithubConnection) (*helper.ApiAsyncClient, errors.Error) {
	// load configuration
	tokens := strings.Split(connection.Token, ",")
	tokenIndex := 0
	// create synchronize api client so we can calculate api rate limit dynamically
	apiClient, err := helper.NewApiClient(taskCtx.GetContext(), connection.Endpoint, nil, 0, connection.Proxy, taskCtx)
	if err != nil {
		return nil, err
	}
	// Rotates token on each request.
	apiClient.SetBeforeFunction(func(req *http.Request) errors.Error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokens[tokenIndex]))
		// Set next token index
		tokenIndex = (tokenIndex + 1) % len(tokens)
		return nil
	})

	// create rate limit calculator
	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: connection.RateLimitPerHour,
		Method:               http.MethodGet,
		DynamicRateLimit: func(res *http.Response) (int, time.Duration, errors.Error) {
			/* calculate by number of remaining requests
			remaining, err := strconv.Atoi(res.Header.Get("X-RateLimit-Remaining"))
			if err != nil {
				return 0,0, errors.Default.New("failed to parse X-RateLimit-Remaining header: %w", err)
			}
			reset, err := strconv.Atoi(res.Header.Get("X-RateLimit-Reset"))
			if err != nil {
				return 0, 0, errors.Default.New("failed to parse X-RateLimit-Reset header: %w", err)
			}
			date, err := http.ParseTime(res.Header.Get("Date"))
			if err != nil {
				return 0, 0, errors.Default.New("failed to parse Date header: %w", err)
			}
			return remaining * len(tokens), time.Unix(int64(reset), 0).Sub(date), nil
			*/
			rateLimit, err := strconv.Atoi(res.Header.Get("X-RateLimit-Limit"))
			if err != nil {
				return 0, 0, errors.Default.Wrap(err, "failed to parse X-RateLimit-Limit header")
			}
			// even though different token could have different rate limit, but it is hard to support it
			// so, we calculate the rate limit of a single token, and presume all tokens are the same, to
			// simplify the algorithm for now
			// TODO: consider different token has different rate-limit
			return rateLimit * len(tokens), 1 * time.Hour, nil
		},
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
