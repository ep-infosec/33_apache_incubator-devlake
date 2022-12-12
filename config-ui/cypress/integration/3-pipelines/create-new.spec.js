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

/// <reference types="cypress" />

context.skip('Create New Pipelines Interface', () => {
  beforeEach(() => {
    cy.visit('/pipelines/create')
  })

  it('provides access to creating a new pipeline', () => {
    cy.get('ul.bp3-breadcrumbs')
      .find('a.bp3-breadcrumb-current')
      .contains(/create pipeline run/i)
      .should('be.visible')
      .should('have.attr', 'href', '/pipelines/create')

    cy.get('.headlineContainer')
      .find('h1')
      .contains(/create pipeline run/i)
      .should('be.visible')
  })

  it('has form control for pipeline name', () => {
    cy.get('h2')
    .contains(/pipeline name/i)
    .should('be.visible')

    cy.get('input#pipeline-name')
    .should('be.visible')
  })

  it('has plugin support for gitlab data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-gitlab')
      .should('be.visible')
  })

  it('has plugin support for github data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-github')
      .should('be.visible')
  })

  it('has plugin support for jenkins data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-jenkins')
      .should('be.visible')
  })

  it('has plugin support for jira data provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-jira')
      .should('be.visible')
  })

  it('has plugin support for refdiff plugin provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-refdiff')
      .should('be.visible')
  })

  it('has plugin support for gitextractor plugin provider', () => {
    cy.get('.data-providers')
      .find('.data-provider-row.data-provider-gitextractor')
      .should('be.visible')
  })

  it('has form button control for running pipeline', () => {
    cy.get('.btn-run-pipeline')
      .should('be.visible')
  })

  it('has form button control for resetting pipeline', () => {
    cy.get('.btn-reset-pipeline')
      .should('be.visible')
  })

  it('has form button control for viewing all pipelines (manage)', () => {
    cy.get('.btn-view-jobs')
      .should('be.visible')
  })

  it('supports advanced-mode user interface options', () => {
    cy.get('.advanced-mode-toggleswitch')
      .should('be.visible')
      .find('.bp3-control-indicator')
      .click()

    cy.get('h2')
      .contains(/pipeline name \(advanced\)/i)
      .should('be.visible')
  })


})

