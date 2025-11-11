# Implementation Tasks: Add Sanity App SDK

**Project ID**: z6sd4zry (conneroh.com)
**Content Types**: Post, Project, Tag, Employment
**Polling Strategy**: Poll-based (not real-time)
**Content Source**: Published documents only

## Overview
This is the ordered checklist of work items to integrate Sanity App SDK into the existing conneroh.com Sanity project. Complete each task in order and mark complete when done.

## Phase 0: Sanity Studio Schema Setup (Optional - if not already defined)

- [ ] **T0.1** Review existing Sanity schema
  - Check if Post, Project, Tag, Employment types already exist
  - Run: `bun sanity schema list` to see defined types
  - Review schema files in Sanity Studio if they exist
  - Validation: Document what schema already exists vs. needs creation

- [ ] **T0.2** Create schema types based on design document
  - Follow schema-design.md for Post, Project, Tag, Employment types
  - Create schema files in appropriate Sanity project directory
  - Define Portable Text blockContent type
  - Include validation rules and relationships
  - Validation: Schema files compile without errors

- [ ] **T0.3** Deploy schema to Sanity project
  - Run: `bun sanity deploy` to apply schema changes
  - Verify in Sanity Studio UI that content types appear
  - Validation: Can create new documents of each type in Studio

- [ ] **T0.4** Test schema with sample content
  - Create sample Post, Project, Tag, Employment documents
  - Test all relationships and references work
  - Validate slug uniqueness and required fields
  - Validation: Sample documents save without errors

## Phase 1: Installation & Configuration

- [ ] **T1.1** Install Sanity dependencies
  - Sanity CLI already installed as dev dependency (`sanity@4.15.0`)
  - Install `@sanity/client` package for website
  - Run `bun add @sanity/client` to add to dependencies
  - Verify package.json has both `sanity` (dev) and `@sanity/client`
  - Validation: `bun list | grep sanity` shows both packages

- [ ] **T1.2** Create Sanity environment configuration
  - Add `VITE_SANITY_PROJECT_ID=z6sd4zry` to `.env.local`
  - Add `VITE_SANITY_DATASET=production` to `.env.local`
  - Create API token in Sanity (project z6sd4zry): Project Settings → API → Add token
  - Add `SANITY_API_TOKEN=<token-here>` to `.env.local` (for server-side operations)
  - Create `.env.example` with template showing these variables
  - Validation: Environment variables accessible via `import.meta.env` and `process.env`

- [ ] **T1.3** Create Sanity client wrapper
  - Create `src/lib/sanity-client.ts`
  - Initialize `@sanity/client` with project ID and dataset
  - Export `createSanityClient()` factory function
  - Add TypeScript types for client
  - Validation: `import { createSanityClient } from '@/lib/sanity-client'` works without errors

- [ ] **T1.4** Add error handling utilities
  - Create `src/lib/sanity-errors.ts`
  - Define error types (ValidationError, NetworkError, AuthError)
  - Implement retry logic with exponential backoff
  - Add logging utilities
  - Validation: Error handling covers all major failure cases

## Phase 2: Content Fetching Primitives

- [ ] **T2.1** Create Solid.js content hooks
  - Create `src/hooks/use-sanity-content.ts`
  - Implement `createSanityQuery()` for reactive queries
  - Implement `createSanityDocument()` for single documents
  - Add Suspense compatibility
  - Validation: Hooks work with Solid.js `Show` and `Suspense` components

- [ ] **T2.2** Create server-side fetching utilities
  - Create `src/lib/sanity-server.ts`
  - Implement `fetchSanityContent()` for server-side data loading
  - Add support for route loaders in TanStack Start
  - Validation: Server function returns typed content without client overhead

- [ ] **T2.3** Create GROQ query builder utilities
  - Create `src/lib/sanity-queries.ts`
  - Implement helpers for common query patterns (list, bySlug, byId)
  - Add pagination utilities
  - Add query parameter validation
  - Validation: Queries are properly formatted GROQ strings

- [ ] **T2.4** Add query caching layer
  - Create `src/lib/sanity-cache.ts`
  - Implement memory cache for frequent queries
  - Add cache invalidation helpers
  - Set appropriate TTLs based on content type
  - Validation: Cache reduces redundant API calls in tests

## Phase 3: Type Safety with TypeGen

- [ ] **T3.1** Verify TypeGen is available
  - Sanity CLI already installed with typegen support
  - Verify: `bun sanity typegen --help`
  - Check for sanity.json or sanity.config.ts in project root
  - Validation: TypeGen command responds without errors

- [ ] **T3.2** Configure TypeGen for project z6sd4zry
  - Ensure Sanity config has projectId: 'z6sd4zry' and dataset: 'production'
  - Run: `bun sanity typegen generate` to generate types
  - Types generated from: Post, Project, Tag, Employment schemas
  - Validation: TypeGen successfully connects and generates types

