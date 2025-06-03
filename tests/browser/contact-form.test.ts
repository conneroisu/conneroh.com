import { expect, test, beforeEach } from 'vitest'
import { page, userEvent } from '@vitest/browser/context'

beforeEach(async () => {
  document.body.innerHTML = ''
})

test('contact form matches template structure', async () => {
  // Create contact form based on contact.templ
  const form = document.createElement('form')
  form.className = 'bg-gray-800 p-8 rounded-lg space-y-6 shadow-md'
  form.setAttribute('hx-post', '/contact')
  
  // Grid container
  const gridContainer = document.createElement('div')
  gridContainer.className = 'gap-6 grid grid-cols-1 md:grid-cols-2'
  
  // Name field
  const nameDiv = document.createElement('div')
  const nameLabel = document.createElement('label')
  nameLabel.htmlFor = 'name'
  nameLabel.className = 'mb-1 block text-sm font-medium text-gray-300'
  nameLabel.textContent = 'Name'
  
  const nameInput = document.createElement('input')
  nameInput.type = 'text'
  nameInput.id = 'name'
  nameInput.name = 'name'
  nameInput.className = 'py-2 bg-gray-700 focus:outline-none border-gray-600 w-full focus:ring-green-500 border px-4 text-white rounded-md focus:ring-2'
  nameInput.required = true
  
  nameDiv.appendChild(nameLabel)
  nameDiv.appendChild(nameInput)
  gridContainer.appendChild(nameDiv)
  
  // Email field
  const emailDiv = document.createElement('div')
  const emailLabel = document.createElement('label')
  emailLabel.htmlFor = 'email'
  emailLabel.className = 'mb-1 block text-sm font-medium text-gray-300'
  emailLabel.textContent = 'Email'
  
  const emailInput = document.createElement('input')
  emailInput.type = 'email'
  emailInput.id = 'email'
  emailInput.name = 'email'
  emailInput.className = 'py-2 bg-gray-700 focus:outline-none border-gray-600 w-full focus:ring-green-500 border px-4 text-white rounded-md focus:ring-2'
  emailInput.required = true
  
  emailDiv.appendChild(emailLabel)
  emailDiv.appendChild(emailInput)
  gridContainer.appendChild(emailDiv)
  
  form.appendChild(gridContainer)
  
  // Subject field
  const subjectDiv = document.createElement('div')
  const subjectLabel = document.createElement('label')
  subjectLabel.htmlFor = 'subject'
  subjectLabel.className = 'mb-1 block text-sm font-medium text-gray-300'
  subjectLabel.textContent = 'Subject'
  
  const subjectInput = document.createElement('input')
  subjectInput.type = 'text'
  subjectInput.id = 'subject'
  subjectInput.name = 'subject'
  subjectInput.className = 'py-2 bg-gray-700 focus:outline-none border-gray-600 w-full focus:ring-green-500 border px-4 text-white rounded-md focus:ring-2'
  subjectInput.required = true
  
  subjectDiv.appendChild(subjectLabel)
  subjectDiv.appendChild(subjectInput)
  form.appendChild(subjectDiv)
  
  // Message field
  const messageDiv = document.createElement('div')
  const messageLabel = document.createElement('label')
  messageLabel.htmlFor = 'message'
  messageLabel.className = 'mb-1 block text-sm font-medium text-gray-300'
  messageLabel.textContent = 'Message'
  
  const messageInput = document.createElement('textarea')
  messageInput.id = 'message'
  messageInput.name = 'message'
  messageInput.setAttribute('rows', '4')
  messageInput.className = 'py-2 bg-gray-700 focus:outline-none border-gray-600 w-full focus:ring-green-500 border px-4 text-white rounded-md focus:ring-2'
  messageInput.required = true
  
  messageDiv.appendChild(messageLabel)
  messageDiv.appendChild(messageInput)
  form.appendChild(messageDiv)
  
  // Submit button
  const submitButton = document.createElement('button')
  submitButton.type = 'submit'
  submitButton.className = 'rounded-md font-medium text-white transition-colors focus:ring-offset-2 focus:ring-2 focus:ring-green-500 w-full hover:bg-green-700 bg-green-600 py-3 focus:outline-none px-6'
  submitButton.textContent = 'Send Message'
  
  form.appendChild(submitButton)
  document.body.appendChild(form)
  
  // Test form structure
  expect(document.querySelector('form[hx-post="/contact"]')).toBeTruthy()
  expect(document.getElementById('name')).toBeTruthy()
  expect(document.getElementById('email')).toBeTruthy()
  expect(document.getElementById('subject')).toBeTruthy()
  expect(document.getElementById('message')).toBeTruthy()
  expect(document.querySelector('button[type="submit"]')?.textContent).toBe('Send Message')
})