context.skip('RUN / Trigger New Pipelines', () => {
  beforeEach(() => {
    cy.visit('/pipelines/create')
  })

  it('can run a jenkins pipeline', () => {
    cy.fixture('new-jenkins-pipeline').then((newJenkinsPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newJenkinsPipelineJSON.id}`, { body: newJenkinsPipelineJSON }).as('JenkinsPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newJenkinsPipelineJSON }).as('runJenkinsPipeline')
      cy.fixture('new-jenkins-pipeline-tasks').then((newJenkinsPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newJenkinsPipelineJSON.id}/tasks`, { body: newJenkinsPipelineTasksJSON }).as('JenkinsPipelineTasks')
      })
    })

    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT JENKINS ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-jenkins')
      .should('be.visible')
      .click()

    cy.get('button#btn-run-pipeline').click()
    cy.wait('@JenkinsPipeline')
    cy.wait('@JenkinsPipelineTasks')
    cy.wait('@runJenkinsPipeline').then(({ response }) => {
      const NewJenkinsRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewJenkinsRun.id}`)
    })
  })

  it('can run a gitlab pipeline', () => {
    cy.fixture('new-gitlab-pipeline').then((newGitlabPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newGitlabPipelineJSON.id}`, { body: newGitlabPipelineJSON }).as('GitlabPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newGitlabPipelineJSON }).as('runGitlabPipeline')
      cy.fixture('new-gitlab-pipeline-tasks').then((newGitlabPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newGitlabPipelineJSON.id}/tasks`, { body: newGitlabPipelineTasksJSON }).as('GitlabPipelineTasks')
      })
    })
    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT GITLAB ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-gitlab')
      .should('be.visible')
      .click()

    cy.get('.input-project-id').find('input').type('278964{enter}')

    cy.get('button#btn-run-pipeline').click()
    cy.wait('@GitlabPipeline')
    cy.wait('@GitlabPipelineTasks')
    cy.wait('@runGitlabPipeline').then(({ response }) => {
      const NewGitlabRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewGitlabRun.id}`)
    })
  })

  it('can run a github pipeline', () => {
    cy.fixture('new-github-pipeline').then((newGithubPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newGithubPipelineJSON.id}`, { body: newGithubPipelineJSON }).as('GithubPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newGithubPipelineJSON }).as('runGithubPipeline')
      cy.fixture('new-github-pipeline-tasks').then((newGithubPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newGithubPipelineJSON.id}/tasks`, { body: newGithubPipelineTasksJSON }).as('GithubPipelineTasks')
      })
    })
    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT GITHUB ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-github')
      .should('be.visible')
      .click()
      .trigger('mouseleave')

    cy.get('input#owner').click().type('merico-dev', {force: true})
    cy.get('input#repository-name').type('lake')

    cy.get('button#btn-run-pipeline').click()
    cy.wait('@GithubPipeline')
    cy.wait('@GithubPipelineTasks')
    cy.wait('@runGithubPipeline').then(({ response }) => {
      const NewGithubRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewGithubRun.id}`)
    })
  })

  it('can run a gitextractor pipeline', () => {
    cy.fixture('new-gitextractor-pipeline').then((newGitExtractorPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newGitExtractorPipelineJSON.id}`, { body: newGitExtractorPipelineJSON }).as('GitExtractorPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newGitExtractorPipelineJSON }).as('runGitExtractorPipeline')
      cy.fixture('new-github-pipeline-tasks').then((newGitExtractorPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newGitExtractorPipelineJSON.id}/tasks`, { body: newGitExtractorPipelineTasksJSON }).as('GitExtractorPipelineTasks')
      })
    })
    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT GITEXTRACTOR ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-gitextractor')
      .should('be.visible')
      .click()
      .trigger('mouseleave')

    cy.get('input#gitextractor-url').click().type('https://github.com/apache/incubator-devlake.git')
    cy.get('input#gitextractor-repo-id').type('github:GithubRepo:384111310')

    cy.get('button#btn-run-pipeline').click()
    cy.wait('@GitExtractorPipeline')
    cy.wait('@GitExtractorPipelineTasks')
    cy.wait('@runGitExtractorPipeline').then(({ response }) => {
      const NewGitExtractorRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewGitExtractorRun.id}`)
    })
  })


  it('can run a jira pipeline', () => {
    cy.fixture('new-jira-pipeline').then((newJiraPipelineJSON) => {
      cy.intercept('GET', `/api/pipelines/${newJiraPipelineJSON.id}`, { body: newJiraPipelineJSON }).as('JiraPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newJiraPipelineJSON }).as('runJiraPipeline')
      cy.fixture('new-jira-pipeline-tasks').then((newJiraPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newJiraPipelineJSON.id}/tasks`, { body: newJiraPipelineTasksJSON }).as('JiraPipelineTasks')
      })
    })
    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT JIRA ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-jira')
      .should('be.visible')
      .click()

    cy.get('button.btn-connection-id-selector').click()
    cy.wait(500)
    cy.get('.bp3-select-popover.source-id-popover')
      .find('ul.bp3-menu li')
      .first()
      .click()
    cy.wait(500)
    cy.get('.input-board-id').find('input').type('1{enter}')


    cy.get('button#btn-run-pipeline').click()
    cy.wait('@JiraPipeline')
    cy.wait('@JiraPipelineTasks')
    cy.wait('@runJiraPipeline').then(({ response }) => {
      const NewJiraRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewJiraRun.id}`)
    })
  })

  it('can run a refdiff pipeline', () => {
    cy.fixture('new-refdiff-pipeline').then((newRefDiffJSON) => {
      cy.intercept('GET', `/api/pipelines/${newRefDiffJSON.id}`, { body: newRefDiffJSON }).as('RefDiffPipeline')
      cy.intercept('POST', '/api/pipelines', { body: newRefDiffJSON }).as('runRefDiffPipeline')
      cy.fixture('new-refdiff-pipeline-tasks').then((newRefDiffPipelineTasksJSON) => {
        cy.intercept('GET', `/api/pipelines/${newRefDiffJSON.id}/tasks`, { body: newRefDiffPipelineTasksJSON }).as('RefDiffPipelineTasks')
      })
    })
    cy.get('input#pipeline-name').type(`{selectAll}{backSpace}COLLECT REFDIFF ${Date.now()}`)
    cy.get('.provider-toggle-switch.switch-refdiff')
      .should('be.visible')
      .click()
      .trigger('mouseleave')

    cy.get('input#refdiff-repo-id').click().type('github:GithubRepo:384111310')
    cy.get('input#refdiff-pair-newref').type('refs/tags/v0.2.0')
    cy.get('input#refdiff-pair-oldref').type('refs/tags/v0.1.0')
    cy.get('button.btn-add-tagpair').click()

    cy.get('button#btn-run-pipeline').click()
    cy.wait('@RefDiffPipeline')
    cy.wait('@RefDiffPipelineTasks')
    cy.wait('@runRefDiffPipeline').then(({ response }) => {
      const NewRefDiffRun = response.body
      cy.url().should('include', `/pipelines/activity/${NewRefDiffRun.id}`)
    })
  })


})
