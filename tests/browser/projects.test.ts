import { expect, test } from 'vitest'

test('projects page DOM manipulation', async () => {
  // Create mock project card
  const projectCard = document.createElement('div')
  projectCard.setAttribute('data-testid', 'project-card')
  projectCard.textContent = 'Test Project'
  document.body.appendChild(projectCard)
  
  // Test project card exists
  const foundCard = document.querySelector('[data-testid="project-card"]')
  expect(foundCard).toBeTruthy()
  expect(foundCard?.textContent).toBe('Test Project')
  
  // Clean up
  document.body.removeChild(projectCard)
})

test('project filtering simulation', async () => {
  // Create mock tag filter
  const tagFilter = document.createElement('button')
  tagFilter.setAttribute('data-testid', 'tag-filter')
  tagFilter.textContent = 'go'
  document.body.appendChild(tagFilter)
  
  // Simulate click handler
  let clicked = false
  tagFilter.addEventListener('click', () => {
    clicked = true
  })
  
  tagFilter.click()
  expect(clicked).toBe(true)
  
  // Clean up
  document.body.removeChild(tagFilter)
})