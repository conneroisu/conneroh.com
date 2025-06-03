import { expect, test, describe } from 'vitest'
import { page } from '@vitest/browser/context'

describe('Posts Page', () => {
  test('displays blog posts correctly', async () => {
    await page.goto('http://localhost:8080/posts')
    
    // Check page title
    await expect.element(page.getByText('All Posts')).toBeInTheDocument()
    
    // Check if posts are rendered
    const postCards = page.locator('[data-testid="post-card"]')
    const count = await postCards.count()
    expect(count).toBeGreaterThan(0)
  })

  test('shows post metadata', async () => {
    await page.goto('http://localhost:8080/posts')
    
    const firstPost = page.locator('[data-testid="post-card"]').first()
    
    // Check for required metadata
    const title = firstPost.locator('[data-testid="post-title"]')
    const date = firstPost.locator('[data-testid="post-date"]')
    const description = firstPost.locator('[data-testid="post-description"]')
    
    await expect.element(title).toBeVisible()
    await expect.element(date).toBeVisible()
    await expect.element(description).toBeVisible()
    
    // Verify date format
    const dateText = await date.textContent()
    expect(dateText).toMatch(/\w+ \d{1,2}, \d{4}/)
  })

  test('navigates to individual post', async () => {
    await page.goto('http://localhost:8080/posts')
    
    // Get first post title
    const firstPostTitle = await page.locator('[data-testid="post-title"]').first().textContent()
    
    // Click on the post
    const firstPostLink = page.locator('[data-testid="post-link"]').first()
    await firstPostLink.click()
    
    // Verify navigation
    expect(page.url()).toContain('/posts/')
    
    // Verify post content loaded
    await expect.element(page.getByText(firstPostTitle.trim())).toBeInTheDocument()
    
    // Check for post content elements
    await expect.element(page.locator('[data-testid="post-content"]')).toBeVisible()
  })

  test('filters posts by tag', async () => {
    await page.goto('http://localhost:8080/posts')
    
    // Find a tag on a post
    const tagLink = page.locator('[data-testid="post-tag"]').first()
    const tagExists = await tagLink.isVisible()
    
    if (tagExists) {
      const tagText = await tagLink.textContent()
      await tagLink.click()
      
      // Verify filtered view
      expect(page.url()).toContain('/tags/')
      await expect.element(page.getByText(`Posts tagged with "${tagText.trim()}"`)).toBeInTheDocument()
    }
  })
})