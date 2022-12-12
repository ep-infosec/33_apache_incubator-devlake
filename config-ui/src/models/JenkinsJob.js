/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import Entity from '@/models/Entity'

/**
 * @typedef {object} JenkinsJob
 * @property {number?} id
 * @property {number?} key
 * @property {number?} projectId
 * @property {string|number?} name
 * @property {string|number?} value
 * @property {string|number?} title
 * @property {boolean?} useApi
 * @property {project|board|job?} variant
 * @property {string?} providerId
 */
class JenkinsJob extends Entity {
  constructor(data = {}) {
    super(data)
    this.id = data?.id || null
    this.key = data?.key || this.id || null
    this.name = data?.name || null
    this.value = data?.value || this.name || this.id || null
    this.title = data?.title || this.name || this.id || null
    this.jobPath = data?.jobPath || ''

    this.useApi = data?.useApi || true
    this.variant = data?.variant || 'job'
    this.providerId = 'jenkins'
  }

  getConfiguredEntityId() {
    return this.id
  }

  getTransformationScopeOptions() {
    return {
      jobName: this.value,
      jobPath: this.jobPath
    }
  }
}

export default JenkinsJob
