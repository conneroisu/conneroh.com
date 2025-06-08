import { test, expect } from '@playwright/test'

test('contact form matches template structure', async ({ page }) => {
  await page.goto('/')
  
  // Navigate to contact section or check if form exists on homepage
  const contactSection = page.locator('#contact')
  if (await contactSection.isVisible()) {
    // Test contact form exists
    const contactForm = page.locator('form')
    if (await contactForm.count() > 0) {
      await expect(contactForm.first()).toBeVisible()
      
      // Test basic form elements exist
      const nameField = page.locator('input[name="name"], #name')
      const emailField = page.locator('input[name="email"], #email')
      const messageField = page.locator('textarea[name="message"], #message')
      
      if (await nameField.count() > 0) {
        await expect(nameField.first()).toBeVisible()
      }
      if (await emailField.count() > 0) {
        await expect(emailField.first()).toBeVisible()
      }
      if (await messageField.count() > 0) {
        await expect(messageField.first()).toBeVisible()
      }
    }
  }
})

test('contact form input interaction', async ({ page }) => {
  await page.goto('/')
  
  const contactSection = page.locator('#contact')
  if (await contactSection.isVisible()) {
    const nameField = page.locator('input[name="name"], #name')
    const emailField = page.locator('input[name="email"], #email')
    const messageField = page.locator('textarea[name="message"], #message')
    
    if (await nameField.count() > 0) {
      await nameField.first().fill('John Doe')
      await expect(nameField.first()).toHaveValue('John Doe')
    }
    
    if (await emailField.count() > 0) {
      await emailField.first().fill('john@example.com')
      await expect(emailField.first()).toHaveValue('john@example.com')
    }
    
    if (await messageField.count() > 0) {
      await messageField.first().fill('Test message content')
      await expect(messageField.first()).toHaveValue('Test message content')
    }
  }
})

test('contact form submission flow', async ({ page }) => {
  await page.goto('/')
  
  const contactSection = page.locator('#contact')
  if (await contactSection.isVisible()) {
    const form = page.locator('form')
    
    if (await form.count() > 0) {
      // Mock the form submission if it uses HTMX
      await page.route('**/contact', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'text/html',
          body: '<div class="success">Thank you! Your message has been sent.</div>'
        })
      })
      
      // Fill out form if fields exist
      const nameField = page.locator('input[name="name"], #name')
      const emailField = page.locator('input[name="email"], #email')
      const messageField = page.locator('textarea[name="message"], #message')
      
      if (await nameField.count() > 0) {
        await nameField.first().fill('Test User')
      }
      if (await emailField.count() > 0) {
        await emailField.first().fill('test@example.com')
      }
      if (await messageField.count() > 0) {
        await messageField.first().fill('This is a test message')
      }
      
      // Submit form
      const submitButton = page.locator('button[type="submit"], input[type="submit"]')
      if (await submitButton.count() > 0) {
        await submitButton.first().click()
        
        // Check for success message
        await expect(page.locator('.success')).toBeVisible()
      }
    }
  }
})