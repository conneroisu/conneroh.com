import { expect, test } from 'vitest'

test('HTMX attributes and events', async () => {
  // Create element with HTMX-like attributes
  const link = document.createElement('a')
  link.setAttribute('hx-get', '/projects')
  link.setAttribute('hx-target', '#content')
  link.textContent = 'Projects'
  document.body.appendChild(link)
  
  // Test HTMX attributes exist
  expect(link.getAttribute('hx-get')).toBe('/projects')
  expect(link.getAttribute('hx-target')).toBe('#content')
  
  // Simulate HTMX event
  const htmxEvent = new CustomEvent('htmx:afterRequest', {
    detail: { xhr: { status: 200 } }
  })
  
  let eventFired = false
  document.addEventListener('htmx:afterRequest', () => {
    eventFired = true
  })
  
  document.dispatchEvent(htmxEvent)
  expect(eventFired).toBe(true)
  
  // Clean up
  document.body.removeChild(link)
})

test('form with HTMX attributes', async () => {
  const form = document.createElement('form')
  form.setAttribute('hx-post', '/submit')
  form.setAttribute('hx-target', '#result')
  
  const input = document.createElement('input')
  input.type = 'text'
  input.name = 'test'
  input.value = 'test value'
  
  form.appendChild(input)
  document.body.appendChild(form)
  
  // Test form attributes
  expect(form.getAttribute('hx-post')).toBe('/submit')
  expect(form.getAttribute('hx-target')).toBe('#result')
  expect(input.value).toBe('test value')
  
  // Clean up
  document.body.removeChild(form)
})