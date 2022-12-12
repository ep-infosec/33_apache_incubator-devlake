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

import styled from 'styled-components'
import { Colors } from '@blueprintjs/core'

export * from '../styled'

export const ConnectionList = styled.ul`
  margin: 0;
  padding: 12px;
  list-style: none;

  li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 4px 0;
    border-bottom: 1px solid #f0f0f0;

    .name {
      font-weight: 600;
    }

    .status {
      display: flex;
      align-items: center;

      &.online {
        color: ${Colors.GREEN3};
      }

      &.offline {
        color: ${Colors.RED3};
      }
    }
  }
`

export const Tips = styled.p`
  margin: 24px 0 0;

  span:last-child {
    color: #7497f7;
    cursor: pointer;
  }
`
