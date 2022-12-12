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
import React, { useState } from 'react'
import { Card, Elevation } from '@blueprintjs/core'
import StageTaskName from '@/components/pipelines/StageTaskName'
import StageTaskIndicator from '@/components/pipelines/StageTaskIndicator'
import StageTaskCaption from '@/components/pipelines/StageTaskCaption'

import { ReactComponent as SyncIcon } from '@/images/icons/sync.svg'

const StageTask = (props) => {
  const { task, rerunTask } = props

  const [taskModuleOpened, setTaskModuleOpened] = useState(null)

  const handleRerunTask = (e) => {
    e.stopPropagation()
    rerunTask(task.id)
  }

  const generateStageTaskCssClasses = () => {
    return `pipeline-task-module task-${task.status
      .split('_')[1]
      .toLowerCase()} ${task.ID === taskModuleOpened?.ID ? 'active' : ''}`
  }

  const determineCardElevation = (status, isElevated = false) => {
    let elevation = Elevation.ZERO
    if (status === 'TASK_RUNNING' && isElevated) {
      elevation = Elevation.ZERO
    } else if (status === 'TASK_RUNNING' && !isElevated) {
      elevation = Elevation.ZERO
    } else if (isElevated) {
      elevation = Elevation.ZERO
    } else {
      elevation = Elevation.ZERO
    }
    return elevation
  }

  return (
    <>
      <Card
        elevation={determineCardElevation(
          task.status,
          taskModuleOpened !== null
        )}
        className={generateStageTaskCssClasses()}
        onClick={() => setTaskModuleOpened(task)}
        style={{}}
      >
        <StageTaskIndicator task={task} />
        <div
          className='task-module-name'
          style={{
            flex: 1
          }}
        >
          <div style={{ padding: '4px 2px 4px 0', width: '100%' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <StageTaskName
                task={task}
                showDetails={taskModuleOpened}
                onClose={() => setTaskModuleOpened(null)}
              />
              <span onClick={handleRerunTask}>
                <SyncIcon />
              </span>
            </div>
            <StageTaskCaption task={task} options={task.options} />
          </div>
        </div>
        {/* <div
          className='task-module-progress'
          style={{
            fontWeight: task.status === 'TASK_COMPLETED' ? 700 : 300
          }}
        >
          {Number(task.status === 'TASK_COMPLETED' ? 100 : (task.progress / 1) * 100).toFixed(0)}%
        </div> */}
      </Card>
    </>
  )
}

export default StageTask
