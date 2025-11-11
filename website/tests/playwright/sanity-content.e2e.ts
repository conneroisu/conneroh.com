import { test, expect } from '@playwright/test';

/**
 * End-to-End tests for Sanity CMS content rendering
 *
 * These tests verify the complete content pipeline from Sanity CMS to browser:
 * 1. Content fetching from Sanity API
 * 2. Server-side rendering with TanStack Start
 * 3. Client-side hydration
 * 4. Content display in browser
 * 5. Error handling with fallbacks
 *
 * Prerequisites:
 * - Sanity Studio running with sample content
 * - Environment variables set (VITE_SANITY_PROJECT_ID, VITE_SANITY_DATASET)
 * - Dev server running on http://localhost:3000
 *
 * Note: These tests can be run against:
 * - Mock Sanity data (for CI/CD pipelines)
 * - Real Sanity project (for integration testing)
 */

test.describe('Sanity Content E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Set up any necessary test context
    // For now, we'll assume the dev server is running
  });

  /**
   * Test: Content Rendering
   * Verify that content from Sanity renders correctly in the browser
   */
  test.describe('Content Rendering', () => {
    test.skip('should render blog posts list from Sanity', async ({ page }) => {
      // Navigate to posts page (adjust URL based on your routes)
      await page.goto('/posts');

      // Wait for content to load
      await page.waitForLoadState('networkidle');

      // Check for posts container
      const postsContainer = page.locator('[data-testid="posts-list"]');
      await expect(postsContainer).toBeVisible();

      // Verify at least one post is rendered
      const posts = page.locator('[data-testid="post-item"]');
      await expect(posts.first()).toBeVisible();

      // Check post structure
      const firstPost = posts.first();
      await expect(firstPost.locator('h2')).toBeVisible(); // Title
      await expect(firstPost.locator('[data-testid="post-excerpt"]')).toBeVisible(); // Excerpt
    });

    test.skip('should render single post detail page', async ({ page }) => {
      // Navigate to a specific post (adjust slug based on your test data)
      await page.goto('/posts/test-post');

      // Wait for content to load
      await page.waitForLoadState('networkidle');

      // Verify post title
      const title = page.locator('h1');
      await expect(title).toBeVisible();
      await expect(title).toHaveText(/./); // Has some text

      // Verify post content
      const content = page.locator('[data-testid="post-content"]');
      await expect(content).toBeVisible();

      // Verify published date
      const date = page.locator('[data-testid="post-date"]');
      await expect(date).toBeVisible();
    });

    test.skip('should render project listings', async ({ page }) => {
      await page.goto('/projects');
      await page.waitForLoadState('networkidle');

      const projectsContainer = page.locator('[data-testid="projects-list"]');
      await expect(projectsContainer).toBeVisible();

      const projects = page.locator('[data-testid="project-item"]');
      await expect(projects.first()).toBeVisible();
    });

    test.skip('should render employment history', async ({ page }) => {
      await page.goto('/experience');
      await page.waitForLoadState('networkidle');

      const employmentContainer = page.locator('[data-testid="employment-list"]');
      await expect(employmentContainer).toBeVisible();

      const jobs = page.locator('[data-testid="employment-item"]');
      await expect(jobs.first()).toBeVisible();
    });
  });

  /**
   * Test: SSR and Hydration
   * Verify server-side rendering and client-side hydration
   */
  test.describe('SSR and Hydration', () => {
    test.skip('should have content in initial HTML (SSR)', async ({ page }) => {
      // Navigate to posts page
      const response = await page.goto('/posts');

      // Get the raw HTML response
      const html = await response?.text();

      // Verify content is in the initial HTML (not loaded via JS)
      expect(html).toContain('post'); // Should contain post-related content
    });

    test.skip('should hydrate correctly without flickering', async ({ page }) => {
      // Record console errors
      const consoleErrors: string[] = [];
      page.on('console', msg => {
        if (msg.type() === 'error') {
          consoleErrors.push(msg.text());
        }
      });

      await page.goto('/posts');
      await page.waitForLoadState('networkidle');

      // Check for hydration errors
      expect(consoleErrors).toEqual([]);

      // Verify content is still visible after hydration
      const postsContainer = page.locator('[data-testid="posts-list"]');
      await expect(postsContainer).toBeVisible();
    });
  });

  /**
   * Test: Error Handling
   * Verify error states and fallbacks work correctly
   */
  test.describe('Error Handling', () => {
    test.skip('should show error message for non-existent post', async ({ page }) => {
      await page.goto('/posts/non-existent-post-slug');
      await page.waitForLoadState('networkidle');

      // Check for 404 or error message
      const errorMessage = page.locator('[data-testid="error-message"]');
      await expect(errorMessage).toBeVisible();
      await expect(errorMessage).toContainText(/not found|404/i);
    });

    test.skip('should show loading state while fetching', async ({ page }) => {
      // Slow down network to see loading state
      await page.route('**/*', route => {
        setTimeout(() => route.continue(), 1000);
      });

      const navigationPromise = page.goto('/posts');

      // Check for loading indicator
      const loadingIndicator = page.locator('[data-testid="loading"]');
      await expect(loadingIndicator).toBeVisible();

      await navigationPromise;
    });

    test.skip('should handle network errors gracefully', async ({ page, context }) => {
      // Simulate offline mode
      await context.setOffline(true);

      await page.goto('/posts');
      await page.waitForLoadState('networkidle');

      // Check for error message
      const errorMessage = page.locator('[data-testid="error-message"]');
      await expect(errorMessage).toBeVisible();
    });
  });

  /**
   * Test: Navigation and Linking
   * Verify links and navigation between content pages
   */
  test.describe('Navigation', () => {
    test.skip('should navigate from posts list to post detail', async ({ page }) => {
      await page.goto('/posts');
      await page.waitForLoadState('networkidle');

      // Click on first post
      const firstPostLink = page.locator('[data-testid="post-link"]').first();
      await firstPostLink.click();

      // Wait for navigation
      await page.waitForLoadState('networkidle');

      // Verify we're on a post detail page
      const postTitle = page.locator('h1');
      await expect(postTitle).toBeVisible();
    });

    test.skip('should navigate between related posts', async ({ page }) => {
      await page.goto('/posts/test-post');
      await page.waitForLoadState('networkidle');

      // Check for related posts section
      const relatedPosts = page.locator('[data-testid="related-posts"]');
      await expect(relatedPosts).toBeVisible();

      // Click on related post
      const relatedPostLink = relatedPosts.locator('[data-testid="post-link"]').first();
      await relatedPostLink.click();

      await page.waitForLoadState('networkidle');

      // Verify navigation worked
      const newPostTitle = page.locator('h1');
      await expect(newPostTitle).toBeVisible();
    });

    test.skip('should filter posts by tag', async ({ page }) => {
      await page.goto('/tags/typescript');
      await page.waitForLoadState('networkidle');

      // Verify tag page shows filtered posts
      const tagTitle = page.locator('[data-testid="tag-title"]');
      await expect(tagTitle).toContainText('TypeScript');

      const posts = page.locator('[data-testid="post-item"]');
      await expect(posts.first()).toBeVisible();
    });
  });

  /**
   * Test: Performance
   * Verify performance metrics for content loading
   */
  test.describe('Performance', () => {
    test.skip('should load posts page within acceptable time', async ({ page }) => {
      const startTime = Date.now();

      await page.goto('/posts');
      await page.waitForLoadState('networkidle');

      const loadTime = Date.now() - startTime;

      // Verify page loads in under 3 seconds
      expect(loadTime).toBeLessThan(3000);
    });

    test.skip('should use CDN for content images', async ({ page }) => {
      await page.goto('/posts');
      await page.waitForLoadState('networkidle');

      // Check image sources use Sanity CDN
      const images = page.locator('img[src*="sanity"]');
      const count = await images.count();

      if (count > 0) {
        const firstImageSrc = await images.first().getAttribute('src');
        expect(firstImageSrc).toContain('cdn.sanity.io');
      }
    });
  });

  /**
   * Test: Accessibility
   * Verify content meets accessibility standards
   */
  test.describe('Accessibility', () => {
    test.skip('should have proper heading hierarchy', async ({ page }) => {
      await page.goto('/posts');
      await page.waitForLoadState('networkidle');

      // Check for h1
      const h1 = page.locator('h1');
      await expect(h1).toBeVisible();

      // Verify heading structure
      const headings = page.locator('h1, h2, h3, h4, h5, h6');
      const count = await headings.count();
      expect(count).toBeGreaterThan(0);
    });

    test.skip('should have alt text for images', async ({ page }) => {
      await page.goto('/posts');
      await page.waitForLoadState('networkidle');

      const images = page.locator('img');
      const count = await images.count();

      // Check all images have alt attribute
      for (let i = 0; i < count; i++) {
        const alt = await images.nth(i).getAttribute('alt');
        expect(alt).toBeTruthy();
      }
    });
  });
});

/**
 * Note: All tests are currently skipped (.skip) because they require:
 * 1. Routes to be implemented (/posts, /projects, /experience, etc.)
 * 2. Components with data-testid attributes for reliable testing
 * 3. Sample content in Sanity Studio
 * 4. Dev server running at http://localhost:3000
 *
 * To enable these tests:
 * 1. Remove .skip from test() calls
 * 2. Implement the required routes
 * 3. Add data-testid attributes to components
 * 4. Ensure Sanity has sample content
 *
 * Run with: bun run test:e2e or npx playwright test
 */
