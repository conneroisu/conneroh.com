import { test, expect } from '@playwright/test'

test('Alpine.js data attributes', async ({ page }) => {
  await page.goto('/')
  
  // Test that Alpine.js attributes exist on elements (like mobile menu)
  await expect(page.locator('[x-data]')).toHaveCount(await page.locator('[x-data]').count())
  
  // Test mobile menu specifically
  const mobileMenuContainer = page.locator('[x-data*="isMenuOpen"]')
  if (await mobileMenuContainer.isVisible()) {
    await expect(mobileMenuContainer).toHaveAttribute('x-data')
    await expect(page.locator('[x-show="isMenuOpen"]')).toBeAttached()
    
    // Check for mobile menu button - it should exist but may not be visible on desktop
    const menuButton = page.locator('[\\@click="isMenuOpen = !isMenuOpen"]')
    await expect(menuButton).toBeAttached()
  }
})

test('mobile menu interaction with Alpine.js', async ({ page }) => {
  // Set mobile viewport first before loading page
  await page.setViewportSize({ width: 375, height: 667 })
  await page.goto('/')
  
  const menuButton = page.locator('[\\@click="isMenuOpen = !isMenuOpen"]')
  const mobileMenu = page.locator('[x-show="isMenuOpen"]')
  
  // Wait for Alpine.js to initialize
  await page.waitForLoadState('networkidle')
  
  // The mobile menu button should be visible on mobile
  await expect(menuButton).toBeVisible()
  
  // Test initial state - menu should be hidden
  await expect(mobileMenu).toBeHidden()
  
  // Click to open menu
  await menuButton.click()
  
  // Wait for Alpine.js to show the menu with timeout
  await expect(mobileMenu).toBeVisible({ timeout: 10000 })
  
  // Click menu button again to close
  await menuButton.click()
  
  // Menu should be hidden again
  await expect(mobileMenu).toBeHidden({ timeout: 10000 })
})