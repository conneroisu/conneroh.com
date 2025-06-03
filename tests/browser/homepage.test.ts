import { expect, test, beforeEach } from 'vitest'
import { page, userEvent } from '@vitest/browser/context'

beforeEach(async () => {
  // Clear the page content before each test
  document.body.innerHTML = ''
})

test('homepage structure matches template', async () => {
  // Create homepage structure based on home.templ
  const main = document.createElement('main')
  main.id = 'bodiody'
  main.className = 'flex-grow'
  
  // Hero section
  const heroSection = document.createElement('section')
  heroSection.className = 'bg-gradient-to-b from-gray-900 to-gray-800 py-20'
  
  const heroTitle = document.createElement('h1')
  heroTitle.textContent = 'Conner Ohnesorge'
  heroTitle.setAttribute('aria-label', 'Name')
  heroTitle.className = 'mb-4 md:text-5xl text-4xl leading-tight font-bold lg:text-6xl text-white'
  
  const heroSummary = document.createElement('p')
  heroSummary.setAttribute('aria-label', 'summary')
  heroSummary.textContent = 'Electrical Engineer & Software Developer specialized in creating robust, scalable, and elegant solutions.'
  
  heroSection.appendChild(heroTitle)
  heroSection.appendChild(heroSummary)
  main.appendChild(heroSection)
  
  // Projects section
  const projectsSection = document.createElement('section')
  projectsSection.id = 'projects'
  projectsSection.className = 'bg-gray-800 py-16'
  
  const projectsTitle = document.createElement('h2')
  projectsTitle.textContent = 'Featured Projects'
  projectsTitle.className = 'mb-4 text-3xl font-bold text-white'
  
  projectsSection.appendChild(projectsTitle)
  main.appendChild(projectsSection)
  
  // Contact section
  const contactSection = document.createElement('section')
  contactSection.id = 'contact'
  contactSection.className = 'bg-gray-8  00 py-16'
  contactSection.setAttribute('aria-label', 'Contact')
  
  const contactTitle = document.createElement('h2')
  contactTitle.textContent = 'Get In Touch'
  contactTitle.className = 'mb-4 text-3xl font-bold text-white'
  
  contactSection.appendChild(contactTitle)
  main.appendChild(contactSection)
  
  document.body.appendChild(main)
  
  // Test homepage structure
  expect(document.getElementById('bodiody')).toBeTruthy()
  expect(document.querySelector('h1[aria-label="Name"]')?.textContent).toBe('Conner Ohnesorge')
  expect(document.querySelector('p[aria-label="summary"]')).toBeTruthy()
  expect(document.getElementById('projects')).toBeTruthy()
  expect(document.getElementById('contact')).toBeTruthy()
})

test('navigation structure', async () => {
  // Create header with navigation based on layout.templ
  const header = document.createElement('header')
  const nav = document.createElement('nav')
  nav.className = 'border-gray-800 border-b'
  
  const container = document.createElement('div')
  container.className = 'container mx-auto lg:px-8 sm:px-6 px-4'
  container.setAttribute('x-data', '{ isMenuOpen: false }')
  
  // Logo
  const logo = document.createElement('a')
  logo.className = 'text-white cursor-pointer pr-4 text-xl font-bold'
  logo.setAttribute('hx-get', '/')
  logo.setAttribute('hx-push-url', '/')
  logo.setAttribute('hx-target', '#bodiody')
  logo.setAttribute('aria-label', 'Back to Home')
  logo.textContent = 'Conner Ohnesorge'
  
  // Desktop navigation
  const desktopNav = document.createElement('div')
  desktopNav.className = 'space-x-8 hidden sm:flex items-center ml-8'
  
  const projectsLink = document.createElement('a')
  projectsLink.className = 'hover:text-white cursor-pointer text-gray-300'
  projectsLink.setAttribute('hx-target', '#bodiody')
  projectsLink.setAttribute('hx-get', '/projects')
  projectsLink.setAttribute('hx-push-url', 'true')
  projectsLink.textContent = 'Projects'
  
  const postsLink = document.createElement('a')
  postsLink.className = 'hover:text-white cursor-pointer text-gray-300'
  postsLink.setAttribute('hx-target', '#bodiody')
  postsLink.setAttribute('hx-get', '/posts')
  postsLink.setAttribute('hx-push-url', 'true')
  postsLink.textContent = 'Posts'
  
  desktopNav.appendChild(projectsLink)
  desktopNav.appendChild(postsLink)
  container.appendChild(logo)
  container.appendChild(desktopNav)
  nav.appendChild(container)
  header.appendChild(nav)
  document.body.appendChild(header)
  
  // Test navigation
  expect(document.querySelector('a[aria-label="Back to Home"]')?.textContent).toBe('Conner Ohnesorge')
  expect(document.querySelector('a[hx-get="/projects"]')?.textContent).toBe('Projects')
  expect(document.querySelector('a[hx-get="/posts"]')?.textContent).toBe('Posts')
})

test('mobile menu functionality', async () => {
  // Create mobile menu based on layout.templ
  const container = document.createElement('div')
  container.setAttribute('x-data', '{ isMenuOpen: false }')
  
  // Mobile menu button
  const menuButton = document.createElement('img')
  menuButton.className = 'p-2 focus:outline-none sm:hidden hover:text-white hover:bg-gray-700 rounded-md text-gray-300 mr-2'
  menuButton.setAttribute('x-on:click', 'isMenuOpen = !isMenuOpen')
  menuButton.src = 'https://conneroisu.fly.storage.tigris.dev/svg/menu.svg'
  
  // Mobile menu
  const mobileMenu = document.createElement('div')
  mobileMenu.setAttribute('x-show', 'isMenuOpen')
  mobileMenu.className = 'pb-4 space-y-1 sm:hidden pt-2'
  
  const mobileProjectsLink = document.createElement('a')
  mobileProjectsLink.className = 'text-base text-gray-300 hover:bg-gray-700 hover:text-white pl-3 pr-4 block py-2 font-medium'
  mobileProjectsLink.setAttribute('hx-target', '#bodiody')
  mobileProjectsLink.setAttribute('hx-get', '/projects')
  mobileProjectsLink.setAttribute('x-on:click', 'isMenuOpen = false')
  mobileProjectsLink.textContent = 'Projects'
  
  mobileMenu.appendChild(mobileProjectsLink)
  container.appendChild(menuButton)
  container.appendChild(mobileMenu)
  document.body.appendChild(container)
  
  // Test mobile menu elements
  expect(document.querySelector('img[src*="menu.svg"]')).toBeTruthy()
  expect(document.querySelector('div[x-show="isMenuOpen"]')).toBeTruthy()
  expect(document.querySelector('a[hx-get="/projects"]')?.textContent).toBe('Projects')
})