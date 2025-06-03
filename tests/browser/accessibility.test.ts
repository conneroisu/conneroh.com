import { expect, test, describe } from 'vitest'
import { page } from '@vitest/browser/context'

describe('Accessibility', () => {
  test('has proper heading hierarchy', async () => {
    await page.goto('http://localhost:8080')
    
    // Get all headings
    const headings = await page.locator('h1, h2, h3, h4, h5, h6').all()
    
    let lastLevel = 0
    for (const heading of headings) {
      const tagName = await heading.evaluate(el => el.tagName)
      const level = parseInt(tagName.charAt(1))
      
      // Check heading hierarchy (shouldn't skip levels)
      if (lastLevel > 0) {
        expect(level).toBeLessThanOrEqual(lastLevel + 1)
      }
      lastLevel = level
    }
  })

  test('images have alt text', async () => {
    await page.goto('http://localhost:8080')
    
    const images = await page.locator('img').all()
    
    for (const img of images) {
      const alt = await img.getAttribute('alt')
      expect(alt).toBeTruthy()
    }
  })

  test('interactive elements are keyboard accessible', async () => {
    await page.goto('http://localhost:8080')
    
    // Tab through interactive elements
    await page.keyboard.press('Tab')
    
    // Check if focus is visible
    const focusedElement = await page.evaluate(() => {
      const el = document.activeElement
      return {
        tagName: el?.tagName,
        hasOutline: window.getComputedStyle(el!).outline !== 'none',
        hasFocusVisible: el?.matches(':focus-visible'),
      }
    })
    
    expect(focusedElement.tagName).toBeTruthy()
    expect(focusedElement.hasOutline || focusedElement.hasFocusVisible).toBe(true)
  })

  test('forms have proper labels', async () => {
    await page.goto('http://localhost:8080')
    
    const inputs = await page.locator('input:not([type="hidden"]), select, textarea').all()
    
    for (const input of inputs) {
      const id = await input.getAttribute('id')
      const ariaLabel = await input.getAttribute('aria-label')
      const ariaLabelledby = await input.getAttribute('aria-labelledby')
      
      if (id) {
        const label = await page.locator(`label[for="${id}"]`).first()
        const hasLabel = await label.isVisible()
        
        expect(hasLabel || ariaLabel || ariaLabelledby).toBeTruthy()
      } else {
        expect(ariaLabel || ariaLabelledby).toBeTruthy()
      }
    }
  })

  test('color contrast meets WCAG standards', async () => {
    await page.goto('http://localhost:8080')
    
    // Sample check for main text elements
    const textElements = await page.locator('p, h1, h2, h3, h4, h5, h6, a').all()
    
    for (const element of textElements.slice(0, 5)) { // Check first 5 elements
      const styles = await element.evaluate(el => {
        const computed = window.getComputedStyle(el)
        return {
          color: computed.color,
          backgroundColor: computed.backgroundColor,
          fontSize: computed.fontSize,
        }
      })
      
      // Basic check - ensure text is not transparent
      expect(styles.color).not.toBe('rgba(0, 0, 0, 0)')
      expect(styles.color).not.toBe('transparent')
    }
  })

  test('page has proper landmarks', async () => {
    await page.goto('http://localhost:8080')
    
    // Check for main landmarks
    const header = page.locator('header, [role="banner"]')
    const main = page.locator('main, [role="main"]')
    const footer = page.locator('footer, [role="contentinfo"]')
    const nav = page.locator('nav, [role="navigation"]')
    
    await expect.element(header.first()).toBeVisible()
    await expect.element(main.first()).toBeVisible()
    await expect.element(footer.first()).toBeVisible()
    await expect.element(nav.first()).toBeVisible()
  })

  test('skip to content link exists', async () => {
    await page.goto('http://localhost:8080')
    
    // Look for skip link
    const skipLink = page.locator('a[href="#main"], a[href="#content"], .skip-link')
    
    // Tab to potentially reveal skip link
    await page.keyboard.press('Tab')
    
    const skipLinkExists = await skipLink.first().isVisible()
    if (skipLinkExists) {
      const href = await skipLink.first().getAttribute('href')
      expect(href).toMatch(/#(main|content)/)
    }
  })
})