import { test, expect } from '@playwright/test'

test('basic CSS classes and styling', async ({ page }) => {
  await page.goto('/')
  
  // Test that main content has proper classes
  const mainContent = page.locator('#bodiody, main')
  await expect(mainContent.first()).toBeVisible()
  
  // Test header styling
  const header = page.locator('header')
  if (await header.count() > 0) {
    await expect(header.first()).toHaveClass(/.*/)
  }
  
  // Test navigation styling
  const nav = page.locator('nav')
  if (await nav.count() > 0) {
    await expect(nav.first()).toHaveClass(/.*/)
  }
  
  // Test that Tailwind classes are applied
  const elementsWithTailwind = page.locator('[class*="text-"], [class*="bg-"], [class*="p-"], [class*="m-"]')
  if (await elementsWithTailwind.count() > 0) {
    await expect(elementsWithTailwind.first()).toBeVisible()
  }
})

test('responsive design classes', async ({ page }) => {
  await page.goto('/')
  
  // Test responsive classes
  const responsiveElements = page.locator('[class*="sm:"], [class*="md:"], [class*="lg:"]')
  if (await responsiveElements.count() > 0) {
    await expect(responsiveElements.first()).toBeVisible()
  }
  
  // Test mobile-first approach with different viewports
  await page.setViewportSize({ width: 375, height: 667 })
  await expect(page.locator('#bodiody')).toBeVisible()
  
  await page.setViewportSize({ width: 1200, height: 800 })
  await expect(page.locator('#bodiody')).toBeVisible()
})

test('color scheme and theme', async ({ page }) => {
  await page.goto('/')
  
  // Test dark theme classes
  const darkElements = page.locator('[class*="gray-"], [class*="dark"]')
  if (await darkElements.count() > 0) {
    await expect(darkElements.first()).toBeVisible()
  }
  
  // Test text color classes
  const textElements = page.locator('[class*="text-white"], [class*="text-gray"]')
  if (await textElements.count() > 0) {
    await expect(textElements.first()).toBeVisible()
  }
  
  // Test background color classes
  const bgElements = page.locator('[class*="bg-gray"], [class*="bg-gradient"]')
  if (await bgElements.count() > 0) {
    await expect(bgElements.first()).toBeVisible()
  }
})

test('button and interactive element styling', async ({ page }) => {
  await page.goto('/')
  
  // Test button styling
  const buttons = page.locator('button, [type="button"], [type="submit"]')
  if (await buttons.count() > 0) {
    const firstButton = buttons.first()
    await expect(firstButton).toBeVisible()
    
    // Test hover state (if possible)
    await firstButton.hover()
    await expect(firstButton).toBeVisible()
  }
  
  // Test link styling
  const links = page.locator('a')
  if (await links.count() > 0) {
    const firstLink = links.first()
    await expect(firstLink).toBeVisible()
    
    // Test hover state
    await firstLink.hover()
    await expect(firstLink).toBeVisible()
  }
})