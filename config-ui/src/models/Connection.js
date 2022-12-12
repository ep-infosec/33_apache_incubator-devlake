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

/**
 * @typedef {object} Connection
 * @property {number?} id
 * @property {number?} connectionId
 * @property {string?} name
 * @property {string?} endpoint
 * @property {string?} endpointUrl
 * @property {string?} proxy
 * @property {string?} token
 * @property {object?} initialTokenStore
 * @property {string?} username
 * @property {string?} password
 * @property {number?} rateLimitPerHour
 * @property {Date?} createdAt
 * @property {Date?} updatedAt
 * @property {'plain'|'token'?} authentication
 * @property {string|object?} plugin
 * @property {string|object?} provider
 * @property {string?} providerId
 * @property {Array<DataDomain>} dataDomains
 * @property {boolean} multiConnection
 * @property {number|string?} status
 * @property {Array<object>?} transformations
 * @property {Array<string>?} errors
 */
class Connection {
  constructor(data = {}) {
    this.id = parseInt(data?.id, 10) || null
    this.connectionId = parseInt(data?.connectionId, 10) || this.id
    this.name = data?.name || ''
    this.endpoint = data?.endpoint || ''
    // @todo: will be refactored and replaced by $this.endpoint!
    this.endpointUrl = data?.endpointUrl || this.endpoint
    this.proxy = data?.proxy || ''
    this.token = data?.token || ''
    this.initialTokenStore = data?.initialTokensTore || { 0: '', 1: '', 2: '' }
    this.username = data?.username || ''
    this.password = data?.password || ''
    this.rateLimitPerHour = data?.rateLimitPerHour || 0
    this.enableGraphql = data?.enableGraphql || false
    this.createdAt = data?.createdAt || null
    this.updatedAt = data?.updatedAt || null

    this.authentication = data?.authentication || 'plain'
    this.plugin = data?.plugin || null
    // @todo: will be replaced out by $this.plugin
    this.provider = data?.provider || null
    this.providerId = data?.providerId || null
    this.dataDomains = data?.dataDomains || []
    this.multiConnection = data?.multiConnection || true
    this.status = data?.status || null

    this.transformations = data?.transformations || []

    this.errors = data?.errors || []

    this.determineAuthentication()
  }

  get(property) {
    return this[property]
  }

  set(property, value) {
    this[property] = value
    return this.property
  }

  determineAuthentication() {
    if (this.token !== null && this.token !== '') {
      this.authentication = 'token'
    } else {
      this.authentication = 'plain'
    }
  }
}

export default Connection
