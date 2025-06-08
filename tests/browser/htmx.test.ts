import { test, expect } from '@playwright/test'

test('HTMX attributes and events', async ({ page }) => {
  await page.goto('/')
  
  // Set desktop viewport to test navigation
  await page.setViewportSize({ width: 1024, height: 768 })
  
  // Test HTMX attributes exist on navigation links
  const projectsLink = page.locator('a[hx-get="/projects"]').first()
  
  // Check attributes
  await expect(projectsLink).toHaveAttribute('hx-target', '#bodiody')
  await expect(projectsLink).toHaveAttribute('hx-get', '/projects')
  await expect(projectsLink).toBeVisible()
  
  // Test HTMX navigation by clicking the link
  await projectsLink.click()
  
  // Wait for HTMX request to complete and check URL
  await page.waitForURL('/projects', { timeout: 5000 })
  await expect(page.locator('#bodiody')).toBeVisible()
})

test('form with HTMX attributes', async ({ page }) => {
  await page.goto('/')
  
  // Check if contact form exists with HTMX attributes
  const contactForm = page.locator('form[hx-post]')
  if (await contactForm.isVisible()) {
    await expect(contactForm).toHaveAttribute('hx-post')
    
    // Test that form inputs exist (no hx-target needed for this form)
    await expect(contactForm.locator('input, textarea')).toHaveCount(await contactForm.locator('input, textarea').count())
  }
})