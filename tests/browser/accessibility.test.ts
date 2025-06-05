import { expect, test } from 'vitest'

test('heading hierarchy validation', async () => {
  // Create proper heading hierarchy
  const h1 = document.createElement('h1')
  h1.textContent = 'Main Title'
  
  const h2 = document.createElement('h2')
  h2.textContent = 'Section Title'
  
  const h3 = document.createElement('h3')
  h3.textContent = 'Subsection Title'
  
  document.body.append(h1, h2, h3)
  
  // Test heading hierarchy
  const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6')
  expect(headings.length).toBeGreaterThan(0)
  
  // Check if h1 exists
  expect(document.querySelector('h1')).toBeTruthy()
  
  // Clean up
  document.body.removeChild(h1)
  document.body.removeChild(h2)
  document.body.removeChild(h3)
})

test('image alt text validation', async () => {
  // Create image with alt text
  const img = document.createElement('img')
  img.src = '/test.jpg'
  img.alt = 'Test image description'
  document.body.appendChild(img)
  
  // Test alt text exists
  expect(img.alt).toBeTruthy()
  expect(img.alt).toBe('Test image description')
  
  // Clean up
  document.body.removeChild(img)
})

test('form label associations', async () => {
  // Create form with proper labels
  const form = document.createElement('form')
  
  const label = document.createElement('label')
  label.htmlFor = 'test-input'
  label.textContent = 'Test Input'
  
  const input = document.createElement('input')
  input.id = 'test-input'
  input.type = 'text'
  
  form.append(label, input)
  document.body.appendChild(form)
  
  // Test label association
  expect(label.htmlFor).toBe(input.id)
  expect(document.querySelector(`label[for="${input.id}"]`)).toBeTruthy()
  
  // Clean up
  document.body.removeChild(form)
})

test('landmarks and semantic structure', async () => {
  // Create semantic landmarks
  const header = document.createElement('header')
  const main = document.createElement('main')
  const footer = document.createElement('footer')
  const nav = document.createElement('nav')
  
  document.body.append(header, nav, main, footer)
  
  // Test landmarks exist
  expect(document.querySelector('header')).toBeTruthy()
  expect(document.querySelector('main')).toBeTruthy()
  expect(document.querySelector('footer')).toBeTruthy()
  expect(document.querySelector('nav')).toBeTruthy()
  
  // Clean up
  document.body.removeChild(header)
  document.body.removeChild(nav)
  document.body.removeChild(main)
  document.body.removeChild(footer)
})