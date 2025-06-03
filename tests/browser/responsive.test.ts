import { expect, test, describe } from 'vitest'
import { page } from '@vitest/browser/context'

describe('Responsive Design', () => {
  const viewports = [
    { name: 'mobile', width: 375, height: 667 },
    { name: 'tablet', width: 768, height: 1024 },
    { name: 'desktop', width: 1920, height: 1080 },
  ]

  viewports.forEach(({ name, width, height }) => {
    test(`renders correctly on ${name}`, async () => {
      await page.setViewport({ width, height })
      await page.goto('http://localhost:8080')
      
      // Check if navigation is visible/hidden based on viewport
      const mobileMenu = page.locator('[data-testid="mobile-menu"]')
      const desktopNav = page.locator('[data-testid="desktop-nav"]')
      
      if (width < 768) {
        // Mobile view
        await expect.element(mobileMenu).toBeVisible()
        await expect.element(desktopNav).not.toBeVisible()
        
        // Test mobile menu toggle
        const menuToggle = page.locator('[data-testid="menu-toggle"]')
        await menuToggle.click()
        
        const mobileNavItems = page.locator('[data-testid="mobile-nav-items"]')
        await expect.element(mobileNavItems).toBeVisible()
      } else {
        // Desktop view
        await expect.element(desktopNav).toBeVisible()
        const mobileMenuVisible = await mobileMenu.isVisible()
        expect(mobileMenuVisible).toBe(false)
      }
    })
  })

  test('images are responsive', async () => {
    await page.goto('http://localhost:8080')
    
    // Check images have responsive attributes
    const images = page.locator('img')
    const imageCount = await images.count()
    
    for (let i = 0; i < Math.min(imageCount, 5); i++) {
      const img = images.nth(i)
      
      // Check for responsive classes or attributes
      const hasResponsiveClass = await img.evaluate((el) => {
        return el.classList.toString().includes('responsive') || 
               el.classList.toString().includes('w-full') ||
               el.style.maxWidth === '100%'
      })
      
      expect(hasResponsiveClass).toBe(true)
    }
  })

  test('text scales properly', async () => {
    // Mobile
    await page.setViewport({ width: 375, height: 667 })
    await page.goto('http://localhost:8080')
    
    const mobileH1Size = await page.locator('h1').first().evaluate(el => 
      window.getComputedStyle(el).fontSize
    )
    
    // Desktop
    await page.setViewport({ width: 1920, height: 1080 })
    await page.reload()
    
    const desktopH1Size = await page.locator('h1').first().evaluate(el => 
      window.getComputedStyle(el).fontSize
    )
    
    // Desktop font should be larger
    expect(parseInt(desktopH1Size)).toBeGreaterThan(parseInt(mobileH1Size))
  })

  test('grid/flex layouts adapt', async () => {
    await page.goto('http://localhost:8080/projects')
    
    // Desktop - multi column
    await page.setViewport({ width: 1920, height: 1080 })
    const desktopGrid = page.locator('[data-testid="projects-grid"]').first()
    const desktopColumns = await desktopGrid.evaluate(el => {
      const style = window.getComputedStyle(el)
      return style.gridTemplateColumns || style.display
    })
    
    // Mobile - single column
    await page.setViewport({ width: 375, height: 667 })
    const mobileColumns = await desktopGrid.evaluate(el => {
      const style = window.getComputedStyle(el)
      return style.gridTemplateColumns || style.display
    })
    
    expect(desktopColumns).not.toBe(mobileColumns)
  })
})