- [ ] **T3.3** Generate TypeScript types from schema
  - Run: `bun sanity typegen generate`
  - Creates types for Post, Project, Tag, Employment, BlockContent
  - Output to: `src/lib/sanity-types.ts`
  - Add generated file to `.gitignore`
  - Validation: Generated file exists and contains all 4 document types

- [ ] **T3.4** Integrate generated types into client
  - Update `src/lib/sanity-client.ts` to import and use types
  - Add generic type parameters to fetch functions
  - Update `src/hooks/use-sanity-content.ts` with typed responses
  - Use Post, Project, Tag, Employment types in GROQ queries
  - Validation: TypeScript strict mode passes, no unused type errors

## Phase 4: Integration Tests

- [ ] **T4.1** Create Sanity client tests
  - Create `src/lib/sanity-client.spec.ts`
  - Test client initialization
  - Mock `@sanity/client` for unit tests
  - Test error handling and retries
  - Validation: `bun run test` passes for sanity-client tests

- [ ] **T4.2** Create content hook tests
  - Create `src/hooks/use-sanity-content.spec.ts`
  - Test query execution
  - Test Suspense integration
  - Test error states
  - Validation: Hooks properly handle async content loading

- [ ] **T4.3** Create integration tests with mock data
  - Create `tests/sanity-integration.spec.ts`
  - Test complete fetch → component render flow
  - Test SSR behavior
  - Validation: Full content pipeline works end-to-end

- [ ] **T4.4** Add E2E tests (optional phase 1)
  - Create Playwright test for content pages
  - Verify content renders correctly
  - Test error fallbacks
  - Validation: E2E tests pass with real Sanity connection

## Phase 5: Documentation & Examples

- [ ] **T5.1** Create setup documentation
  - Document environment variable setup
  - Add Sanity project creation guide
  - Include schema example
  - Create `docs/sanity-setup.md`
  - Validation: Documentation is complete and accurate

- [ ] **T5.2** Create usage examples
  - Example: Fetching posts list
  - Example: Single page/document
  - Example: Dynamic routing with slugs
  - Example: Real-time content sync
  - Validation: Examples are executable and correct

- [ ] **T5.3** Create API documentation
  - Document exported functions and hooks
  - Include TypeScript signatures
  - Add return type documentation
  - Validation: All public APIs are documented

- [ ] **T5.4** Update project README
  - Add Sanity section to main README
  - Link to Sanity setup documentation
  - Note any required environment variables
  - Validation: README is complete and discoverable

## Phase 6: Example Route Implementation

- [ ] **T6.1** Create post listing route
  - Create `src/routes/posts/index.tsx`
  - Implement route loader fetching: `*[_type == "post"] | order(createdAt desc)`
  - Display post list with title, excerpt, date, tags
  - Implement pagination for large lists
  - Validation: Route renders post list without errors

- [ ] **T6.2** Create dynamic post detail route
  - Create `src/routes/posts/[slug].tsx`
  - Implement loader fetching post by slug with: `*[_type == "post" && slug.current == $slug][0]`
  - Populate related posts, projects, employments
  - Render content using Portable Text renderer
  - Handle 404 for missing posts
  - Validation: Dynamic route works with real post slugs

- [ ] **T6.3** Create project listing and detail routes (optional)
  - Create `src/routes/projects/index.tsx` and `src/routes/projects/[slug].tsx`
  - Follow same pattern as posts
  - Include project-specific fields (tags, related posts)
  - Validation: Project routes functional

- [ ] **T6.4** Create tag listing routes (optional)
  - Create `src/routes/tags/index.tsx` and `src/routes/tags/[slug].tsx`
  - Show all posts/projects for given tag
  - Validation: Tag routes functional

- [ ] **T6.5** Create employment/experience route (optional)
  - Create `src/routes/experience/index.tsx` and `src/routes/experience/[slug].tsx`
  - Show employment history with start/end dates
  - Link to related projects, posts, skills
  - Validation: Experience routes functional

## Quality Gates

- [ ] All tests pass: `bun run test`
- [ ] No TypeScript errors: `tsc --noEmit`
- [ ] No ESLint issues: `bun run lint`
- [ ] Code formatted: `bun run format --write`
- [ ] Package.json has correct Sanity dependencies
- [ ] Environment variables documented
- [ ] Git commits are clean and descriptive

## Notes

- Each task should be completed fully before moving to the next
- Tests should be written/updated as you implement each feature
- Performance should be verified (check bundle size impact)
- Documentation should be kept up-to-date as implementation progresses
- Consider creating a feature branch: `feature/sanity-integration`
