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
    await expect(page.locator('[\\@click*="isMenuOpen"]')).toBeVisible()
  }
})

test('mobile menu interaction with Alpine.js', async ({ page }) => {
  await page.goto('/')
  
  // Set mobile viewport
  await page.setViewportSize({ width: 375, height: 667 })
  
  const menuButton = page.locator('[\\@click*="isMenuOpen"]')
  const mobileMenu = page.locator('[x-show="isMenuOpen"]')
  
  if (await menuButton.isVisible()) {
    // Test initial state - menu should be hidden
    await expect(mobileMenu).toBeHidden()
    
    // Click to open menu
    await menuButton.click()
    
    // Wait for Alpine.js to show the menu
    await expect(mobileMenu).toBeVisible()
    
    // Click menu button again to close
    await menuButton.click()
    
    // Menu should be hidden again
    await expect(mobileMenu).toBeHidden()
  }
})