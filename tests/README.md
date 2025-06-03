# Testing Guide

This project uses Vitest with Playwright for comprehensive testing including unit tests and browser-based integration tests.

## Setup

All testing dependencies are managed through Nix. Enter the development shell:

```bash
nix develop
```

Then install JavaScript dependencies:

```bash
bun install
```

## Running Tests

### All Tests
```bash
bun test
```

### Test with UI
```bash
bun test:ui
# or
test-ui  # Nix command
```

### Run Tests Once (CI Mode)
```bash
bun test:run
# or
test-ci  # Nix command
```

### Coverage Report
```bash
bun test:coverage
```

## Test Structure

```
tests/
├── browser/          # Browser integration tests
│   ├── homepage.test.ts
│   ├── projects.test.ts
│   ├── posts.test.ts
│   ├── htmx.test.ts
│   ├── alpine.test.ts
│   ├── responsive.test.ts
│   └── accessibility.test.ts
├── unit/            # Unit tests
│   └── example.test.ts
├── setup.ts         # Test setup and utilities
└── README.md        # This file
```

## Writing Tests

### Browser Tests

Browser tests use Playwright to test the application in real browsers:

```typescript
import { expect, test } from 'vitest'
import { page } from '@vitest/browser/context'

test('my browser test', async () => {
  await page.goto('http://localhost:8080')
  await expect.element(page.getByText('Hello')).toBeInTheDocument()
})
```

### Unit Tests

Standard Vitest tests for testing functions and components:

```typescript
import { describe, expect, it } from 'vitest'

describe('myFunction', () => {
  it('should work', () => {
    expect(myFunction()).toBe(true)
  })
})
```

## CI Integration

Tests run automatically on GitHub Actions for:
- Every push to main
- Every pull request

The CI workflow:
1. Sets up Nix environment
2. Installs dependencies
3. Generates necessary files
4. Runs unit tests
5. Starts the application
6. Runs browser tests
7. Uploads coverage reports

## Test Configuration

- **vitest.config.ts**: Main configuration for development
- **vitest.config.browser.ts**: Browser-specific configuration for CI

## Debugging

To debug tests:

1. Use the UI mode: `bun test:ui`
2. Add `test.only()` to run a single test
3. Use `page.pause()` in browser tests to pause execution
4. Check browser console for errors

## Best Practices

1. **Test IDs**: Use `data-testid` attributes for reliable element selection
2. **Async Operations**: Always await browser interactions
3. **Cleanup**: Tests should be independent and not affect each other
4. **Accessibility**: Include accessibility tests for all new features
5. **Coverage**: Aim for high coverage but prioritize meaningful tests