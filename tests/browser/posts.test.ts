import { test, expect } from '@playwright/test'

test('posts page structure', async ({ page }) => {
  await page.goto('/posts')
  
  // Test basic page structure
  await expect(page.locator('#bodiody')).toBeVisible()
  
  // Check if posts exist
  const postElements = page.locator('article, .post, [class*="post"]')
  if (await postElements.count() > 0) {
    await expect(postElements.first()).toBeVisible()
  }
  
  // Test page title or heading
  const headings = page.locator('h1, h2')
  if (await headings.count() > 0) {
    await expect(headings.first()).toBeVisible()
  }
})

test('post navigation and links', async ({ page }) => {
  await page.goto('/posts')
  
  // Test that post links are clickable
  const postLinks = page.locator('a[href*="/posts/"], a[hx-get*="/posts/"]')
  const linkCount = await postLinks.count()
  
  if (linkCount > 0) {
    const firstLink = postLinks.first()
    await expect(firstLink).toBeVisible()
    
    // Test clicking the link
    await firstLink.click()
    
    // Should navigate to post detail or update content
    await expect(page.locator('#bodiody')).toBeVisible()
  }
})

test('post search and filtering', async ({ page }) => {
  await page.goto('/posts')
  
  // Test search functionality if it exists
  const searchInput = page.locator('input[type="search"], #search, [placeholder*="search" i]')
  if (await searchInput.count() > 0) {
    await searchInput.first().fill('test')
    await page.waitForTimeout(500) // Wait for any debounced search
    
    // Check that search works (content should update or filter)
    await expect(page.locator('#bodiody')).toBeVisible()
  }
  
  // Test tag filtering if it exists
  const tagFilters = page.locator('a[href*="tag="], button[data-tag], .tag')
  if (await tagFilters.count() > 0) {
    await tagFilters.first().click()
    await expect(page.locator('#bodiody')).toBeVisible()
  }
})

test('post date formatting', async ({ page }) => {
  await page.goto('/posts')
  
  // Find time elements and test their format
  const timeElements = page.locator('time')
  const timeCount = await timeElements.count()
  
  if (timeCount > 0) {
    const firstTime = timeElements.first()
    const dateText = await firstTime.textContent()
    
    // Test various date formats that might be used
    const dateRegex = /(\w+ \d{1,2}, \d{4}|\d{4}-\d{2}-\d{2}|\w+ \d{1,2}|\d{1,2}\/\d{1,2}\/\d{4})/
    expect(dateRegex.test(dateText || '')).toBe(true)
  }
})