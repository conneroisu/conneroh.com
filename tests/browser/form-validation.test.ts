import { expect, test, beforeEach } from 'vitest'
import { page, userEvent } from '@vitest/browser/context'

beforeEach(async () => {
  document.body.innerHTML = ''
})

test('form validation with expect.element', async () => {
  // Create a realistic contact form
  const form = document.createElement('form')
  form.setAttribute('data-testid', 'contact-form')
  
  const emailInput = document.createElement('input')
  emailInput.type = 'email'
  emailInput.setAttribute('data-testid', 'email-input')
  emailInput.required = true
  
  const submitButton = document.createElement('button')
  submitButton.type = 'submit'
  submitButton.setAttribute('data-testid', 'submit-btn')
  submitButton.textContent = 'Submit'
  
  form.append(emailInput, submitButton)
  document.body.appendChild(form)
  
  // Use proper Vitest browser API
  const formElement = page.getByTestId('contact-form')
  const emailElement = page.getByTestId('email-input')
  const submitElement = page.getByTestId('submit-btn')
  
  // Test initial state
  await expect.element(formElement).toBeInTheDocument()
  await expect.element(emailElement).toBeRequired()
  await expect.element(submitElement).toBeEnabled()
  
  // Test invalid email
  await userEvent.fill(emailElement, 'invalid-email')
  await expect.element(emailElement).toBeInvalid()
  
  // Test valid email
  await userEvent.fill(emailElement, 'test@example.com')
  await expect.element(emailElement).toBeValid()
  await expect.element(emailElement).toHaveValue('test@example.com')
})

test('form accessibility attributes', async () => {
  // Create form with proper accessibility
  const nameLabel = document.createElement('label')
  nameLabel.htmlFor = 'name-input'
  nameLabel.textContent = 'Full Name'
  
  const nameInput = document.createElement('input')
  nameInput.id = 'name-input'
  nameInput.type = 'text'
  nameInput.setAttribute('data-testid', 'name-field')
  nameInput.setAttribute('aria-required', 'true')
  nameInput.required = true
  
  const errorDiv = document.createElement('div')
  errorDiv.id = 'name-error'
  errorDiv.setAttribute('role', 'alert')
  errorDiv.textContent = 'Name is required'
  errorDiv.style.display = 'none'
  
  nameInput.setAttribute('aria-describedby', 'name-error')
  
  document.body.append(nameLabel, nameInput, errorDiv)
  
  const nameField = page.getByTestId('name-field')
  
  // Test accessibility attributes
  await expect.element(nameField).toBeRequired()
  await expect.element(nameField).toHaveAttribute('aria-required', 'true')
  await expect.element(nameField).toHaveAttribute('aria-describedby', 'name-error')
  
  // Test form validation triggers error
  nameInput.value = ''
  nameInput.dispatchEvent(new Event('blur'))
  
  // Show error
  errorDiv.style.display = 'block'
  nameInput.setAttribute('aria-invalid', 'true')
  
  await expect.element(nameField).toBeInvalid()
})

test('form styling and classes', async () => {
  // Create form with CSS classes
  const input = document.createElement('input')
  input.className = 'py-2 bg-gray-700 focus:outline-none border-gray-600 w-full focus:ring-green-500 border px-4 text-white rounded-md focus:ring-2'
  input.setAttribute('data-testid', 'styled-input')
  
  const button = document.createElement('button')
  button.className = 'rounded-md font-medium text-white transition-colors bg-green-600 py-3 px-6'
  button.setAttribute('data-testid', 'styled-button')
  button.style.backgroundColor = 'rgb(34, 197, 94)' // green-600
  
  document.body.append(input, button)
  
  const inputElement = page.getByTestId('styled-input')
  const buttonElement = page.getByTestId('styled-button')
  
  // Test CSS classes
  await expect.element(inputElement).toHaveClass('bg-gray-700', 'text-white', 'rounded-md')
  await expect.element(buttonElement).toHaveClass('font-medium', 'bg-green-600')
  
  // Test computed styles
  await expect.element(buttonElement).toHaveStyle({
    backgroundColor: 'rgb(34, 197, 94)'
  })
})

test('form interactions and events', async () => {
  // Create interactive form
  const textarea = document.createElement('textarea')
  textarea.setAttribute('data-testid', 'message-area')
  textarea.placeholder = 'Enter your message...'
  
  const charCount = document.createElement('div')
  charCount.setAttribute('data-testid', 'char-count')
  charCount.textContent = '0/500'
  
  // Add event listener that updates on input
  textarea.addEventListener('input', () => {
    charCount.textContent = `${textarea.value.length}/500`
  })
  
  document.body.append(textarea, charCount)
  
  const messageArea = page.getByTestId('message-area')
  const countDisplay = page.getByTestId('char-count')
  
  // Test initial state
  await expect.element(messageArea).toHaveValue('')
  await expect.element(countDisplay).toHaveTextContent('0/500')
  
  // Use userEvent.type which properly triggers input events
  await userEvent.type(messageArea, 'Hello, test!')
  
  await expect.element(messageArea).toHaveValue('Hello, test!')
  await expect.element(countDisplay).toHaveTextContent('12/500')
  
  // Test clearing using userEvent.clear
  await userEvent.clear(messageArea)
  
  await expect.element(messageArea).toHaveValue('')
  await expect.element(countDisplay).toHaveTextContent('0/500')
})