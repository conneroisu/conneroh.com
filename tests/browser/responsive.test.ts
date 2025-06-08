import { test, expect } from '@playwright/test'

test('mobile responsiveness (375px)', async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 667 })
  await page.goto('/')
  
  // Test page loads properly on mobile
  await expect(page.locator('#bodiody')).toBeVisible()
  
  // Test mobile menu button is visible
  const mobileMenuButton = page.locator('[\\@click*="isMenuOpen"]')
  if (await mobileMenuButton.count() > 0) {
    await expect(mobileMenuButton.first()).toBeVisible()
  }
  
  // Test mobile menu functionality
  const mobileMenu = page.locator('[x-show="isMenuOpen"]')
  if (await mobileMenu.count() > 0) {
    // Menu should be hidden initially
    await expect(mobileMenu.first()).toBeHidden()
    
    // Click menu button to open
    await mobileMenuButton.first().click()
    
    // Menu should now be visible
    await expect(mobileMenu.first()).toBeVisible()
  }
})

test('tablet responsiveness (768px)', async ({ page }) => {
  await page.setViewportSize({ width: 768, height: 1024 })
  await page.goto('/')
  
  // Test page loads properly on tablet
  await expect(page.locator('#bodiody')).toBeVisible()
  
  // Test navigation is accessible
  const nav = page.locator('nav')
  if (await nav.count() > 0) {
    await expect(nav.first()).toBeVisible()
  }
  
  // Test content layout
  const mainContent = page.locator('main, #bodiody')
  await expect(mainContent.first()).toBeVisible()
})

test('desktop responsiveness (1200px)', async ({ page }) => {
  await page.setViewportSize({ width: 1200, height: 800 })
  await page.goto('/')
  
  // Test page loads properly on desktop
  await expect(page.locator('#bodiody')).toBeVisible()
  
  // Test navigation is visible (not hidden in mobile menu)
  const desktopNav = page.locator('nav a[hx-get]')
  if (await desktopNav.count() > 0) {
    await expect(desktopNav.first()).toBeVisible()
  }
  
  // Test layout elements are properly positioned
  const header = page.locator('header')
  const main = page.locator('main, #bodiody')
  const footer = page.locator('footer')
  
  await expect(main.first()).toBeVisible()
  if (await header.count() > 0) {
    await expect(header.first()).toBeVisible()
  }
  if (await footer.count() > 0) {
    await expect(footer.first()).toBeVisible()
  }
})

test('responsive navigation behavior', async ({ page }) => {
  // Test mobile behavior
  await page.setViewportSize({ width: 375, height: 667 })
  await page.goto('/')
  
  // Check mobile menu exists
  const mobileMenuToggle = page.locator('[\\@click*="isMenuOpen"]')
  if (await mobileMenuToggle.count() > 0) {
    await expect(mobileMenuToggle.first()).toBeVisible()
  }
  
  // Test desktop behavior
  await page.setViewportSize({ width: 1200, height: 800 })
  
  // Desktop navigation should be visible
  const desktopNav = page.locator('nav a[hx-get]')
  if (await desktopNav.count() > 0) {
    await expect(desktopNav.first()).toBeVisible()
  }
})