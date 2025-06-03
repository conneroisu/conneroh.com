import { expect, test, describe } from 'vitest'
import { page, userEvent } from '@vitest/browser/context'

describe('Alpine.js Interactions', () => {
  test('dropdown menus work correctly', async () => {
    await page.goto('http://localhost:8080')
    
    // Find Alpine dropdown (if exists)
    const dropdown = page.locator('[x-data*="open"]').first()
    const dropdownExists = await dropdown.isVisible()
    
    if (dropdownExists) {
      // Click to open dropdown
      const trigger = dropdown.locator('[x-on\\:click], [@click]').first()
      await trigger.click()
      
      // Check if dropdown content is visible
      const dropdownContent = dropdown.locator('[x-show]').first()
      await expect.element(dropdownContent).toBeVisible()
      
      // Click outside to close
      await page.click('body', { position: { x: 0, y: 0 } })
      
      // Verify dropdown closed
      await expect.element(dropdownContent).not.toBeVisible()
    }
  })

  test('modal dialogs work correctly', async () => {
    await page.goto('http://localhost:8080')
    
    // Find modal trigger
    const modalTrigger = page.locator('[x-data*="modal"], [x-data*="dialog"]').first()
    const modalExists = await modalTrigger.isVisible()
    
    if (modalExists) {
      // Open modal
      await modalTrigger.click()
      
      // Check modal is visible
      const modal = page.locator('[x-show*="modal"], [x-show*="dialog"]').first()
      await expect.element(modal).toBeVisible()
      
      // Close modal via close button
      const closeButton = modal.locator('[x-on\\:click*="close"], [@click*="close"]').first()
      if (await closeButton.isVisible()) {
        await closeButton.click()
        await expect.element(modal).not.toBeVisible()
      }
    }
  })

  test('theme toggle with Alpine', async () => {
    await page.goto('http://localhost:8080')
    
    // Find theme toggle managed by Alpine
    const themeToggle = page.locator('[x-data*="theme"], [x-data*="dark"]').first()
    const toggleExists = await themeToggle.isVisible()
    
    if (toggleExists) {
      // Get current theme
      const isDark = await page.evaluate(() => {
        return document.documentElement.classList.contains('dark')
      })
      
      // Click toggle
      await themeToggle.click()
      
      // Verify theme changed
      const isDarkAfter = await page.evaluate(() => {
        return document.documentElement.classList.contains('dark')
      })
      
      expect(isDarkAfter).not.toBe(isDark)
      
      // Verify persistence (if implemented)
      await page.reload()
      
      const isDarkAfterReload = await page.evaluate(() => {
        return document.documentElement.classList.contains('dark')
      })
      
      expect(isDarkAfterReload).toBe(isDarkAfter)
    }
  })

  test('form validation with Alpine', async () => {
    await page.goto('http://localhost:8080')
    
    // Find forms with Alpine validation
    const alpineForm = page.locator('form[x-data]').first()
    const formExists = await alpineForm.isVisible()
    
    if (formExists) {
      // Try to submit empty form
      const submitButton = alpineForm.locator('button[type="submit"]')
      await submitButton.click()
      
      // Check for validation messages
      const errorMessage = alpineForm.locator('[x-show*="error"], [x-text*="error"]').first()
      const hasValidation = await errorMessage.isVisible()
      
      if (hasValidation) {
        expect(await errorMessage.textContent()).toBeTruthy()
      }
    }
  })
})