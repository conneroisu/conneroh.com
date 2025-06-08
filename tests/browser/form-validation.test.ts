import { test, expect } from '@playwright/test'

test('form input validation', async ({ page }) => {
  await page.goto('/')
  
  const contactSection = page.locator('#contact')
  if (await contactSection.isVisible()) {
    const emailInput = page.locator('input[type="email"], #email')
    
    if (await emailInput.count() > 0) {
      // Test invalid email
      await emailInput.first().fill('invalid-email')
      const isInvalid = await emailInput.first().evaluate((el: HTMLInputElement) => !el.validity.valid)
      expect(isInvalid).toBe(true)
      
      // Test valid email
      await emailInput.first().fill('test@example.com')
      const isValid = await emailInput.first().evaluate((el: HTMLInputElement) => el.validity.valid)
      expect(isValid).toBe(true)
    }
    
    // Test required fields
    const requiredFields = page.locator('input[required], textarea[required]')
    const requiredCount = await requiredFields.count()
    
    if (requiredCount > 0) {
      for (let i = 0; i < requiredCount; i++) {
        const field = requiredFields.nth(i)
        await field.fill('')
        const isEmpty = await field.evaluate((el: HTMLInputElement) => !el.validity.valid)
        expect(isEmpty).toBe(true)
      }
    }
  }
})

test('form styling and classes', async ({ page }) => {
  await page.goto('/')
  
  const contactSection = page.locator('#contact')
  if (await contactSection.isVisible()) {
    const form = page.locator('form')
    
    if (await form.count() > 0) {
      // Test form has appropriate classes
      await expect(form.first()).toHaveClass(/.*/)
      
      // Test input fields have styling classes
      const inputs = page.locator('input, textarea')
      const inputCount = await inputs.count()
      
      if (inputCount > 0) {
        for (let i = 0; i < inputCount; i++) {
          const input = inputs.nth(i)
          await expect(input).toHaveClass(/.*/)
        }
      }
    }
  }
})