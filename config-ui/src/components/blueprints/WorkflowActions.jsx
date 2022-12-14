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
import React, { Fragment, useEffect, useState, useCallback } from 'react'
import { ENVIRONMENT } from '@/config/environment'
import { CSSTransition } from 'react-transition-group'
import {
  PopoverInteractionKind,
  Button,
  Icon,
  Intent,
  Elevation,
  Popover,
  Card,
  Colors
} from '@blueprintjs/core'
import FormValidationErrors from '@/components/messages/FormValidationErrors'

const WorkflowActions = (props) => {
  const {
    activeStep,
    blueprintSteps = [],
    setShowBlueprintInspector = () => {},
    validationErrors = [],
    onPrev = () => {},
    onNext = () => {},
    onSave = () => {},
    onSaveAndRun = () => {},
    isLoading = false,
    isValid = true,
    canGoNext = true,
    canGoPrev = true
  } = props

  return (
    <div className='workflow-actions'>
      <Button
        loading={isLoading}
        disabled={activeStep?.id === 1 || isLoading}
        intent={Intent.PRIMARY}
        text='Previous Step'
        onClick={onPrev}
      />

      {activeStep?.id === blueprintSteps.length ? (
        <div style={{ marginLeft: 'auto' }}>
          <Button
            loading={isLoading}
            disabled={isLoading}
            intent={Intent.PRIMARY}
            text='Save Blueprint'
            // disabled
            onClick={onSave}
          />
          <Button
            loading={isLoading}
            disabled={isLoading}
            intent={Intent.DANGER}
            text='Save and Run Now'
            style={{ marginLeft: '5px' }}
            // disabled
            onClick={onSaveAndRun}
          />
        </div>
      ) : (
        <div style={{ display: 'flex', marginLeft: 'auto' }}>
          {ENVIRONMENT !== 'production' && (
            <Button
              loading={isLoading}
              intent={Intent.PRIMARY}
              icon='code'
              text='Inspect'
              onClick={() => setShowBlueprintInspector(true)}
              style={{ marginRight: '8px' }}
              minimal
              small
            />
          )}
          <Button
            loading={isLoading}
            disabled={isLoading || !canGoNext || !isValid}
            intent={Intent.PRIMARY}
            text='Next Step'
            onClick={onNext}
            rightIcon={
              validationErrors.length > 0 ? (
                <Popover
                  interactionKind={PopoverInteractionKind.HOVER_TARGET_ONLY}
                  defaultIsOpen={true}
                  enforceFocus={false}
                >
                  <Icon
                    icon='warning-sign'
                    size={12}
                    color={Colors.ORANGE5}
                    style={{ outline: 'none', margin: '0 3px 2px 3px' }}
                  />
                  <div style={{ padding: '5px' }}>
                    <FormValidationErrors errors={validationErrors} />
                  </div>
                </Popover>
              ) : null
            }
          />
        </div>
      )}
    </div>
  )
}

export default WorkflowActions
