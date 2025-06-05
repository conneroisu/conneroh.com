import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    browser: {
      provider: 'playwright',
      enabled: true,
      headless: true,
      instances: [
        { browser: 'chromium' },
        { browser: 'firefox' },
        { browser: 'webkit' },
      ],
    },
    // Only run browser tests
    include: [
      'tests/browser/**/*.{test,spec}.{js,ts,jsx,tsx}',
      '**/*.browser.{test,spec}.{js,ts,jsx,tsx}',
    ],
    // Exclude patterns
    exclude: [
      'node_modules/**',
      'dist/**',
      'build/**',
      '.direnv/**',
      'tests/unit/**',
    ],
    // Reporter configuration for CI
    reporters: ['verbose', 'json'],
    // Timeout for browser tests
    testTimeout: 30000,
  },
})