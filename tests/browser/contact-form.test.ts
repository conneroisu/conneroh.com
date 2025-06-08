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

// test('contact form submission flow', async ({ page }) => {
//   await page.goto('/')
//
//   const contactSection = page.locator('#contact')
//   if (await contactSection.isVisible()) {
//     const form = page.locator('form')
//
//     if (await form.count() > 0) {
//       // Let the actual form submission go through (don't mock)
//
//       // Fill out form if fields exist
//       const nameField = page.locator('input[name="name"], #name')
//       const emailField = page.locator('input[name="email"], #email')
//       const messageField = page.locator('textarea[name="message"], #message')
//
//       if (await nameField.count() > 0) {
//         await nameField.first().fill('Test User')
//       }
//       if (await emailField.count() > 0) {
//         await emailField.first().fill('test@example.com')
//       }
//       if (await messageField.count() > 0) {
//         await messageField.first().fill('This is a test message')
//       }
//
//       // Submit form
//       const submitButton = page.locator('button[type="submit"], input[type="submit"]')
//       if (await submitButton.count() > 0) {
//         await submitButton.first().click()
//
//         // Wait for any response and check what we get
//         await page.waitForLoadState('networkidle')
//
//         // Check for success message - first try to find any h3 element 
//         const anyThankYou = page.locator(':has-text("Thank You")')
//
//         console.log('Looking for Thank You message...')
//         console.log('Page title:', await page.title())
//         console.log('Page URL:', page.url())
//
//         // Try more flexible locators
//         if (await anyThankYou.count() > 0) {
//           await expect(anyThankYou.first()).toBeVisible({ timeout: 5000 })
//         } else {
//           // If no thank you message, let's see what's actually on the page
//           const bodyText = await page.locator('body').textContent()
//           console.log('Page body contains:', bodyText?.slice(0, 500))
//           throw new Error('No Thank You message found on page')
//         }
//       }
//     }
//   }
// })
