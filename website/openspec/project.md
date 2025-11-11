# Project Context

## Purpose
A TanStack Start + Bun production-ready web application with intelligent static asset loading and optimized performance. This website leverages Solid.js for reactive UI, Tailwind CSS for styling, and Bun as the runtime. The project implements a high-performance production server with memory-efficient asset management and comprehensive testing infrastructure.

## Tech Stack
- **Runtime**: Bun (high-performance JavaScript runtime)
- **Framework**: TanStack Start (meta-framework for full-stack applications)
- **UI Library**: Solid.js (reactive, fine-grained reactivity framework)
- **Routing**: TanStack Solid Router
- **Styling**: Tailwind CSS v4 with Vite integration
- **Build Tool**: Vite 7
- **Testing**: Vitest, @testing-library (DOM and Solid.js)
- **Language**: TypeScript 5.9+ (strict mode)
- **Code Quality**: ESLint (TanStack config), Prettier
- **Dev Tools**: TanStack Devtools, Router devtools
- **Package Manager**: Bun

## Project Conventions

### Code Style
- **Formatter**: Prettier with following rules:
  - No semicolons (`semi: false`)
  - Single quotes (`singleQuote: true`)
  - Trailing commas for all arguments/items (`trailingComma: 'all'`)
- **Linter**: ESLint with TanStack configuration (@tanstack/eslint-config)
- **TypeScript**: Strict mode enabled with no unused locals/parameters
- **Module Resolution**: Bundler mode with ESNext target (ES2022)
- **Path Aliases**: `@/*` maps to `./src/*`
- **JSX**: Solid.js via `jsxImportSource: 'solid-js'`

### Architecture Patterns
- **Full-Stack Framework**: TanStack Start for seamless server/client code sharing
- **Reactive UI**: Solid.js fine-grained reactivity for optimal performance
- **Server-Side Rendering**: SSR enabled via TanStack Start + Solid plugin
- **Asset Strategy**: Intelligent hybrid loading
  - Small files (<5MB default): Preloaded into memory at startup
  - Large files: Served on-demand from disk
  - Configurable include/exclude patterns via environment variables
- **Dev Tools Integration**: TanStack devtools for debugging UI and routing state
- **Vite Configuration**: Path aliases, Solid plugin SSR, Tailwind CSS Vite integration

### Testing Strategy
- **Framework**: Vitest for unit and integration tests
- **DOM Testing**: @testing-library/dom for DOM-level testing
- **Solid.js Testing**: @solidjs/testing-library for component testing
- **Test Location**: Collocated with source (e.g., `*.spec.ts`, `*.spec.tsx`)
- **Running Tests**: `bun run test`
- **Coverage Requirement**: Tests should verify core functionality and edge cases

### Git Workflow
- **Main Branch**: `main` (production-ready)
- **Current Development Branch**: `latest-mcs` (current working branch)
- **Feature Branches**: Descriptive names (e.g., `feature/component-name`, `codex/feature-description`)
- **Commit Style**: Descriptive messages
  - Format: `action: description` or `[topic] description`
  - Examples: "Enhance htmx-indicator styles", "added openspec inline-bun-nix-generation"
  - Be clear about what changed and why
- **Branching Strategy**: Feature branches off `main`, merged via pull requests

## Domain Context
- **Production Server**: Uses Bun with intelligent asset preloading for memory efficiency
- **Environment Variables** (server configuration):
  - `PORT`: Server port (default: 3000)
  - `STATIC_PRELOAD_MAX_BYTES`: Max file size for in-memory preload (default: 5MB)
  - `STATIC_PRELOAD_INCLUDE`: Comma-separated patterns to include (e.g., `*.js,*.css,*.woff2`)
  - `STATIC_PRELOAD_EXCLUDE`: Comma-separated patterns to exclude (e.g., `*.map,*.txt`)
  - `STATIC_PRELOAD_VERBOSE`: Enable detailed logging (true/false)
- **Development Workflow**:
  - Dev server: `bun run dev` (Vite dev server on port 3000)
  - Production build: `bun run build` (outputs to `dist/`)
  - Production run: `bun run start` (runs `server.ts`)
- **Performance Focus**: Vite-like output, cache headers, Lighthouse optimization

## Important Constraints
- **TypeScript Strict Mode**: All code must pass strict type checking
- **No Unused Code**: `noUnusedLocals` and `noUnusedParameters` enforced
- **Side Effect Imports**: `noUncheckedSideEffectImports` enabled
- **File Size Optimization**: Server intelligently manages in-memory preloading vs. on-demand serving
- **Solid.js Reactivity**: Components must leverage fine-grained reactivity (not React patterns)
- **SSR Compatibility**: All components must work in both server and client contexts

## External Dependencies
- **TanStack Ecosystem**: Solid Router, TanStack Start, Devtools
- **Testing Libraries**: Vitest, @testing-library
- **Styling**: Tailwind CSS (Vite integration)
- **Development**: Prettier, ESLint with TanStack config
- **Runtime**: Bun for both development and production
