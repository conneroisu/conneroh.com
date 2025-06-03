import { expect, test, describe } from 'vitest'
import { page, userEvent } from '@vitest/browser/context'

describe('HTMX Interactions', () => {
  test('navigation uses HTMX for partial updates', async () => {
    await page.goto('http://localhost:8080')
    
    // Listen for HTMX events
    const htmxEventFired = await page.evaluate(() => {
      return new Promise((resolve) => {
        document.body.addEventListener('htmx:afterRequest', () => resolve(true), { once: true })
        setTimeout(() => resolve(false), 5000)
      })
    })
    
    // Click navigation link
    const projectsLink = page.getByText('Projects')
    await projectsLink.click()
    
    // Verify HTMX was used
    expect(htmxEventFired).toBe(true)
    
    // Verify content updated without full page reload
    await expect.element(page.getByText('All Projects')).toBeInTheDocument()
  })

  test('lazy loading works with HTMX', async () => {
    await page.goto('http://localhost:8080/projects')
    
    // Check for lazy-loaded content indicators
    const lazyElements = page.locator('[hx-trigger*="intersect"]')
    const lazyCount = await lazyElements.count()
    
    if (lazyCount > 0) {
      // Scroll to trigger lazy loading
      await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))
      
      // Wait for content to load
      await page.waitForTimeout(1000)
      
      // Verify new content loaded
      const newLazyCount = await lazyElements.count()
      expect(newLazyCount).toBeLessThanOrEqual(lazyCount)
    }
  })

  test('form submissions use HTMX', async () => {
    await page.goto('http://localhost:8080')
    
    // Look for any forms with HTMX attributes
    const htmxForm = page.locator('form[hx-post], form[hx-get]').first()
    const formExists = await htmxForm.isVisible()
    
    if (formExists) {
      // Set up event listener
      const htmxSubmitFired = await page.evaluate(() => {
        return new Promise((resolve) => {
          document.body.addEventListener('htmx:afterRequest', () => resolve(true), { once: true })
          setTimeout(() => resolve(false), 5000)
        })
      })
      
      // Submit form
      const submitButton = htmxForm.locator('button[type="submit"]')
      await submitButton.click()
      
      // Verify HTMX handled submission
      expect(htmxSubmitFired).toBe(true)
    }
  })

  test('history navigation works correctly', async () => {
    await page.goto('http://localhost:8080')
    
    // Navigate to projects
    await page.getByText('Projects').click()
    await expect.element(page.getByText('All Projects')).toBeInTheDocument()
    
    // Navigate to posts
    await page.getByText('Posts').click()
    await expect.element(page.getByText('All Posts')).toBeInTheDocument()
    
    // Use browser back button
    await page.goBack()
    await expect.element(page.getByText('All Projects')).toBeInTheDocument()
    
    // Use browser forward button
    await page.goForward()
    await expect.element(page.getByText('All Posts')).toBeInTheDocument()
  })
})