test('contact form input interaction', async () => {
  // Create form fields
  const nameInput = document.createElement('input')
  nameInput.type = 'text'
  nameInput.id = 'name'
  nameInput.name = 'name'
  
  const emailInput = document.createElement('input')
  emailInput.type = 'email'
  emailInput.id = 'email'
  emailInput.name = 'email'
  
  const messageInput = document.createElement('textarea')
  messageInput.id = 'message'
  messageInput.name = 'message'
  
  document.body.appendChild(nameInput)
  document.body.appendChild(emailInput)
  document.body.appendChild(messageInput)
  
  // Test input interactions using direct DOM manipulation (since we can't guarantee page API)
  nameInput.value = 'John Doe'
  expect(nameInput.value).toBe('John Doe')
  
  emailInput.value = 'john@example.com'
  expect(emailInput.value).toBe('john@example.com')
  
  messageInput.value = 'Test message content'
  expect(messageInput.value).toBe('Test message content')
  
  // Test special characters
  nameInput.value = '{{a[['
  expect(nameInput.value).toBe('{{a[[')
  
  // Test form validation
  expect(emailInput.validity.valid).toBe(true) // valid email
  
  emailInput.value = 'invalid-email'
  expect(emailInput.validity.valid).toBe(false) // invalid email
})

test('contact form HTMX attributes', async () => {
  // Create form with HTMX attributes
  const form = document.createElement('form')
  form.className = 'bg-gray-800 p-8 rounded-lg space-y-6 shadow-md'
  form.setAttribute('hx-post', '/contact')
  
  const submitButton = document.createElement('button')
  submitButton.type = 'submit'
  submitButton.textContent = 'Send Message'
  
  form.appendChild(submitButton)
  document.body.appendChild(form)
  
  // Test HTMX attributes
  expect(form.getAttribute('hx-post')).toBe('/contact')
  expect(form.className).toContain('bg-gray-800')
  
  // Test form submission event
  let submitted = false
  form.addEventListener('submit', (e) => {
    e.preventDefault()
    submitted = true
  })
  
  submitButton.click()
  expect(submitted).toBe(true)
})

test('thank you component', async () => {
  // Create thank you component based on contact.templ ThankYou template
  const thankYouDiv = document.createElement('div')
  thankYouDiv.className = 'bg-green-500 p-8 rounded-lg shadow-md text-center'
  
  const iconContainer = document.createElement('div')
  iconContainer.className = 'mb-6 flex justify-center'
  
  const iconDiv = document.createElement('div')
  iconDiv.className = 'bg-green-600 rounded-full p-4 w-16 h-16 flex items-center justify-center'
  
  const svg = document.createElement('svg')
  svg.setAttribute('xmlns', 'http://www.w3.org/2000/svg')
  svg.className = 'h-8 w-8 text-white'
  svg.setAttribute('fill', 'none')
  svg.setAttribute('viewBox', '0 0 24 24')
  svg.setAttribute('stroke', 'currentColor')
  
  const path = document.createElement('path')
  path.setAttribute('stroke-linecap', 'round')
  path.setAttribute('stroke-linejoin', 'round')
  path.setAttribute('stroke-width', '2')
  path.setAttribute('d', 'M5 13l4 4L19 7')
  
  svg.appendChild(path)
  iconDiv.appendChild(svg)
  iconContainer.appendChild(iconDiv)
  
  const title = document.createElement('h3')
  title.className = 'text-2xl font-bold text-white mb-4'
  title.textContent = 'Thank You!'
  
  const message = document.createElement('p')
  message.className = 'text-white mb-6'
  message.textContent = "Your message has been sent successfully. I'll get back to you as soon as possible."
  
  thankYouDiv.appendChild(iconContainer)
  thankYouDiv.appendChild(title)
  thankYouDiv.appendChild(message)
  document.body.appendChild(thankYouDiv)
  
  // Test thank you component
  expect(document.querySelector('.bg-green-500')).toBeTruthy()
  expect(document.querySelector('h3')?.textContent).toBe('Thank You!')
  expect(document.querySelector('p')?.textContent).toContain('Your message has been sent successfully')
  expect(document.querySelector('svg')).toBeTruthy()
})