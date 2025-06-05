import { beforeAll, afterAll, beforeEach } from 'vitest'

// Global test setup
beforeAll(async () => {
  // Setup code that runs once before all tests
  console.log('Starting test suite...')
})

afterAll(async () => {
  // Cleanup code that runs once after all tests
  console.log('Test suite completed.')
})

beforeEach(async () => {
  // Reset any global state before each test
  // Clear localStorage, sessionStorage, etc. if needed
})

// Test utilities
export const waitForElement = async (selector: string, timeout = 5000) => {
  const start = Date.now()
  while (Date.now() - start < timeout) {
    const element = document.querySelector(selector)
    if (element) return element
    await new Promise(resolve => setTimeout(resolve, 100))
  }
  throw new Error(`Element ${selector} not found within ${timeout}ms`)
}

export const mockFetch = (responses: Record<string, any>) => {
  window.fetch = async (url: string | URL | Request) => {
    const urlString = typeof url === 'string' ? url : url.toString()
    const response = responses[urlString]
    
    if (!response) {
      throw new Error(`No mock response for ${urlString}`)
    }
    
    return new Response(JSON.stringify(response), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    })
  }
}

// Custom matchers
export const customMatchers = {
  toBeAccessible: async (element: Element) => {
    // Basic accessibility checks
    const hasAltText = element.tagName === 'IMG' ? element.hasAttribute('alt') : true
    const hasAriaLabel = element.hasAttribute('aria-label') || element.hasAttribute('aria-labelledby')
    const hasRole = element.hasAttribute('role')
    
    const pass = hasAltText && (hasAriaLabel || hasRole || element.textContent?.trim())
    
    return {
      pass,
      message: () => pass 
        ? 'Element is accessible' 
        : 'Element is not accessible - missing alt text, aria labels, or text content',
    }
  },
}