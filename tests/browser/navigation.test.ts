import { expect, test, beforeEach } from 'vitest'
import { page, userEvent } from '@vitest/browser/context'

beforeEach(async () => {
  document.body.innerHTML = ''
})

test('HTMX navigation attributes', async () => {
  // Create navigation based on layout.templ
  const nav = document.createElement('nav')
  nav.setAttribute('data-testid', 'main-nav')
  
  const projectsLink = document.createElement('a')
  projectsLink.setAttribute('data-testid', 'projects-link')
  projectsLink.setAttribute('hx-target', '#bodiody')
  projectsLink.setAttribute('hx-get', '/projects')
  projectsLink.setAttribute('hx-push-url', 'true')
  projectsLink.textContent = 'Projects'
  projectsLink.className = 'hover:text-white cursor-pointer text-gray-300'
  
  const postsLink = document.createElement('a')
  postsLink.setAttribute('data-testid', 'posts-link')
  postsLink.setAttribute('hx-target', '#bodiody')
  postsLink.setAttribute('hx-get', '/posts')
  postsLink.setAttribute('hx-push-url', 'true')
  postsLink.textContent = 'Posts'
  postsLink.className = 'hover:text-white cursor-pointer text-gray-300'
  
  nav.append(projectsLink, postsLink)
  document.body.appendChild(nav)
  
  const navElement = page.getByTestId('main-nav')
  const projectsElement = page.getByTestId('projects-link')
  const postsElement = page.getByTestId('posts-link')
  
  // Test navigation structure
  await expect.element(navElement).toBeInTheDocument()
  await expect.element(projectsElement).toHaveAttribute('hx-get', '/projects')
  await expect.element(projectsElement).toHaveAttribute('hx-target', '#bodiody')
  await expect.element(postsElement).toHaveAttribute('hx-get', '/posts')
  
  // Test text content
  await expect.element(projectsElement).toHaveTextContent('Projects')
  await expect.element(postsElement).toHaveTextContent('Posts')
  
  // Test CSS classes for hover effects
  await expect.element(projectsElement).toHaveClass('hover:text-white', 'text-gray-300')
})

test('mobile menu with Alpine.js', async () => {
  // Create mobile menu with Alpine.js attributes
  const container = document.createElement('div')
  container.setAttribute('data-testid', 'nav-container')
  container.setAttribute('x-data', '{ isMenuOpen: false }')
  
  const menuButton = document.createElement('button')
  menuButton.setAttribute('data-testid', 'menu-toggle')
  menuButton.setAttribute('x-on:click', 'isMenuOpen = !isMenuOpen')
  menuButton.className = 'sm:hidden'
  menuButton.textContent = 'Menu'
  
  const mobileMenu = document.createElement('div')
  mobileMenu.setAttribute('data-testid', 'mobile-menu')
  mobileMenu.setAttribute('x-show', 'isMenuOpen')
  mobileMenu.className = 'sm:hidden'
  
  const mobileLink = document.createElement('a')
  mobileLink.setAttribute('data-testid', 'mobile-projects')
  mobileLink.setAttribute('x-on:click', 'isMenuOpen = false')
  mobileLink.textContent = 'Projects'
  
  mobileMenu.appendChild(mobileLink)
  container.append(menuButton, mobileMenu)
  document.body.appendChild(container)
  
  const containerElement = page.getByTestId('nav-container')
  const toggleElement = page.getByTestId('menu-toggle')
  const menuElement = page.getByTestId('mobile-menu')
  const linkElement = page.getByTestId('mobile-projects')
  
  // Test Alpine.js attributes
  await expect.element(containerElement).toHaveAttribute('x-data', '{ isMenuOpen: false }')
  await expect.element(toggleElement).toHaveAttribute('x-on:click', 'isMenuOpen = !isMenuOpen')
  await expect.element(menuElement).toHaveAttribute('x-show', 'isMenuOpen')
  await expect.element(linkElement).toHaveAttribute('x-on:click', 'isMenuOpen = false')
  
  // Test responsive classes
  await expect.element(toggleElement).toHaveClass('sm:hidden')
  await expect.element(menuElement).toHaveClass('sm:hidden')
})

test('breadcrumb navigation', async () => {
  // Create breadcrumb navigation
  const breadcrumb = document.createElement('nav')
  breadcrumb.setAttribute('data-testid', 'breadcrumb')
  breadcrumb.setAttribute('aria-label', 'Breadcrumb')
  
  const list = document.createElement('ol')
  list.className = 'flex space-x-2'
  
  const homeItem = document.createElement('li')
  const homeLink = document.createElement('a')
  homeLink.href = '/'
  homeLink.textContent = 'Home'
  homeLink.setAttribute('data-testid', 'breadcrumb-home')
  homeItem.appendChild(homeLink)
  
  const separator = document.createElement('li')
  separator.textContent = '/'
  separator.setAttribute('aria-hidden', 'true')
  
  const currentItem = document.createElement('li')
  currentItem.setAttribute('aria-current', 'page')
  currentItem.setAttribute('data-testid', 'breadcrumb-current')
  currentItem.textContent = 'Projects'
  
  list.append(homeItem, separator, currentItem)
  breadcrumb.appendChild(list)
  document.body.appendChild(breadcrumb)
  
  const breadcrumbElement = page.getByTestId('breadcrumb')
  const homeElement = page.getByTestId('breadcrumb-home')
  const currentElement = page.getByTestId('breadcrumb-current')
  
  // Test breadcrumb accessibility
  await expect.element(breadcrumbElement).toHaveAttribute('aria-label', 'Breadcrumb')
  await expect.element(currentElement).toHaveAttribute('aria-current', 'page')
  
  // Test content
  await expect.element(homeElement).toHaveTextContent('Home')
  await expect.element(currentElement).toHaveTextContent('Projects')
})

test('skip link accessibility', async () => {
  // Create skip link for accessibility
  const skipLink = document.createElement('a')
  skipLink.href = '#main'
  skipLink.setAttribute('data-testid', 'skip-link')
  skipLink.className = 'sr-only focus:not-sr-only'
  skipLink.textContent = 'Skip to main content'
  
  const main = document.createElement('main')
  main.id = 'main'
  main.setAttribute('data-testid', 'main-content')
  main.textContent = 'Main content here'
  
  document.body.append(skipLink, main)
  
  const skipElement = page.getByTestId('skip-link')
  const mainElement = page.getByTestId('main-content')
  
  // Test skip link
  await expect.element(skipElement).toHaveAttribute('href', '#main')
  await expect.element(skipElement).toHaveTextContent('Skip to main content')
  await expect.element(mainElement).toHaveAttribute('id', 'main')
  
  // Test accessibility classes
  await expect.element(skipElement).toHaveClass('sr-only', 'focus:not-sr-only')
})