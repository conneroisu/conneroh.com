import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    setupFiles: ['./tests/setup.ts'],
    browser: {
      provider: 'playwright',
      enabled: true,
      headless: true,
      instances: [
        { browser: 'chromium' },
      ],
    },
    // Include test files
    include: [
      'tests/**/*.{test,spec}.{js,ts,jsx,tsx}',
      '**/*.browser.{test,spec}.{js,ts,jsx,tsx}',
    ],
    // Exclude patterns
    exclude: [
      'node_modules/**',
      'dist/**',
      'build/**',
      '.direnv/**',
    ],
    // Test environment configuration
    environment: 'happy-dom',
    // Reporter configuration
    reporters: process.env.CI ? ['verbose', 'json'] : ['default'],
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