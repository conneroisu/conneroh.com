import { test, expect } from '@playwright/test'

test('HTMX navigation attributes', async ({ page }) => {
  await page.goto('/')
  
  // Test navigation structure
  await expect(page.locator('nav')).toBeVisible()
  
  const projectsLink = page.locator('a[hx-get="/projects"]')
  const postsLink = page.locator('a[hx-get="/posts"]')
  
  if (await projectsLink.count() > 0) {
    await expect(projectsLink.first()).toHaveAttribute('hx-target', '#bodiody')
    await expect(projectsLink.first()).toHaveAttribute('hx-push-url', 'true')
    await expect(projectsLink.first()).toContainText('Projects')
    await expect(projectsLink.first()).toHaveClass(/text-gray-300/)
  }
  
  if (await postsLink.count() > 0) {
    await expect(postsLink.first()).toHaveAttribute('hx-get', '/posts')
    await expect(postsLink.first()).toContainText('Posts')
  }
})

test('mobile menu with Alpine.js', async ({ page }) => {
  await page.goto('/')
  
  // Set mobile viewport
  await page.setViewportSize({ width: 375, height: 667 })
  
  const container = page.locator('[x-data*="isMenuOpen"]')
  const menuButton = page.locator('[\\@click*="isMenuOpen"]')
  const mobileMenu = page.locator('[x-show="isMenuOpen"]')
  
  // Test Alpine.js attributes
  if (await container.count() > 0) {
    await expect(container.first()).toHaveAttribute('x-data')
  }
  
  if (await menuButton.count() > 0) {
    await expect(menuButton.first()).toBeVisible()
  }
  
  if (await mobileMenu.count() > 0) {
    await expect(mobileMenu.first()).toBeAttached()
  }
})

test('navigation links functionality', async ({ page }) => {
  await page.goto('/')
  
  // Test projects navigation
  const projectsLink = page.locator('a[hx-get="/projects"]').first()
  if (await projectsLink.count() > 0) {
    await projectsLink.click()
    
    // Wait for HTMX navigation to complete
    await expect(page).toHaveURL('/projects')
    await expect(page.locator('#bodiody')).toBeVisible()
  }
  
  // Test posts navigation
  const postsLink = page.locator('a[hx-get="/posts"]').first()
  if (await postsLink.count() > 0) {
    await postsLink.click()
    
    // Wait for HTMX navigation to complete
    await expect(page).toHaveURL('/posts')
    await expect(page.locator('#bodiody')).toBeVisible()
  }
})

test('skip link accessibility', async ({ page }) => {
  await page.goto('/')
  
  // Look for skip links (might be visually hidden)
  const skipLinks = page.locator('a[href*="#"]')
  const mainContent = page.locator('main, #bodiody, [role="main"]').first()
  
  // Test main content is present
  await expect(mainContent).toBeVisible()
  
  // Test keyboard navigation works
  await page.keyboard.press('Tab')
  const focusedElement = await page.locator(':focus').first()
  if (await focusedElement.count() > 0) {
    await expect(focusedElement).toBeVisible()
  }
})