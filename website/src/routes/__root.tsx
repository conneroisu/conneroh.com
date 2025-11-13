import { TanStackDevtools } from '@tanstack/solid-devtools'
import { HeadContent, Scripts, createRootRoute, Link, ErrorComponentProps } from '@tanstack/solid-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/solid-router-devtools'

import { HydrationScript } from 'solid-js/web'
import { Show } from 'solid-js'
import Header from '../components/Header'
import Footer from '../components/Footer'

import appCss from '../styles.css?url'
import type { JSX } from 'solid-js'

export const Route = createRootRoute({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'TanStack Start Starter',
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: appCss,
      },
    ],
  }),

  shellComponent: RootDocument,
  errorComponent: ErrorComponent,
})

function ErrorComponent(props: ErrorComponentProps) {
  const isDevelopment = import.meta.env.DEV

  // Determine error type and message
  const getErrorInfo = () => {
    const error = props.error
    const statusCode = error?.status || 500

    let title = 'Something went wrong'
    let description = 'An unexpected error occurred while loading this page.'

    if (statusCode === 404) {
      title = 'Page not found'
      description = "The page you're looking for doesn't exist or has been moved."
    } else if (statusCode >= 500) {
      title = 'Server error'
      description = 'Our server encountered an error. Please try again later.'
    } else if (statusCode >= 400) {
      title = 'Request error'
      description = 'There was a problem with your request.'
    }

    return { statusCode, title, description }
  }

  const errorInfo = getErrorInfo()

  return (
    <div class="min-h-screen bg-gray-900 text-white flex items-center justify-center px-4">
      <div class="max-w-2xl w-full text-center">
        {/* Error Status Code */}
        <div class="mb-8">
          <h1 class="text-8xl font-bold text-green-400 mb-4">
            {errorInfo.statusCode}
          </h1>
          <h2 class="text-3xl font-semibold text-white mb-4">
            {errorInfo.title}
          </h2>
          <p class="text-lg text-gray-400">
            {errorInfo.description}
          </p>
        </div>

        {/* Error Details (Development Only) */}
        <Show when={isDevelopment && props.error}>
          <div class="mb-8 text-left bg-gray-800 border border-gray-700 rounded-lg p-6">
            <h3 class="text-lg font-semibold text-red-400 mb-3">
              Error Details (Development Mode)
            </h3>
            <div class="space-y-2">
              <div>
                <span class="text-gray-400 font-medium">Message:</span>
                <p class="text-white mt-1 font-mono text-sm break-words">
                  {props.error?.message || 'No error message available'}
                </p>
              </div>
              <Show when={props.error?.stack}>
                <div class="mt-4">
                  <span class="text-gray-400 font-medium">Stack Trace:</span>
                  <pre class="text-white mt-2 font-mono text-xs overflow-x-auto bg-gray-900 p-4 rounded border border-gray-700">
                    {props.error?.stack}
                  </pre>
                </div>
              </Show>
            </div>
          </div>
        </Show>

        {/* Action Buttons */}
        <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
          <button
            onClick={() => window.history.back()}
            class="w-full sm:w-auto px-6 py-3 bg-gray-800 hover:bg-gray-700 text-white rounded-lg transition-colors border border-gray-700 font-medium"
          >
            Go Back
          </button>
          <Link
            to="/"
            class="w-full sm:w-auto px-6 py-3 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors font-medium"
          >
            Return Home
          </Link>
        </div>

        {/* Additional Help */}
        <div class="mt-8 text-gray-500 text-sm">
          <p>
            If this problem persists, please{' '}
            <Link to="/" class="text-green-400 hover:text-green-300 underline">
              contact support
            </Link>
            .
          </p>
        </div>
      </div>
    </div>
  )
}

function RootDocument({ children }: { children: JSX.Element }) {
  return (
    <html lang="en">
      <head>
        <HydrationScript />
      </head>
      <body class="bg-gray-900 text-white flex flex-col min-h-screen">
        <HeadContent />
        <Header />
        <main class="flex-grow">
          {children}
        </main>
        <Footer />
        <TanStackDevtools
          config={{
            position: 'bottom-left',
          }}
          plugins={[
            {
              name: 'Tanstack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
          ]}
        />
        <Scripts />
      </body>
    </html>
  )
}
