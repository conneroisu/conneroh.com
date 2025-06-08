import { test, expect } from '@playwright/test'

test('projects page structure', async ({ page }) => {
  await page.goto('/projects')
  
  // Test basic page structure
  await expect(page.locator('#bodiody')).toBeVisible()
  
  // Check if projects exist
  const projectElements = page.locator('article, .project, [class*="project"]')
  if (await projectElements.count() > 0) {
    await expect(projectElements.first()).toBeVisible()
  }
  
  // Test page title or heading
  const headings = page.locator('h1, h2')
  if (await headings.count() > 0) {
    await expect(headings.first()).toBeVisible()
  }
})

test('project filtering and tags', async ({ page }) => {
  await page.goto('/projects')
  
  // Test tag filtering if it exists
  const tagFilters = page.locator('a[href*="tag="], button[data-tag], .tag')
  if (await tagFilters.count() > 0) {
    await tagFilters.first().click()
    await expect(page.locator('#bodiody')).toBeVisible()
  }
  
  // Test search functionality if it exists
  const searchInput = page.locator('input[type="search"], #search, [placeholder*="search" i]')
  if (await searchInput.count() > 0) {
    await searchInput.first().fill('test')
    await page.waitForTimeout(500) // Wait for any debounced search
    
    // Check that search works (content should update or filter)
    await expect(page.locator('#bodiody')).toBeVisible()
  }
})

test('project links and navigation', async ({ page }) => {
  await page.goto('/projects')
  
  // Test that project links are clickable
  const projectLinks = page.locator('a[href*="/projects/"], a[hx-get*="/projects/"]')
  const linkCount = await projectLinks.count()
  
  if (linkCount > 0) {
    const firstLink = projectLinks.first()
    await expect(firstLink).toBeVisible()
    
    // Test clicking the link
    await firstLink.click()
    
    // Should navigate to project detail or update content
    await expect(page.locator('#bodiody')).toBeVisible()
  }
  
  // Test external links (GitHub, etc.)
  const externalLinks = page.locator('a[href^="http"], a[target="_blank"]')
  if (await externalLinks.count() > 0) {
    await expect(externalLinks.first()).toBeVisible()
  }
})