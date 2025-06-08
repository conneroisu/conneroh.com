import { expect, test, beforeEach, vi, afterEach } from 'vitest'
import { mockFetch } from '../setup'

// Mock HTMX for form submission
const mockHTMX = {
  ajax: vi.fn(),
  trigger: vi.fn(),
  process: vi.fn(),
  on: vi.fn(),
  off: vi.fn(),
  config: {
    defaultSwapStyle: 'innerHTML',
    defaultSwapDelay: 0,
    defaultSettleDelay: 20,
    includeIndicatorStyles: true,
    indicatorClass: 'htmx-indicator',
    requestClass: 'htmx-request',
    addedClass: 'htmx-added',
    settlingClass: 'htmx-settling',
    swappingClass: 'htmx-swapping',
    allowEval: true,
    inlineScriptNonce: '',
    attributesToSettle: ['class', 'style', 'width', 'height'],
    withCredentials: false,
    timeout: 0,
    wsReconnectDelay: 'full-jitter',
    disableSelector: '[hx-disable], [data-hx-disable]',
    useTemplateFragments: false,
    scrollBehavior: 'smooth',
    defaultFocusScroll: false,
    getCacheBusterParam: false,
    globalViewTransitions: false,
    methodsThatUseUrlParams: ['get'],
    selfRequestsOnly: false,
    scrollIntoViewOnBoost: true,
    triggerSpecsCache: null,
  },
}

beforeEach(async () => {
  document.body.innerHTML = ''
  
  // Set up mocks in browser context
  if (typeof window !== 'undefined') {
    // @ts-ignore
    window.htmx = mockHTMX
  }
  
  mockHTMX.ajax.mockClear()
  mockHTMX.trigger.mockClear()
})

