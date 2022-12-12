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
import classNames from 'classnames'
import React, { useEffect, useState, useRef, useCallback } from 'react'
import {
  Button,
  Classes,
  Colors,
  Dialog,
  DialogStep,
  MultistepDialog,
  MultistepDialogNavPosition,
  Elevation,
  FormGroup,
  Icon,
  Intent,
  Label,
  MenuItem,
  Position
} from '@blueprintjs/core'

import { NullBlueprint } from '@/data/NullBlueprint'
import DataScopes from '@/components/blueprints/create-workflow/DataScopes'
import DataTransformations from '@/components/blueprints/create-workflow/DataTransformations'

const Modes = {
  CREATE: 'create',
  EDIT: 'edit'
}

const DialogPanel = (props) => {
  const { children } = props

  return (
    <div
      className={classNames(Classes.DIALOG_BODY)}
      style={{ minHeight: '300px' }}
    >
      {children}
    </div>
  )
}

const BlueprintDataScopesDialog = (props) => {
  const {
    isOpen = false,
    title = 'Change Data Scope',
    blueprintConnections = [],
    blueprint = NullBlueprint,
    provider,
    activeTransformation,
    configuredConnection,
    configuredScopeEntity,
    scopeConnection,
    jiraBoards = [],
    issueTypesList = [],
    fieldsList = [],
    setBoardSearch = () => {},
    gitlabProjects = [],
    fetchGitlabProjects = () => [],
    fetchJenkinsJobs = () => [],
    jenkinsJobs = [],
    dataDomainsGroup = {},
    scopeEntitiesGroup = {},
    mode = Modes.EDIT,
    canOutsideClickClose = false,
    showCloseButtonInFooter = true,
    resetOnClose = true,
    isCloseButtonShown = true,
    initialStepIndex = 0,
    navPosition = 'top',
    usePortal = true,
    hasTitle = true,
    activeStep = null,
    onStepChange = () => {},
    onOpening = () => {},
    onClose = () => {},
    onCancel = () => {},
    onSave = () => {},
    setDataDomainsGroup = () => {},
    setScopeEntitiesGroup = () => {},
    hasConfiguredEntityTransformationChanged = () => false,
    changeConfiguredEntityTransformation = () => {},
    setConfiguredScopeEntity = () => {},
    fieldHasError = () => {},
    getFieldError = () => {},
    isSaving = false,
    isValid = true,
    isTesting = false,
    isFetchingJIRA = false,
    jiraProxyError,
    isFetchingGitlab,
    gitlabProxyError,
    isFetchingJenkins,
    jenkinsProxyError,
    errors = [],
    content = null,
    backButtonProps = {
      // disabled:
      intent: Intent.PRIMARY,
      text: 'Previous Step',
      outlined: true,
      loading:
        isFetchingJIRA || isFetchingGitlab || isFetchingJenkins || isSaving
    },
    nextButtonProps = {
      disabled: !isValid,
      intent: Intent.PRIMARY,
      text: 'Next Step',
      outlined: true,
      loading:
        isFetchingJIRA || isFetchingGitlab || isFetchingJenkins || isSaving
    },
    finalButtonProps = {
      disabled: !isValid,
      intent: Intent.PRIMARY,
      onClick: onSave,
      text: 'Save Changes',
      loading:
        isFetchingJIRA || isFetchingGitlab || isFetchingJenkins || isSaving
    },
    closeButtonProps = {
      // disabled:
      intent: Intent.PRIMARY,
      text: 'Cancel',
      outlined: true,
      loading:
        isFetchingJIRA || isFetchingGitlab || isFetchingJenkins || isSaving
    }
  } = props

  return (
    <>
      <MultistepDialog
        disabled
        isOpen={isOpen}
        className='blueprint-data-scopes-dialog'
        // icon='info-sign'
        navigationPosition={Position.BOTTOM}
        closeButtonProps={closeButtonProps}
        backButtonProps={backButtonProps}
        nextButtonProps={nextButtonProps}
        finalButtonProps={finalButtonProps}
        title={title}
        hasTitle={hasTitle}
        initialStepIndex={initialStepIndex}
        showCloseButtonInFooter={showCloseButtonInFooter}
        isCloseButtonShown={isCloseButtonShown}
        canOutsideClickClose={canOutsideClickClose}
        resetOnClose={resetOnClose}
        onOpening={onOpening}
        onClose={onClose}
        onClosed={() => {}}
        onChange={onStepChange}
      >
        <DialogStep
          id='scopes'
          panel={
            <DialogPanel>
              <DataScopes
                provider={provider}
                activeStep={activeStep}
                blueprintConnections={blueprintConnections}
                scopeEntitiesGroup={scopeEntitiesGroup}
                jiraBoards={jiraBoards}
                setBoardSearch={setBoardSearch}
                fetchGitlabProjects={fetchGitlabProjects}
                gitlabProjects={gitlabProjects}
                isFetchingGitlab={isFetchingGitlab}
                gitlabProxyError={gitlabProxyError}
                fetchJenkinsJobs={fetchJenkinsJobs}
                jenkinsJobs={jenkinsJobs}
                isFetchingJenkins={isFetchingJenkins}
                jenkinsProxyError={jenkinsProxyError}
                dataDomainsGroup={dataDomainsGroup}
                configuredConnection={configuredConnection}
                setDataDomainsGroup={setDataDomainsGroup}
                setScopeEntitiesGroup={setScopeEntitiesGroup}
                isSaving={isSaving}
                isLoading={isFetchingJIRA}
                validationErrors={[]}
                enableConnectionTabs={false}
                elevation={Elevation.ZERO}
                cardStyle={{ padding: 0, backgroundColor: 'transparent' }}
              />
            </DialogPanel>
          }
          title='Data Scopes'
        />
        <DialogStep
          id='transformations'
          panel={
            <DialogPanel>
              <DataTransformations
                provider={provider}
                blueprint={blueprint}
                activeTransformation={activeTransformation}
                blueprintConnections={blueprintConnections}
                dataDomainsGroup={dataDomainsGroup}
                scopeEntitiesGroup={scopeEntitiesGroup}
                issueTypes={issueTypesList}
                fields={fieldsList}
                configuredConnection={configuredConnection}
                configuredScopeEntity={configuredScopeEntity}
                setConfiguredScopeEntity={setConfiguredScopeEntity}
                isSaving={isSaving}
                hasConfiguredEntityTransformationChanged={
                  hasConfiguredEntityTransformationChanged
                }
                changeConfiguredEntityTransformation={
                  changeConfiguredEntityTransformation
                }
                // onSave={handleTransformationSave}
                // onCancel={handleTransformationCancel}
                // onClear={handleTransformationClear}
                jiraProxyError={jiraProxyError}
                isFetchingJIRA={isFetchingJIRA}
                fieldHasError={fieldHasError}
                getFieldError={getFieldError}
                enableConnectionTabs={false}
                enableNoticeAlert={false}
                useDropdownSelector={true}
                enableGoBack={false}
                elevation={Elevation.ZERO}
                cardStyle={{ padding: 0, backgroundColor: 'transparent' }}
              />
            </DialogPanel>
          }
          title='Transformations'
        />
      </MultistepDialog>
    </>
  )
}

export default BlueprintDataScopesDialog
