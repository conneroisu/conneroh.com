import { expect, test } from 'vitest'
import { page } from '@vitest/browser/context'

test('homepage loads correctly', async () => {
  // Navigate to the homepage
  await page.goto('http://localhost:8080')
  
  // Check if the main title exists
  await expect.element(page.getByText('Conner Ohnesorge')).toBeInTheDocument()
  
  // Check navigation elements
  await expect.element(page.getByText('Home')).toBeInTheDocument()
  await expect.element(page.getByText('Projects')).toBeInTheDocument()
  await expect.element(page.getByText('Posts')).toBeInTheDocument()
})

test('navigation works correctly', async () => {
  await page.goto('http://localhost:8080')
  
  // Click on Projects link
  const projectsLink = page.getByText('Projects')
  await projectsLink.click()
  
  // Verify URL changed
  expect(page.url()).toContain('/projects')
  
  // Verify projects page content loads
  await expect.element(page.getByText('All Projects')).toBeInTheDocument()
})

test('dark mode toggle works', async () => {
  await page.goto('http://localhost:8080')
  
  // Find and click the theme toggle button
  const themeToggle = page.getByRole('button', { name: /theme/i })
  
  // Get initial theme
  const initialTheme = await page.evaluate(() => document.documentElement.classList.contains('dark'))
  
  // Click toggle
  await themeToggle.click()
  
  // Verify theme changed
  const newTheme = await page.evaluate(() => document.documentElement.classList.contains('dark'))
  expect(newTheme).not.toBe(initialTheme)
})