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
import { Button, Card, Elevation, Icon, Intent } from '@blueprintjs/core'
import ContentLoader from '@/components/loaders/ContentLoader'
import useIntegrations from '@/hooks/useIntegrations'
import useConnectionManager from '@/hooks/useConnectionManager'
import useConnectionValidation from '@/hooks/useConnectionValidation'
import ConnectionForm from '@/pages/configure/connections/ConnectionForm'
import DeleteAction from '@/components/actions/DeleteAction'
import DeleteConfirmationMessage from '@/components/actions/DeleteConfirmationMessage'
import { NullSettings } from '@/data/NullSettings'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@/styles/configure.scss'

export default function ConfigureConnection() {
  const history = useHistory()
  const { providerId, connectionId } = useParams()

  const {
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
  // const [activeConnection, setActiveConnection] = useState(NullConnection)
  const [showConnectionSettings, setShowConnectionSettings] = useState(true)
  const [deleteId, setDeleteId] = useState(null)

  const [settings, setSettings] = useState(NullSettings)

  const {
    testConnection,
    saveConnection,
    // fetchConnection,
    activeConnection,
    name,
    endpointUrl,
    proxy,
    rateLimitPerHour = 0,
    enableGraphql,
    username,
    password,
    token,
    initialTokenStore,
    errors,
    testStatus,
    testResponse,
    allTestResponses,
    isSaving: isSavingConnection,
    isTesting: isTestingConnection,
    isFetching: isLoadingConnection,
    setName,
    setEndpointUrl,
    setProxy,
    setRateLimitPerHour,
    setEnableGraphql,
    setUsername,
    setPassword,
    setToken,
    setInitialTokenStore,
    // saveComplete: saveConnectionComplete,
    showError: showConnectionError,
    isDeleting: isDeletingConnection,
    deleteConnection,
    deleteComplete
  } = useConnectionManager(
    {
      activeProvider,
      connectionId
    },
    true
  )

  // const {
  //   saveSettings,
  //   // errors: settingsErrors,
  //   isSaving
  //   // isTesting,
  //   // showError,
  // } = useSettingsManager({
  //   activeProvider,
  //   activeConnection,
  //   settings
  // })

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

  // @todo: cleanup unused render helper
  // const renderProviderSettings = useCallback(
  //   (providerId, activeProvider) => {
  //     console.log('>>> RENDERING PROVIDER SETTINGS...')
  //     let settingsComponent = null
  //     if (activeProvider && activeProvider.settings) {
  //       settingsComponent = activeProvider.settings({
  //         activeProvider,
  //         activeConnection,
  //         isSaving,
  //         isSavingConnection,
  //         setSettings
  //       })
  //     } else {
  //       // @todo create & display "fallback/empty settings" view
  //       console.log(
  //         '>> WARNING: NO PROVIDER SETTINGS RENDERED, PROVIDER = ',
  //         activeProvider
  //       )
  //     }
  //     return settingsComponent
  //   },
  //   [activeConnection, isSaving, isSavingConnection]
  // )

  useEffect(() => {
    console.log('>>>> DETECTED PROVIDER ID = ', providerId)
    console.log('>>>> DETECTED CONNECTION ID = ', connectionId)
    if (connectionId && providerId) {
      setActiveProvider(Integrations.find((p) => p.id === providerId))
    } else {
      console.log('NO PARAMS!')
    }
  }, [connectionId, providerId, Integrations, setActiveProvider])

  useEffect(() => {}, [settings])

  useEffect(() => {
    if (deleteComplete) {
      console.log('>>> DELETE COMPLETE!')
      history.replace(`/integrations/${deleteComplete.provider?.id}`)
    }
  }, [deleteComplete, history])

  // useEffect(() => {
  //   // CONNECTION SAVED!
  // }, [saveConnectionComplete])

  return (
    <main className='main'>
      <div className='configureConnection' style={{ width: '100%' }}>
        {!isLoadingConnection && (
          <Link
            style={{
              float: 'right',
              marginLeft: '10px',
              color: '#777777'
            }}
            to={`/integrations/${activeProvider?.id}`}
          >
            <Icon icon='fast-backward' size={16} /> Connection List
          </Link>
        )}
        <div style={{ display: 'flex', justifyContent: 'flex-start' }}>
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
          {isLoadingConnection && (
            <ContentLoader
              title='Loading Connection ...'
              message='Please wait while connection settings are loaded.'
            />
          )}
          {!isLoadingConnection && (
            <div style={{ justifyContent: 'flex-start' }}>
              <div style={{ display: 'flex' }}>
                <h1 style={{ margin: 0 }}>
                  Manage{' '}
                  <strong style={{ fontWeight: 900 }}>
                    {activeProvider?.name}
                  </strong>{' '}
                  Settings
                </h1>
                {activeProvider?.multiConnection && (
                  <div style={{ paddingTop: '5px' }}>
                    <DeleteAction
                      id={deleteId}
                      connection={activeConnection}
                      text='Delete'
                      showConfirmation={() => setDeleteId(activeConnection.id)}
                      onConfirm={deleteConnection}
                      onCancel={(e) => setDeleteId(null)}
                      isDisabled={isDeletingConnection}
                      isLoading={isDeletingConnection}
                    >
                      <DeleteConfirmationMessage
                        title={`DELETE "${activeConnection.name}"`}
                      />
                    </DeleteAction>
                  </div>
                )}
              </div>
              {activeConnection && (
                <>
                  {[Providers.GITLAB, Providers.JIRA, Providers.TAPD].includes(
                    activeProvider?.id
                  ) && (
                    <h2 style={{ margin: 0 }}>
                      #{activeConnection?.id} {activeConnection.name}
                    </h2>
                  )}
                  <p className='page-description'>
                    Manage settings and options for this connection.
                  </p>
                </>
              )}
            </div>
          )}
        </div>
        {!isLoadingConnection && activeProvider && activeConnection && (
          <>
            <Card
              interactive={false}
              elevation={Elevation.ZERO}
              style={{
                backgroundColor: '#f8f8f8',
                width: '100%',
                marginBottom: '20px'
              }}
            >
              <Button
                type='button'
                icon={showConnectionSettings ? 'eye-on' : 'eye-off'}
                intent={
                  showConnectionSettings ? Intent.PRIMARY : Intent.DISABLED
                }
                style={{ margin: '2px', float: 'right' }}
                onClick={() =>
                  setShowConnectionSettings(!showConnectionSettings)
                }
                minimal
                small
              />
              {showConnectionSettings ? (
                <div className='editConnection' style={{ display: 'flex' }}>
                  <ConnectionForm
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
                    onSave={() => saveConnection()}
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
                    isSaving={isSavingConnection}
                    isTesting={isTestingConnection}
                    testStatus={testStatus}
                    testResponse={testResponse}
                    allTestResponses={allTestResponses}
                    errors={errors}
                    showError={showConnectionError}
                    authType={activeProvider?.getAuthenticationType()}
                    showLimitWarning={false}
                    sourceLimits={ProviderConnectionLimits}
                    labels={ProviderFormLabels[activeProvider?.id]}
                    placeholders={ProviderFormPlaceholders[activeProvider?.id]}
                    tooltips={ProviderFormTooltips[activeProvider?.id]}
                  />
                </div>
              ) : (
                <>
                  <h2 style={{ margin: 0 }}>Configure Connection</h2>
                  <p className='description' style={{ margin: 0 }}>
                    ( Click the <strong>Visibility</strong> icon to your right
                    to edit connection )
                  </p>
                </>
              )}
              {/* {validationErrors.length > 0 && (
                      <FormValidationErrors errors={validationErrors} />
                    )} */}
            </Card>
            {/* <div style={{ marginTop: '30px' }}>
                    {renderProviderSettings(providerId, activeProvider)}
                  </div> */}
            {/* <div className='form-actions-block' style={{ display: 'flex', marginTop: '60px', justifyContent: 'space-between' }}>
                    <div />
                    <div>
                      <Button icon='remove' text='Cancel' onClick={cancel} disabled={isSaving} />
                      <Button
                        icon='cloud-upload'
                        intent={Intent.PRIMARY}
                        text='Save Settings'
                        loading={isSaving}
                        disabled={isSaving || providerId === Providers.JENKINS}
                        onClick={saveSettings}
                        style={{ marginLeft: '10px' }}
                      />
                    </div>
                  </div> */}
          </>
        )}
      </div>
    </main>
  )
}
