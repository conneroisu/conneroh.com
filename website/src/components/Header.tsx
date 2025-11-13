import { Link } from '@tanstack/solid-router'
import { createSignal } from 'solid-js'

export default function Header() {
  const [isMenuOpen, setIsMenuOpen] = createSignal(false)

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen())
  }

  const closeMenu = () => {
    setIsMenuOpen(false)
  }

  return (
    <header class="bg-gray-900">
      <nav class="border-b border-gray-800">
        <div class="container mx-auto px-4 lg:px-8">
          <div class="flex justify-between items-center h-16">
            {/* Left side: Logo + Desktop Nav */}
            <div class="flex items-center">
              {/* Mobile menu button */}
              <button
                class="sm:hidden text-gray-300 hover:text-white mr-4"
                onClick={toggleMenu}
                aria-label="Toggle menu"
              >
                <svg
                  class="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                  />
                </svg>
              </button>

              {/* Logo/Site Name */}
              <Link
                to="/"
                class="text-white text-xl font-bold hover:text-green-400 transition-colors"
              >
                Conner Ohnesorge
              </Link>

              {/* Desktop Navigation */}
              <div class="hidden sm:flex space-x-8 ml-8">
                <Link
                  to="/projects"
                  class="text-gray-300 hover:text-white transition-colors"
                >
                  Projects
                </Link>
                <Link
                  to="/posts"
                  class="text-gray-300 hover:text-white transition-colors"
                >
                  Posts
                </Link>
                <Link
                  to="/tags"
                  class="text-gray-300 hover:text-white transition-colors"
                >
                  Tags
                </Link>
                <Link
                  to="/experience"
                  class="text-gray-300 hover:text-white transition-colors"
                >
                  Experience
                </Link>
              </div>
            </div>
          </div>

          {/* Mobile Menu */}
          <div
            class={`sm:hidden overflow-hidden transition-all duration-300 ${
              isMenuOpen() ? 'max-h-64 opacity-100' : 'max-h-0 opacity-0'
            }`}
          >
            <div class="py-4 space-y-2">
              <Link
                to="/projects"
                class="block px-4 py-2 text-gray-300 hover:text-white hover:bg-gray-800 transition-colors rounded"
                onClick={closeMenu}
              >
                Projects
              </Link>
              <Link
                to="/posts"
                class="block px-4 py-2 text-gray-300 hover:text-white hover:bg-gray-800 transition-colors rounded"
                onClick={closeMenu}
              >
                Posts
              </Link>
              <Link
                to="/tags"
                class="block px-4 py-2 text-gray-300 hover:text-white hover:bg-gray-800 transition-colors rounded"
                onClick={closeMenu}
              >
                Tags
              </Link>
              <Link
                to="/experience"
                class="block px-4 py-2 text-gray-300 hover:text-white hover:bg-gray-800 transition-colors rounded"
                onClick={closeMenu}
              >
                Experience
              </Link>
            </div>
          </div>
        </div>
      </nav>
    </header>
  )
}
