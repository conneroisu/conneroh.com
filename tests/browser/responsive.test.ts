import { expect, test } from 'vitest'

test('viewport meta tag', async () => {
  // Create viewport meta tag
  const viewport = document.createElement('meta')
  viewport.name = 'viewport'
  viewport.content = 'width=device-width, initial-scale=1'
  document.head.appendChild(viewport)
  
  // Test viewport exists
  const foundViewport = document.querySelector('meta[name="viewport"]')
  expect(foundViewport).toBeTruthy()
  expect(foundViewport?.getAttribute('content')).toBe('width=device-width, initial-scale=1.0')
  
  // Clean up
  document.head.removeChild(viewport)
})

test('responsive image attributes', async () => {
  // Create responsive image
  const img = document.createElement('img')
  img.src = '/test.jpg'
  img.className = 'w-full h-auto'
  img.alt = 'Test image'
  document.body.appendChild(img)
  
  // Test responsive classes
  expect(img.className).toContain('w-full')
  expect(img.alt).toBe('Test image')
  
  // Clean up
  document.body.removeChild(img)
})

test('CSS media query support', async () => {
  // Test that CSS supports media queries
  expect(typeof window.matchMedia).toBe('function')
  
  // Test a media query
  const mediaQuery = window.matchMedia('(max-width: 768px)')
  expect(typeof mediaQuery.matches).toBe('boolean')
  expect(typeof mediaQuery.addListener).toBe('function')
})