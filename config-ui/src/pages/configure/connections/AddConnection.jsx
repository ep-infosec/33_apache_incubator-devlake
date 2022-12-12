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
import React, { useEffect, useState } from 'react'
import { Link, useHistory, useParams } from 'react-router-dom'
import { Icon } from '@blueprintjs/core'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'
// import { integrationsData } from '@/data/integrations'
// import {
//   ProviderConnectionLimits,
//   ProviderFormLabels,
//   ProviderFormPlaceholders,
//   ProviderLabels,
//   Providers
// } from '@/data/Providers'

import useIntegrations from '@/hooks/useIntegrations'
import useConnectionManager from '@/hooks/useConnectionManager'
import useConnectionValidation from '@/hooks/useConnectionValidation'

import '@/styles/integration.scss'
import '@/styles/connections.scss'

export default function AddConnection() {
  const history = useHistory()
  const { providerId } = useParams()

  const {
    registry,
    plugins: Plugins,
    integrations: Integrations,
    activeProvider,
    Providers,
    ProviderFormLabels,
    ProviderFormPlaceholders,
    ProviderFormTooltips,
    ProviderConnectionLimits,
    setActiveProvider
  } = useIntegrations()

  // @todo: Replace with Integrations Hook
  // const [activeProvider, setActiveProvider] = useState(
  //   integrationsData.find((p) => p.id === providerId)
  // )

  const {
    testConnection,
    saveConnection,
    errors,
    isSaving,
    isTesting,
    showError,
    testStatus,
    testResponse,
    allTestResponses,
    name,
    endpointUrl,
    proxy,
    rateLimitPerHour,
    enableGraphql,
    token,
    initialTokenStore,
    username,
    password,
    setName,
    setEndpointUrl,
    setProxy,
    setRateLimitPerHour,
    setEnableGraphql,
    setUsername,
    setPassword,
    setToken,
    setInitialTokenStore,
    fetchAllConnections,
    connectionLimitReached
    // Providers
  } = useConnectionManager({
    activeProvider
  })

  const {
    validate,
    errors: validationErrors,
    isValid: isValidForm
  } = useConnectionValidation({
    activeProvider,
    name,
    endpointUrl,
    proxy,
    rateLimitPerHour,
    token,
    username,
    password
  })

  const cancel = () => {
    history.push(`/integrations/${activeProvider?.id}`)
  }

  // const resetForm = () => {
  //   setName(null)
  //   setEndpointUrl(null)
  //   setToken(null)
  //   setUsername(null)
  //   setPassword(null)
  // }

  useEffect(() => {
    // @todo: Cleanup Restricted Provider Names (Legacy Feature)
    // Selected Provider
    // if (activeProvider?.id) {
    //   fetchAllConnections()
    //   switch (activeProvider?.id) {
    //     case Providers.JENKINS:
    //       // setName(ProviderLabels.JENKINS)
    //       break
    //     case Providers.GITHUB:
    //     case Providers.GITLAB:
    //     case Providers.JIRA:
    //     case Providers.TAPD:
    //     default:
    //       setName('')
    //       break
    //   }
    // }
  }, [activeProvider?.id, fetchAllConnections, setName])

  useEffect(() => {
    console.log('>>>> DETECTED PROVIDER = ', providerId)
    setActiveProvider(Integrations.find((p) => p.id === providerId))
  }, [providerId, setActiveProvider, Integrations])

  return (
    <main className='main'>
      <div style={{ width: '100%' }}>
        <Link
          style={{ float: 'right', marginLeft: '10px', color: '#777777' }}
          to={`/integrations/${activeProvider?.id}`}
        >
          <Icon icon='undo' size={16} /> Go Back
        </Link>
        <div style={{ display: 'flex' }}>
          <div>
            <span style={{ marginRight: '10px' }}>
              <img
                className='providerIconSvg'
                src={'/' + activeProvider?.icon}
                width={40}
                height={40}
                style={{ width: '40px', height: '40px' }}
              />
            </span>
          </div>
          <div>
            <h1 style={{ margin: 0 }}>{activeProvider?.name} Add Connection</h1>
            <p className='page-description'>
              Create a new connection for this provider.
            </p>
          </div>
        </div>
        <div className='addConnection' style={{ display: 'flex' }}>
          <ConnectionForm
            isLocked={connectionLimitReached}
            isValid={isValidForm}
            validationErrors={validationErrors}
            activeProvider={activeProvider}
            name={name}
            endpointUrl={endpointUrl}
            proxy={proxy}
            rateLimitPerHour={rateLimitPerHour}
            enableGraphql={enableGraphql}
            token={token}
            initialTokenStore={initialTokenStore}
            username={username}
            password={password}
            onSave={() => saveConnection({})}
            onTest={testConnection}
            onCancel={cancel}
            onValidate={validate}
            onNameChange={setName}
            onEndpointChange={setEndpointUrl}
            onProxyChange={setProxy}
            onRateLimitChange={setRateLimitPerHour}
            onEnableGraphqlChange={setEnableGraphql}
            onTokenChange={setToken}
            onUsernameChange={setUsername}
            onPasswordChange={setPassword}
            isSaving={isSaving}
            isTesting={isTesting}
            testStatus={testStatus}
            testResponse={testResponse}
            allTestResponses={allTestResponses}
            errors={errors}
            showError={showError}
            authType={activeProvider?.getAuthenticationType()}
            sourceLimits={ProviderConnectionLimits}
            labels={ProviderFormLabels[activeProvider?.id]}
            placeholders={ProviderFormPlaceholders[activeProvider?.id]}
            tooltips={ProviderFormTooltips[activeProvider?.id]}
          />
        </div>
      </div>
    </main>
  )
}
