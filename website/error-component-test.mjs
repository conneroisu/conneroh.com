import { chromium } from 'playwright';

async function testErrorComponent() {
  console.log('Starting error component test...\n');
  
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext({ viewport: { width: 1280, height: 720 } });
  const page = await context.newPage();
  
  const results = {
    passed: [],
    failed: [],
    screenshots: []
  };
  
  // Enable console logging
  const consoleMessages = [];
  page.on('console', msg => {
    const msgText = msg.text();
    consoleMessages.push({ type: msg.type(), text: msgText });
    console.log('Browser Console [' + msg.type() + ']:', msgText);
  });
  
  try {
    // 1. Navigate to invalid route
    console.log('Test 1: Navigating to invalid route...');
    await page.goto('http://localhost:3000/nonexistent-page', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);
    
    // Take screenshot of error page
    await page.screenshot({ path: '/tmp/error-page.png', fullPage: true });
    results.screenshots.push('/tmp/error-page.png');
    console.log('Screenshot saved: /tmp/error-page.png\n');
    
    // 2. Check error code is displayed
    console.log('Test 2: Checking error code display...');
    const errorCode = await page.locator('h1').first().textContent();
    if (errorCode && errorCode.includes('404')) {
      results.passed.push('Error code 404 is displayed prominently');
      console.log('PASS: Error code 404 is displayed\n');
    } else {
      results.failed.push('Error code not found or incorrect. Found: ' + errorCode);
      console.log('FAIL: Error code not displayed correctly\n');
    }
    
    // 3. Check error message
    console.log('Test 3: Checking error message...');
    const errorTitle = await page.locator('h2').first().textContent();
    const errorDesc = await page.locator('p').first().textContent();
    if (errorTitle && errorTitle.includes('Page not found')) {
      results.passed.push('User-friendly error title is displayed');
      console.log('PASS: Error title: "' + errorTitle + '"\n');
    } else {
      results.failed.push('Error title not found. Found: ' + errorTitle);
    }
    
    if (errorDesc && errorDesc.length > 0) {
      results.passed.push('Error description is displayed');
      console.log('PASS: Error description: "' + errorDesc + '"\n');
    } else {
      results.failed.push('Error description not found');
    }
    
    // 4. Check navigation buttons
    console.log('Test 4: Checking navigation buttons...');
    const goBackButton = await page.locator('button:has-text("Go Back")');
    const returnHomeButton = await page.locator('a:has-text("Return Home")');
    
    const goBackVisible = await goBackButton.isVisible();
    const returnHomeVisible = await returnHomeButton.isVisible();
    
    if (goBackVisible) {
      results.passed.push('"Go Back" button is visible');
      console.log('PASS: "Go Back" button is visible\n');
    } else {
      results.failed.push('"Go Back" button is not visible');
    }
    
    if (returnHomeVisible) {
      results.passed.push('"Return Home" button is visible');
      console.log('PASS: "Return Home" button is visible\n');
    } else {
      results.failed.push('"Return Home" button is not visible');
    }
    
    // 5. Check dark theme styling
    console.log('Test 5: Checking dark theme styling...');
    const bgColor = await page.locator('div.min-h-screen').first().evaluate(el => 
      window.getComputedStyle(el).backgroundColor
    );
    if (bgColor.includes('17, 24, 39') || bgColor.includes('rgb(17, 24, 39)')) {
      results.passed.push('Dark theme background is applied');
      console.log('PASS: Dark theme background detected: ' + bgColor + '\n');
    } else {
      results.failed.push('Dark theme not detected. Background color: ' + bgColor);
    }
    
    // 6. Click "Return Home" button
    console.log('Test 6: Testing "Return Home" button...');
    await returnHomeButton.click();
    await page.waitForTimeout(2000);
    
    const homeUrl = page.url();
    if (homeUrl === 'http://localhost:3000/' || homeUrl === 'http://localhost:3000') {
      results.passed.push('"Return Home" button navigates to home page');
      console.log('PASS: Navigated to home page: ' + homeUrl + '\n');
      
      // Take screenshot of home page
      await page.screenshot({ path: '/tmp/home-page.png', fullPage: true });
      results.screenshots.push('/tmp/home-page.png');
      console.log('Screenshot saved: /tmp/home-page.png\n');
    } else {
      results.failed.push('"Return Home" button did not navigate to home. Current URL: ' + homeUrl);
    }
    
    // 7. Navigate back to error page and test "Go Back"
    console.log('Test 7: Testing "Go Back" button...');
    await page.goto('http://localhost:3000/nonexistent-page', { waitUntil: 'networkidle' });
    await page.waitForTimeout(1000);
    
    const goBackBtn = await page.locator('button:has-text("Go Back")');
    await goBackBtn.click();
    await page.waitForTimeout(1000);
    
    const urlAfterGoBack = page.url();
    if (urlAfterGoBack !== 'http://localhost:3000/nonexistent-page') {
      results.passed.push('"Go Back" button functions correctly');
      console.log('PASS: "Go Back" button navigated away from error page to: ' + urlAfterGoBack + '\n');
    } else {
      results.failed.push('"Go Back" button did not navigate away from error page');
    }
    
    // 8. Check console for warnings about errorComponent
    console.log('Test 8: Checking browser console for errorComponent warnings...');
    const hasErrorComponentWarning = consoleMessages.some(msg => 
      msg.text.includes('errorComponent') || msg.text.includes('Error component')
    );
    
    if (!hasErrorComponentWarning) {
      results.passed.push('No errorComponent warnings in console');
      console.log('PASS: No errorComponent warnings found\n');
    } else {
      results.failed.push('Found errorComponent warnings in console');
      console.log('FAIL: errorComponent warnings found\n');
    }
    
  } catch (error) {
    results.failed.push('Test execution error: ' + error.message);
    console.error('Test error:', error);
  } finally {
    await browser.close();
  }
  
  // Print results
  console.log('\n' + '='.repeat(60));
  console.log('TEST RESULTS SUMMARY');
  console.log('='.repeat(60));
  console.log('Total Passed: ' + results.passed.length);
  console.log('Total Failed: ' + results.failed.length);
  console.log('\nPassed Tests:');
  results.passed.forEach(test => console.log('  PASS: ' + test));
  
  if (results.failed.length > 0) {
    console.log('\nFailed Tests:');
    results.failed.forEach(test => console.log('  FAIL: ' + test));
  }
  
  console.log('\nScreenshots saved:');
  results.screenshots.forEach(path => console.log('  - ' + path));
  console.log('='.repeat(60));
  
  process.exit(results.failed.length > 0 ? 1 : 0);
}

testErrorComponent();
