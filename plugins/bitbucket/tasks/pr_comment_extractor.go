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
	"encoding/json"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiPrCommentsMeta = core.SubTaskMeta{
	Name:             "extractApiPullRequestsComments",
	EntryPoint:       ExtractApiPullRequestsComments,
	EnabledByDefault: true,
	Required:         true,
	Description:      "Extract raw pull requests comments data into tool layer table BitbucketPrComments",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE_REVIEW},
}

type BitbucketPrCommentsResponse struct {
	BitbucketId int        `json:"id"`
	CreatedOn   time.Time  `json:"created_on"`
	UpdatedOn   *time.Time `json:"updated_on"`
	Content     struct {
		Type string `json:"type"`
		Raw  string `json:"raw"`
	} `json:"content"`
	User    *BitbucketAccountResponse `json:"user"`
	Deleted bool                      `json:"deleted"`
	Type    string                    `json:"type"`
	Links   struct {
		Self struct {
			Href string
		} `json:"self"`
		Html struct {
			Href string
		} `json:"html"`
	} `json:"links"`
	PullRequest struct {
		Type  string `json:"type"`
		Id    int    `json:"id"`
		Title string `json:"title"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Html struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
	}
}

func ExtractApiPullRequestsComments(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_COMMENTS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			prComment := &BitbucketPrCommentsResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, prComment))
			if err != nil {
				return nil, err
			}

			toolprComment, err := convertPullRequestComment(prComment)
			toolprComment.ConnectionId = data.Options.ConnectionId
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 2)

			if prComment.User != nil {
				bitbucketUser, err := convertAccount(prComment.User, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				toolprComment.AuthorId = bitbucketUser.AccountId
				toolprComment.AuthorName = bitbucketUser.UserName
				results = append(results, bitbucketUser)
			}
			results = append(results, toolprComment)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertPullRequestComment(prComment *BitbucketPrCommentsResponse) (*models.BitbucketPrComment, errors.Error) {
	bitbucketPrComment := &models.BitbucketPrComment{
		BitbucketId:        prComment.BitbucketId,
		AuthorId:           prComment.User.AccountId,
		PullRequestId:      prComment.PullRequest.Id,
		AuthorName:         prComment.User.DisplayName,
		BitbucketCreatedAt: prComment.CreatedOn,
		BitbucketUpdatedAt: prComment.UpdatedOn,
		Type:               prComment.Type,
		Body:               prComment.Content.Raw,
	}
	return bitbucketPrComment, nil
}
