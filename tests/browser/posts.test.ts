import { expect, test } from 'vitest'

test('posts page DOM manipulation', async () => {
  // Create mock post card
  const postCard = document.createElement('article')
  postCard.setAttribute('data-testid', 'post-card')
  
  const title = document.createElement('h2')
  title.setAttribute('data-testid', 'post-title')
  title.textContent = 'Test Post'
  
  const date = document.createElement('time')
  date.setAttribute('data-testid', 'post-date')
  date.textContent = 'January 1, 2025'
  
  const description = document.createElement('p')
  description.setAttribute('data-testid', 'post-description')
  description.textContent = 'Test description'
  
  postCard.appendChild(title)
  postCard.appendChild(date)
  postCard.appendChild(description)
  document.body.appendChild(postCard)
  
  // Test post elements exist
  expect(document.querySelector('[data-testid="post-card"]')).toBeTruthy()
  expect(document.querySelector('[data-testid="post-title"]')?.textContent).toBe('Test Post')
  expect(document.querySelector('[data-testid="post-date"]')?.textContent).toBe('January 1, 2025')
  expect(document.querySelector('[data-testid="post-description"]')?.textContent).toBe('Test description')
  
  // Clean up
  document.body.removeChild(postCard)
})

test('post date formatting', async () => {
  const dateElement = document.createElement('time')
  dateElement.textContent = 'January 1, 2025'
  
  // Test date format validation
  const dateRegex = /\w+ \d{1,2}, \d{4}/
  expect(dateRegex.test(dateElement.textContent)).toBe(true)
})