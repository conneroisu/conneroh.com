import { chromium } from 'playwright';
import fs from 'fs';
import path from 'path';

const screenshotDir = '/tmp/screenshots';
if (!fs.existsSync(screenshotDir)) {
  fs.mkdirSync(screenshotDir, { recursive: true });
}

const baseUrl = 'http://localhost:3000';

async function runVisualTests() {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 }
  });
  const page = await context.newPage();

  const results = {
    passed: [],
    failed: [],
    warnings: []
  };

  console.log('Starting comprehensive visual testing...\n');

  try {
    // Test 1: Home Page
    console.log('Testing Home Page (/)...');
    await page.goto(baseUrl, { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '01-home-desktop.png'), fullPage: true });
    
    const bodyBg = await page.evaluate(() => {
      const body = document.body;
      return window.getComputedStyle(body).backgroundColor;
    });
    console.log('  Body background: ' + bodyBg);
    
    const header = await page.locator('header').first();
    if (await header.count() > 0) {
      await header.screenshot({ path: path.join(screenshotDir, '02-header.png') });
      results.passed.push('Header exists and rendered');
    } else {
      results.failed.push('Header not found');
    }

    const navLinks = await page.locator('nav a').count();
    console.log('  Navigation links found: ' + navLinks);
    if (navLinks > 0) {
      results.passed.push('Navigation has ' + navLinks + ' links');
    } else {
      results.warnings.push('No navigation links found');
    }

    const footer = await page.locator('footer').first();
    if (await footer.count() > 0) {
      await footer.screenshot({ path: path.join(screenshotDir, '03-footer.png') });
      results.passed.push('Footer exists and rendered');
    } else {
      results.failed.push('Footer not found');
    }

    // Test 2: Projects Page
    console.log('\nTesting Projects Page (/projects)...');
    await page.goto(baseUrl + '/projects', { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '04-projects-desktop.png'), fullPage: true });
    
    const projectCards = await page.locator('[class*="card"]').count();
    console.log('  Project cards found: ' + projectCards);
    if (projectCards > 0) {
      results.passed.push('Projects page has ' + projectCards + ' cards');
    } else {
      results.warnings.push('No project cards found (might be empty state)');
    }

    // Test 3: Posts Page
    console.log('\nTesting Posts Page (/posts)...');
    await page.goto(baseUrl + '/posts', { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '05-posts-desktop.png'), fullPage: true });
    
    const postCards = await page.locator('[class*="card"]').count();
    console.log('  Post cards found: ' + postCards);
    if (postCards > 0) {
      results.passed.push('Posts page has ' + postCards + ' cards');
    } else {
      results.warnings.push('No post cards found (might be empty state)');
    }

    // Test 4: Tags Page
    console.log('\nTesting Tags Page (/tags)...');
    await page.goto(baseUrl + '/tags', { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '06-tags-desktop.png'), fullPage: true });
    
    const tagCards = await page.locator('[class*="card"]').count();
    console.log('  Tag cards found: ' + tagCards);
    if (tagCards > 0) {
      results.passed.push('Tags page has ' + tagCards + ' cards');
    } else {
      results.warnings.push('No tag cards found (might be empty state)');
    }

    // Test 5: Experience Page
    console.log('\nTesting Experience Page (/experience)...');
    await page.goto(baseUrl + '/experience', { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '07-experience-desktop.png'), fullPage: true });
    
    const experienceCards = await page.locator('[class*="card"]').count();
    console.log('  Experience cards found: ' + experienceCards);
    if (experienceCards > 0) {
      results.passed.push('Experience page has ' + experienceCards + ' cards');
    } else {
      results.warnings.push('No experience cards found (might be empty state)');
    }

    // Test 6: Mobile Responsiveness
    console.log('\nTesting Mobile Responsiveness...');
    await page.setViewportSize({ width: 375, height: 667 });
    
    await page.goto(baseUrl, { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '08-home-mobile.png'), fullPage: true });
    
    await page.goto(baseUrl + '/projects', { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '09-projects-mobile.png'), fullPage: true });
    
    results.passed.push('Mobile screenshots captured');

    // Test 7: Tablet Responsiveness
    console.log('\nTesting Tablet Responsiveness...');
    await page.setViewportSize({ width: 768, height: 1024 });
    
    await page.goto(baseUrl, { waitUntil: 'networkidle' });
    await page.screenshot({ path: path.join(screenshotDir, '10-home-tablet.png'), fullPage: true });
    
    results.passed.push('Tablet screenshots captured');

    // Test 8: Color Theme Verification
    console.log('\nVerifying Dark Theme Colors...');
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.goto(baseUrl, { waitUntil: 'networkidle' });
    
    const themeColors = await page.evaluate(() => {
      const results = {};
      
      const body = document.body;
      results.bodyBg = window.getComputedStyle(body).backgroundColor;
      
      const header = document.querySelector('header');
      if (header) {
        results.headerBg = window.getComputedStyle(header).backgroundColor;
      }
      
      const lightBgElements = [];
      const allElements = document.querySelectorAll('*');
      allElements.forEach(el => {
        const bg = window.getComputedStyle(el).backgroundColor;
        if (bg === 'rgb(255, 255, 255)' || bg === 'rgb(249, 250, 251)') {
          lightBgElements.push(el.tagName);
        }
      });
      results.lightBgCount = lightBgElements.length;
      results.lightBgElements = lightBgElements.slice(0, 5);
      
      return results;
    });
    
    console.log('  Body BG: ' + themeColors.bodyBg);
    console.log('  Header BG: ' + (themeColors.headerBg || 'N/A'));
    console.log('  Light backgrounds found: ' + themeColors.lightBgCount);
    
    if (themeColors.lightBgCount === 0) {
      results.passed.push('No light backgrounds detected - dark theme consistent');
    } else {
      results.warnings.push(themeColors.lightBgCount + ' light background elements found');
    }

  } catch (error) {
    console.error('\nTest execution error:', error);
    results.failed.push('Test execution error: ' + error.message);
  } finally {
    await browser.close();
  }

  console.log('\n' + '='.repeat(60));
  console.log('VISUAL TESTING RESULTS');
  console.log('='.repeat(60));
  
  console.log('\nPASSED (' + results.passed.length + '):');
  results.passed.forEach(item => console.log('  ✓ ' + item));
  
  if (results.warnings.length > 0) {
    console.log('\nWARNINGS (' + results.warnings.length + '):');
    results.warnings.forEach(item => console.log('  ⚠ ' + item));
  }
  
  if (results.failed.length > 0) {
    console.log('\nFAILED (' + results.failed.length + '):');
    results.failed.forEach(item => console.log('  ✗ ' + item));
  }
  
  console.log('\n' + '='.repeat(60));
  console.log('Screenshots saved to: ' + screenshotDir);
  console.log('='.repeat(60));
  
  const overallStatus = results.failed.length === 0 ? 'PASS' : 'FAIL';
  console.log('\nOVERALL STATUS: ' + overallStatus + '\n');
  
  fs.writeFileSync(
    path.join(screenshotDir, 'test-results.json'),
    JSON.stringify(results, null, 2)
  );
  
  return results;
}

runVisualTests().catch(console.error);
