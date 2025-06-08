import { test, expect } from '@playwright/test'

test('heading hierarchy validation', async ({ page }) => {
  await page.goto('/')
  
  // Test heading hierarchy
  const headings = await page.locator('h1, h2, h3, h4, h5, h6').count()
  expect(headings).toBeGreaterThan(0)
  
  // Check if h1 exists
  await expect(page.locator('h1')).toBeVisible()
})

test('image alt text validation', async ({ page }) => {
  await page.goto('/')
  
  // Check that all images have alt text with debug output
  const images = await page.locator('img').all()
  for (let i = 0; i < images.length; i++) {
    const img = images[i]
    const alt = await img.getAttribute('alt')
    const src = await img.getAttribute('src')
    
    try {
      expect(alt).toBeTruthy()
    } catch (error) {
      console.log(`Image ${i} with src "${src}" is missing alt text. Alt attribute value: "${alt}"`)
      throw error
    }
  }
})

test('form label associations', async ({ page }) => {
  await page.goto('/')
  
  // Check contact form if it exists
  const contactSection = page.locator('#contact')
  if (await contactSection.isVisible()) {
    const inputs = await contactSection.locator('input, textarea, select').all()
    for (const input of inputs) {
      const id = await input.getAttribute('id')
      if (id) {
        await expect(page.locator(`label[for="${id}"]`)).toBeVisible()
      }
    }
  }
})

test('landmarks and semantic structure', async ({ page }) => {
  await page.goto('/')
  
  // Test landmarks exist
  await expect(page.locator('header')).toBeVisible()
  await expect(page.locator('main')).toBeVisible()
  await expect(page.locator('footer')).toBeVisible()
  await expect(page.locator('nav')).toBeVisible()
})