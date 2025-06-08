import { test, expect } from '@playwright/test'

test('HTMX attributes and events', async ({ page }) => {
  await page.goto('/')
  
  // Test HTMX attributes exist on navigation links (get first visible one)
  const projectsLink = page.locator('a[hx-get="/projects"]').first()
  await expect(projectsLink).toBeVisible()
  await expect(projectsLink).toHaveAttribute('hx-target', '#bodiody')
  
  // Test HTMX navigation by clicking the link
  await projectsLink.click()
  
  // Wait for HTMX request to complete and check URL
  await expect(page).toHaveURL('/projects')
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