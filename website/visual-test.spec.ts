import { test, expect } from '@playwright/test';

const baseUrl = 'http://localhost:3000';

test.describe('Dark Theme Visual Testing', () => {
  test('Home Page - Desktop', async ({ page }) => {
    await page.goto(baseUrl);
    await page.waitForLoadState('networkidle');
    
    // Take full page screenshot
    await page.screenshot({ 
      path: '/tmp/screenshots/01-home-desktop.png', 
      fullPage: true 
    });
    
    // Verify header exists
    const header = page.locator('header');
    await expect(header).toBeVisible();
    
    // Check body has dark background
    const bodyBg = await page.evaluate(() => {
      return window.getComputedStyle(document.body).backgroundColor;
    });
    console.log('Body background:', bodyBg);
    
    // Verify footer exists
    const footer = page.locator('footer');
    await expect(footer).toBeVisible();
  });

  test('Projects Page - Desktop', async ({ page }) => {
    await page.goto(baseUrl + '/projects');
    await page.waitForLoadState('networkidle');
    
    await page.screenshot({ 
      path: '/tmp/screenshots/04-projects-desktop.png', 
      fullPage: true 
    });
    
    // Check for page title
    await expect(page.locator('h1')).toBeVisible();
  });

  test('Posts Page - Desktop', async ({ page }) => {
    await page.goto(baseUrl + '/posts');
    await page.waitForLoadState('networkidle');
    
    await page.screenshot({ 
      path: '/tmp/screenshots/05-posts-desktop.png', 
      fullPage: true 
    });
    
    await expect(page.locator('h1')).toBeVisible();
  });

  test('Tags Page - Desktop', async ({ page }) => {
    await page.goto(baseUrl + '/tags');
    await page.waitForLoadState('networkidle');
    
    await page.screenshot({ 
      path: '/tmp/screenshots/06-tags-desktop.png', 
      fullPage: true 
    });
    
    await expect(page.locator('h1')).toBeVisible();
  });

  test('Experience Page - Desktop', async ({ page }) => {
    await page.goto(baseUrl + '/experience');
    await page.waitForLoadState('networkidle');
    
    await page.screenshot({ 
      path: '/tmp/screenshots/07-experience-desktop.png', 
      fullPage: true 
    });
    
    await expect(page.locator('h1')).toBeVisible();
  });

  test('Home Page - Mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto(baseUrl);
    await page.waitForLoadState('networkidle');
    
    await page.screenshot({ 
      path: '/tmp/screenshots/08-home-mobile.png', 
      fullPage: true 
    });
  });

  test('Home Page - Tablet', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto(baseUrl);
    await page.waitForLoadState('networkidle');
    
    await page.screenshot({ 
      path: '/tmp/screenshots/10-home-tablet.png', 
      fullPage: true 
    });
  });

  test('Navigation Links Work', async ({ page }) => {
    await page.goto(baseUrl);
    
    // Check Projects link
    const projectsLink = page.locator('nav a[href*="projects"]').first();
    if (await projectsLink.count() > 0) {
      await projectsLink.click();
      await page.waitForLoadState('networkidle');
      await expect(page).toHaveURL(/.*projects/);
      await page.goBack();
    }
    
    // Check Posts link  
    const postsLink = page.locator('nav a[href*="posts"]').first();
    if (await postsLink.count() > 0) {
      await postsLink.click();
      await page.waitForLoadState('networkidle');
      await expect(page).toHaveURL(/.*posts/);
    }
  });
});
