import { expect, test, describe } from 'vitest'
import { page } from '@vitest/browser/context'

describe('Projects Page', () => {
  test('displays project list correctly', async () => {
    await page.goto('http://localhost:8080/projects')
    
    // Check page title
    await expect.element(page.getByText('All Projects')).toBeInTheDocument()
    
    // Check if projects are rendered
    const projectCards = page.locator('[data-testid="project-card"]')
    const count = await projectCards.count()
    expect(count).toBeGreaterThan(0)
  })

  test('filters projects by tag', async () => {
    await page.goto('http://localhost:8080/projects')
    
    // Find and click a tag filter
    const tagButton = page.locator('[data-testid="tag-filter"]').first()
    const tagText = await tagButton.textContent()
    
    await tagButton.click()
    
    // Verify URL updated
    expect(page.url()).toContain('/tags/')
    
    // Verify filtered view shows
    await expect.element(page.getByText(`Projects tagged with "${tagText}"`)).toBeInTheDocument()
  })

  test('project card interactions', async () => {
    await page.goto('http://localhost:8080/projects')
    
    // Test hover effect on project card
    const projectCard = page.locator('[data-testid="project-card"]').first()
    
    // Get initial styles
    const initialOpacity = await projectCard.evaluate(el => 
      window.getComputedStyle(el).opacity
    )
    
    // Hover over card
    await projectCard.hover()
    
    // Check if hover effect applied
    const hoverOpacity = await projectCard.evaluate(el => 
      window.getComputedStyle(el).opacity
    )
    
    expect(hoverOpacity).not.toBe(initialOpacity)
  })

  test('pagination works correctly', async () => {
    await page.goto('http://localhost:8080/projects')
    
    // Check if pagination exists (if there are enough projects)
    const paginationExists = await page.locator('[data-testid="pagination"]').isVisible()
    
    if (paginationExists) {
      // Click next page
      const nextButton = page.locator('[data-testid="pagination-next"]')
      await nextButton.click()
      
      // Verify URL updated
      expect(page.url()).toContain('page=2')
      
      // Verify different content loaded
      const firstProjectTitle = await page.locator('[data-testid="project-title"]').first().textContent()
      
      // Go back to first page
      const prevButton = page.locator('[data-testid="pagination-prev"]')
      await prevButton.click()
      
      const newFirstProjectTitle = await page.locator('[data-testid="project-title"]').first().textContent()
      expect(firstProjectTitle).not.toBe(newFirstProjectTitle)
    }
  })
})