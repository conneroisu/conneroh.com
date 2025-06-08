import { test, expect } from '@playwright/test'

test('homepage structure matches template', async ({ page }) => {
  await page.goto('/')

  // Test homepage structure
  await expect(page.locator('#bodiody')).toBeVisible()

  // Test name header with debug output
  const nameHeader = page.locator('h1[aria-label="Name"]')
  try {
    await expect(nameHeader).toContainText('Conner Ohnesorge')
  } catch (error) {
    const actualText = await nameHeader.textContent()
    console.log(`Name header selector 'h1[aria-label="Name"]' contains: "${actualText}"`)
    throw error
  }

  await expect(page.locator('p[aria-label="summary"]')).toBeVisible()
  await expect(page.locator('#projects')).toBeVisible()
  await expect(page.locator('#contact')).toBeVisible()
})

test('navigation structure', async ({ page }) => {
  await page.goto('/')

  // Test navigation with debug output
  const homeLink = page.locator('a[aria-label="Back to Home"]')
  try {
    await expect(homeLink).toContainText('Conner Ohnesorge')
  } catch (error) {
    const actualText = await homeLink.textContent()
    console.log(`Home link selector 'a[aria-label="Back to Home"]' contains: "${actualText}"`)
    throw error
  }

  // Desktop projects link - use first visible one
  const projectsLink = page.locator('a[hx-get="/projects"]').first()
  try {
    await expect(projectsLink).toContainText('Projects')
  } catch (error) {
    const actualText = await projectsLink.textContent()
    console.log(`Projects link selector 'a[hx-get="/projects"]' contains: "${actualText}"`)
    throw error
  }

  // Desktop posts link - use first visible one
  const postsLink = page.locator('a[hx-get="/posts"]').first()
  try {
    await expect(postsLink).toContainText('Posts')
  } catch (error) {
    const actualText = await postsLink.textContent()
    console.log(`Posts link selector 'a[hx-get="/posts"]' contains: "${actualText}"`)
    throw error
  }
})

test('mobile menu functionality', async ({ page }) => {
  await page.goto('/')

  // Set mobile viewport
  await page.setViewportSize({ width: 375, height: 667 })

  // Test mobile menu elements
  await expect(page.locator('img[src*="menu.svg"]')).toBeVisible()
  await expect(page.locator('div[x-show="isMenuOpen"]')).toBeAttached()

  // Test that mobile menu initially hidden
  await expect(page.locator('div[x-show="isMenuOpen"]')).toBeHidden()

  // Click menu button to open menu
  await page.locator('img[src*="menu.svg"]').click()

  // Test that mobile menu is now visible
  await expect(page.locator('div[x-show="isMenuOpen"]')).toBeVisible()

  // Test mobile menu projects link with debug output - mobile menu has different selectors
  const mobileProjectsLink = page.locator('div[x-show="isMenuOpen"] a[hx-get="/projects"]')
  try {
    await expect(mobileProjectsLink).toContainText('Projects')
  } catch (error) {
    const actualText = await mobileProjectsLink.textContent()
    console.log(`Mobile projects link selector 'div[x-show="isMenuOpen"] a[hx-get="/projects"]' contains: "${actualText}"`)
    throw error
  }
})