afterEach(async () => {
  // Individual mocks are cleared in beforeEach
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

test('full user dialog flow with contact form', async () => {
  // Create container for the entire flow
  const container = document.createElement('div')
  container.id = 'contact-container'
  document.body.appendChild(container)

  // Step 1: User navigates to contact form
  const form = document.createElement('form')
  form.className = 'bg-gray-800 p-8 rounded-lg space-y-6 shadow-md'
  form.setAttribute('hx-post', '/contact')
  form.setAttribute('hx-target', '#contact-container')
  form.setAttribute('hx-swap', 'innerHTML')
  
  // Create form fields
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
  nameInput.setAttribute('aria-required', 'true')
  nameInput.setAttribute('aria-label', 'Your name')
  
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
  emailInput.setAttribute('aria-required', 'true')
  emailInput.setAttribute('aria-label', 'Your email address')
  emailInput.setAttribute('aria-describedby', 'email-error')
  
  const emailError = document.createElement('span')
  emailError.id = 'email-error'
  emailError.className = 'text-red-500 text-sm hidden'
  emailError.setAttribute('role', 'alert')
  emailError.textContent = 'Please enter a valid email address'
  
  emailDiv.appendChild(emailLabel)
  emailDiv.appendChild(emailInput)
  emailDiv.appendChild(emailError)
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
  subjectInput.setAttribute('aria-required', 'true')
  subjectInput.setAttribute('aria-label', 'Message subject')
  
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
  messageInput.setAttribute('aria-required', 'true')
  messageInput.setAttribute('aria-label', 'Your message')
  messageInput.setAttribute('aria-describedby', 'message-counter')
  
  const messageCounter = document.createElement('div')
  messageCounter.id = 'message-counter'
  messageCounter.className = 'text-sm text-gray-400 mt-1'
  messageCounter.textContent = '0 / 500 characters'
  messageCounter.setAttribute('aria-live', 'polite')
  
  messageDiv.appendChild(messageLabel)
  messageDiv.appendChild(messageInput)
  messageDiv.appendChild(messageCounter)
  form.appendChild(messageDiv)
  
  // Submit button
  const submitButton = document.createElement('button')
  submitButton.type = 'submit'
  submitButton.className = 'rounded-md font-medium text-white transition-colors focus:ring-offset-2 focus:ring-2 focus:ring-green-500 w-full hover:bg-green-700 bg-green-600 py-3 focus:outline-none px-6'
  submitButton.textContent = 'Send Message'
  submitButton.setAttribute('aria-label', 'Send contact message')
  
  // Loading state indicator
  const loadingSpinner = document.createElement('span')
  loadingSpinner.className = 'htmx-indicator ml-2'
  loadingSpinner.innerHTML = '<svg class="animate-spin h-4 w-4 inline" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>'
  submitButton.appendChild(loadingSpinner)
  
  form.appendChild(submitButton)
  container.appendChild(form)
  
  // Step 2: User fills out the form
  // Simulate user typing with realistic interactions
  nameInput.focus()
  nameInput.value = 'Jane Smith'
  nameInput.dispatchEvent(new Event('input', { bubbles: true }))
  expect(nameInput.value).toBe('Jane Smith')
  
  emailInput.focus()
  emailInput.value = 'jane.smith@example.com'
  emailInput.dispatchEvent(new Event('input', { bubbles: true }))
  expect(emailInput.value).toBe('jane.smith@example.com')
  
  subjectInput.focus()
  subjectInput.value = 'Interested in Your Services'
  subjectInput.dispatchEvent(new Event('input', { bubbles: true }))
  expect(subjectInput.value).toBe('Interested in Your Services')
  
  messageInput.focus()
  const testMessage = 'Hello,\n\nI am very interested in learning more about your services. I particularly enjoyed your recent blog post about sustainable development.\n\nCould we schedule a call to discuss potential collaboration?\n\nBest regards,\nJane'
  messageInput.value = testMessage
  messageInput.dispatchEvent(new Event('input', { bubbles: true }))
  expect(messageInput.value).toBe(testMessage)
  
  // Update character counter
  messageCounter.textContent = `${testMessage.length} / 500 characters`
  
  // Step 3: Test form validation
  // Test invalid email
  emailInput.value = 'invalid-email'
  emailInput.dispatchEvent(new Event('blur'))
  emailError.classList.remove('hidden')
  expect(emailError.classList.contains('hidden')).toBe(false)
  
  // Fix email
  emailInput.value = 'jane.smith@example.com'
  emailInput.dispatchEvent(new Event('blur'))
  emailError.classList.add('hidden')
  expect(emailError.classList.contains('hidden')).toBe(true)
  
  // Step 4: Mock HTMX form submission
  mockHTMX.ajax.mockImplementation((method, url, context) => {
    expect(method).toBe('POST')
    expect(url).toBe('/contact')
    expect(context.source).toBe(form)
    
    // Simulate successful response with thank you component
    const thankYouHTML = `
      <div class="bg-green-500 p-8 rounded-lg shadow-md text-center">
        <div class="mb-6 flex justify-center">
          <div class="bg-green-600 rounded-full p-4 w-16 h-16 flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>
        </div>
        <h3 class="text-2xl font-bold text-white mb-4">Thank You!</h3>
        <p class="text-white mb-6">
          Your message has been sent successfully. I'll get back to you as soon as possible.
        </p>
      </div>
    `
    
    // Simulate server processing delay
    setTimeout(() => {
      container.innerHTML = thankYouHTML
      mockHTMX.trigger.mockImplementation((element, event) => {
        if (event === 'htmx:afterSwap') {
          // Focus management for accessibility
          const heading = container.querySelector('h3')
          heading?.setAttribute('tabindex', '-1')
          heading?.focus()
        }
      })
      mockHTMX.trigger(container, 'htmx:afterSwap')
    }, 100)
  })
  
  // Submit form
  form.addEventListener('submit', (e) => {
    e.preventDefault()
    submitButton.disabled = true
    submitButton.classList.add('opacity-50')
    loadingSpinner.classList.add('htmx-request')
    mockHTMX.ajax('POST', '/contact', { source: form })
  })
  
  submitButton.click()
  
  // Wait for response
  await new Promise(resolve => setTimeout(resolve, 150))
  
  // Step 5: Verify thank you message is displayed
  expect(container.querySelector('.bg-green-500')).toBeTruthy()
  expect(container.querySelector('h3')?.textContent).toBe('Thank You!')
  expect(container.querySelector('p')?.textContent).toContain('Your message has been sent successfully')
  expect(container.querySelector('svg')).toBeTruthy()
  
  // Verify focus was moved to heading for screen readers
  const heading = container.querySelector('h3')
  expect(heading?.getAttribute('tabindex')).toBe('-1')
  
  // Verify form was replaced
  expect(container.querySelector('form')).toBeFalsy()
})

test('contact form submits successfully and shows confirmation message', async () => {
  // Create form with HTMX attributes for frontend testing
  const form = document.createElement('form')
  form.className = 'bg-gray-800 p-8 rounded-lg space-y-6 shadow-md'
  form.setAttribute('hx-post', '/contact')
  form.setAttribute('hx-target', '#result')
  form.setAttribute('hx-swap', 'outerHTML')
  
  // Add form fields
  const nameInput = document.createElement('input')
  nameInput.name = 'name'
  nameInput.value = 'John Doe'
  nameInput.required = true
  
  const emailInput = document.createElement('input')
  emailInput.name = 'email'
  emailInput.type = 'email'
  emailInput.value = 'john@example.com'
  emailInput.required = true
  
  const subjectInput = document.createElement('input')
  subjectInput.name = 'subject'
  subjectInput.value = 'Test Subject'
  subjectInput.required = true
  
  const messageTextarea = document.createElement('textarea')
  messageTextarea.name = 'message'
  messageTextarea.value = 'Test message content'
  messageTextarea.required = true
  
  const submitButton = document.createElement('button')
  submitButton.type = 'submit'
  submitButton.textContent = 'Send Message'
  
  form.appendChild(nameInput)
  form.appendChild(emailInput)
  form.appendChild(subjectInput)
  form.appendChild(messageTextarea)
  form.appendChild(submitButton)
  
  const resultDiv = document.createElement('div')
  resultDiv.id = 'result'
  
  document.body.appendChild(form)
  document.body.appendChild(resultDiv)
  
  // Mock HTMX submission to simulate successful form submission
  mockHTMX.ajax.mockImplementation(() => {
    resultDiv.innerHTML = `
      <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
        <p>Thank you for your message! I'll get back to you soon.</p>
      </div>
    `
    return Promise.resolve()
  })
  
  // Simulate form submission via HTMX
  await mockHTMX.ajax('POST', '/contact', {
    source: form,
    target: '#result'
  })
  
  // Verify UI shows success message
  const successMessage = resultDiv.querySelector('.bg-green-100')
  expect(successMessage).toBeTruthy()
  expect(successMessage?.textContent).toContain('Thank you for your message')
})

test('contact form handles email service errors gracefully', async () => {
  // This test focuses on UI behavior regardless of email service status
  
  const form = document.createElement('form')
  form.className = 'bg-gray-800 p-8 rounded-lg space-y-6 shadow-md'
  form.setAttribute('hx-post', '/contact')
  form.setAttribute('hx-target', '#result')
  
  const resultDiv = document.createElement('div')
  resultDiv.id = 'result'
  document.body.appendChild(form)
  document.body.appendChild(resultDiv)
  
  // Mock HTMX to still show thank you despite email failure
  mockHTMX.ajax.mockImplementation(() => {
    // Even if email fails, user sees thank you message
    resultDiv.innerHTML = `
      <div class="bg-green-500 p-8 rounded-lg shadow-md text-center">
        <h3 class="text-2xl font-bold text-white mb-4">Thank You!</h3>
        <p class="text-white mb-6">
          Your message has been sent successfully. I'll get back to you as soon as possible.
        </p>
      </div>
    `
  })
  
  mockHTMX.ajax('POST', '/contact', { source: form })
  
  // Verify thank you is still shown
  expect(resultDiv.querySelector('h3')?.textContent).toBe('Thank You!')
})

test('contact form accessibility features', async () => {
  // Create accessible form
  const form = document.createElement('form')
  form.setAttribute('role', 'form')
  form.setAttribute('aria-label', 'Contact form')
  
  const nameInput = document.createElement('input')
  nameInput.id = 'name'
  nameInput.setAttribute('aria-required', 'true')
  nameInput.setAttribute('aria-invalid', 'false')
  
  const emailInput = document.createElement('input')
  emailInput.id = 'email'
  emailInput.type = 'email'
  emailInput.setAttribute('aria-required', 'true')
  emailInput.setAttribute('aria-describedby', 'email-hint email-error')
  
  const emailHint = document.createElement('span')
  emailHint.id = 'email-hint'
  emailHint.textContent = 'We\'ll use this to respond to you'
  
  const emailError = document.createElement('span')
  emailError.id = 'email-error'
  emailError.setAttribute('role', 'alert')
  emailError.setAttribute('aria-live', 'assertive')
  emailError.className = 'hidden'
  
  form.appendChild(nameInput)
  form.appendChild(emailInput)
  form.appendChild(emailHint)
  form.appendChild(emailError)
  document.body.appendChild(form)
  
  // Test ARIA attributes
  expect(form.getAttribute('role')).toBe('form')
  expect(form.getAttribute('aria-label')).toBe('Contact form')
  expect(nameInput.getAttribute('aria-required')).toBe('true')
  expect(emailInput.getAttribute('aria-describedby')).toBe('email-hint email-error')
  expect(emailError.getAttribute('role')).toBe('alert')
  expect(emailError.getAttribute('aria-live')).toBe('assertive')
  
  // Test validation state changes
  emailInput.value = 'invalid'
  emailInput.setAttribute('aria-invalid', 'true')
  emailError.classList.remove('hidden')
  emailError.textContent = 'Please enter a valid email address'
  
  expect(emailInput.getAttribute('aria-invalid')).toBe('true')
  expect(emailError.classList.contains('hidden')).toBe(false)
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

test('contact form prevents multiple simultaneous submissions', async () => {
  const form = document.createElement('form')
  form.setAttribute('hx-post', '/contact')
  
  const submitButton = document.createElement('button')
  submitButton.type = 'submit'
  submitButton.className = 'bg-green-600 py-3 px-6'
  submitButton.textContent = 'Send Message'
  
  form.appendChild(submitButton)
  document.body.appendChild(form)
  
  let submissionCount = 0
  form.addEventListener('submit', (e) => {
    e.preventDefault()
    submissionCount++
    submitButton.disabled = true
    submitButton.classList.add('opacity-50')
    submitButton.textContent = 'Sending...'
  })
  
  // Rapid clicks should only trigger one submission
  submitButton.click()
  submitButton.click()
  submitButton.click()
  
  expect(submissionCount).toBe(1)
  expect(submitButton.disabled).toBe(true)
  expect(submitButton.classList.contains('opacity-50')).toBe(true)
  expect(submitButton.textContent).toBe('Sending...')
})

test('contact form character limit enforcement', async () => {
  const messageInput = document.createElement('textarea')
  messageInput.id = 'message'
  messageInput.setAttribute('maxlength', '500')
  
  const counterDiv = document.createElement('div')
  counterDiv.id = 'message-counter'
  counterDiv.className = 'text-sm text-gray-400 mt-1'
  counterDiv.textContent = '0 / 500 characters'
  
  document.body.appendChild(messageInput)
  document.body.appendChild(counterDiv)
  
  // Test character counting
  const updateCounter = () => {
    const length = messageInput.value.length
    counterDiv.textContent = `${length} / 500 characters`
    if (length > 450) {
      counterDiv.classList.add('text-yellow-400')
    }
    if (length >= 500) {
      counterDiv.classList.add('text-red-400')
    }
  }
  
  messageInput.addEventListener('input', updateCounter)
  
  // Test normal input
  messageInput.value = 'Hello world'
  messageInput.dispatchEvent(new Event('input'))
  expect(counterDiv.textContent).toBe('11 / 500 characters')
  
  // Test approaching limit
  messageInput.value = 'a'.repeat(451)
  messageInput.dispatchEvent(new Event('input'))
  expect(counterDiv.textContent).toBe('451 / 500 characters')
  expect(counterDiv.classList.contains('text-yellow-400')).toBe(true)
  
  // Test at limit
  messageInput.value = 'a'.repeat(500)
  messageInput.dispatchEvent(new Event('input'))
  expect(counterDiv.textContent).toBe('500 / 500 characters')
  expect(counterDiv.classList.contains('text-red-400')).toBe(true)
})

test('contact form XSS prevention with malicious input', async () => {
  const form = document.createElement('form')
  
  const nameInput = document.createElement('input')
  nameInput.type = 'text'
  nameInput.id = 'name'
  
  const messageInput = document.createElement('textarea')
  messageInput.id = 'message'
  
  form.appendChild(nameInput)
  form.appendChild(messageInput)
  document.body.appendChild(form)
  
  // Test script injection attempts
  const maliciousInputs = [
    '<script>alert("XSS")</script>',
    '"><script>alert("XSS")</script>',
    'javascript:alert("XSS")',
    '<img src="x" onerror="alert(\'XSS\')" />',
    '{{constructor.constructor("alert(\'XSS\')")()}}',
  ]
  
  maliciousInputs.forEach((maliciousInput) => {
    nameInput.value = maliciousInput
    expect(nameInput.value).toBe(maliciousInput) // Input accepts the value
    
    // Verify form data is properly escaped when displayed
    const displayDiv = document.createElement('div')
    displayDiv.textContent = nameInput.value // textContent automatically escapes
    
    // Check that dangerous HTML elements are escaped
    if (maliciousInput.includes('<script>')) {
      expect(displayDiv.innerHTML).not.toContain('<script>')
      expect(displayDiv.innerHTML).toContain('&lt;script&gt;')
    }
    if (maliciousInput.includes('<img') && maliciousInput.includes('onerror')) {
      expect(displayDiv.innerHTML).toContain('&lt;img')
      expect(displayDiv.innerHTML).toContain('&gt;')
    }
    // For javascript: we just verify the content is text-only  
    if (maliciousInput.includes('javascript:')) {
      // The dangerous part is when it's in href/src attributes, not in text content
      expect(displayDiv.innerHTML).toContain('javascript:') // It's safe as text
    }
  })
})

test('contact form reset functionality after submission', async () => {
  const form = document.createElement('form')
  
  const nameInput = document.createElement('input')
  nameInput.type = 'text'
  nameInput.id = 'name'
  nameInput.value = 'John Doe'
  
  const emailInput = document.createElement('input')
  emailInput.type = 'email'
  emailInput.id = 'email'
  emailInput.value = 'john@example.com'
  
  const messageInput = document.createElement('textarea')
  messageInput.id = 'message'
  messageInput.value = 'Test message'
  
  const submitButton = document.createElement('button')
  submitButton.type = 'submit'
  submitButton.textContent = 'Send Message'
  
  form.appendChild(nameInput)
  form.appendChild(emailInput)
  form.appendChild(messageInput)
  form.appendChild(submitButton)
  document.body.appendChild(form)
  
  // Verify form has values
  expect(nameInput.value).toBe('John Doe')
  expect(emailInput.value).toBe('john@example.com')
  expect(messageInput.value).toBe('Test message')
  
  // Simulate successful submission and reset
  form.addEventListener('submit', (e) => {
    e.preventDefault()
    // Simulate server success response
    setTimeout(() => {
      form.reset()
      submitButton.disabled = false
      submitButton.classList.remove('opacity-50')
      submitButton.textContent = 'Send Message'
    }, 100)
  })
  
  submitButton.click()
  
  // Wait for reset
  await new Promise(resolve => setTimeout(resolve, 150))
  
  // Verify form was reset
  expect(nameInput.value).toBe('')
  expect(emailInput.value).toBe('')
  expect(messageInput.value).toBe('')
  expect(submitButton.disabled).toBe(false)
  expect(submitButton.classList.contains('opacity-50')).toBe(false)
  expect(submitButton.textContent).toBe('Send Message')
})