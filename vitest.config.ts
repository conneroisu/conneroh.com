import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    setupFiles: ['./tests/setup.ts'],
    // Include test files
    include: [
      'tests/**/*.{test,spec}.{js,ts,jsx,tsx}',
    ],
    // Exclude patterns
    exclude: [
      'node_modules/**',
      'dist/**',
      'build/**',
      '.direnv/**',
    ],
    // Browser configuration for browser tests
    browser: {
      provider: 'playwright',
      enabled: true,
      headless: true,
      name: 'chromium',
      screenshotFailures: true,
      viewport: {
        width: 1280,
        height: 720,
      },
    },
    // Reporter configuration
    reporters: process.env.CI ? ['verbose'] : ['default'],
    // Coverage configuration
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/**',
        'tests/**',
        '**/*.d.ts',
        '**/*.config.*',
        '**/mockData.ts',
        'cmd/conneroh/_static/**',
      ],
    },
  },
})
