# Specification: Sanity Client Integration

**Capability ID**: `sanity-client`

**Version**: 1.0.0

**Status**: Proposed

## Overview

The Sanity Client capability provides content management and fetching functionality for the TanStack Start + Solid.js website. It enables querying, caching, and real-time synchronization of content from a Sanity backend.

## ADDED Requirements

### Requirement: System SHALL initialize Sanity client with project configuration

The client SHALL be initialized with `projectId` and `dataset` from environment variables. The client SHALL support optional API token for authenticated requests. The client SHALL use lazy initialization, and the implementation SHALL use singleton pattern to avoid multiple client instances.

**Implementation**:
- `src/lib/sanity-client.ts` exports `createSanityClient()` factory
- Configuration from `VITE_SANITY_PROJECT_ID`, `VITE_SANITY_DATASET`, optional `SANITY_API_TOKEN`

#### Scenario: Basic Client Creation
```typescript
// Environment: VITE_SANITY_PROJECT_ID=abc123, VITE_SANITY_DATASET=production
const client = createSanityClient()
// Result: Client instance configured with project and dataset
```

#### Scenario: Client with Authentication Token
```typescript
// Environment: SANITY_API_TOKEN=abc123token (set server-side)
const client = createSanityClient({ token: process.env.SANITY_API_TOKEN })
// Result: Authenticated client for write operations
```

---

### Requirement: System SHALL execute GROQ queries against Sanity content

The system SHALL fetch documents using GROQ query language. The system SHALL support query parameters for parameterized queries. The system SHALL return strongly-typed results when types are available. The system SHALL provide proper error handling for invalid queries.

**Implementation**:
- `client.fetch(query, params)` from `@sanity/client`
- Wrapped with error handling in `src/lib/sanity-client.ts`

#### Scenario: Fetch All Posts
```typescript
const posts = await client.fetch(
  '*[_type == "post"] | order(_createdAt desc)',
  { limit: 10 }
)
// Result: Array of posts sorted by creation date, limited to 10
```

#### Scenario: Fetch with Parameters
```typescript
const post = await client.fetch(
  '*[_type == "post" && slug.current == $slug][0]',
  { slug: 'my-first-post' }
)
// Result: Single post matching slug parameter
```

---

### Requirement: System SHALL implement query caching to reduce API calls

The system SHALL implement in-memory cache for query results. The system SHALL provide configurable TTL per query type. The system SHALL provide cache invalidation helpers. The cache key SHALL be based on query + parameters.

**Implementation**:
- `src/lib/sanity-cache.ts` provides `getCachedQuery()` and `invalidateCache()`
- Integrated into fetch wrapper

#### Scenario: Repeated Query Uses Cache
```typescript
// First call - API request
const posts1 = await getCachedQuery('*[_type == "post"]')

// Second call within TTL - returns cached result
const posts2 = await getCachedQuery('*[_type == "post"]')

// Result: posts1 === posts2 (same reference), no second API call
```

#### Scenario: Cache Invalidation on Content Change
```typescript
// After content update webhook
await invalidateCache('*[_type == "post"]')

// Next query makes fresh API call
const posts = await getCachedQuery('*[_type == "post"]')
```

---

### Requirement: System SHALL handle errors and retry failed requests

The system SHALL retry failed requests with exponential backoff. The system SHALL limit retries to maximum 3 attempts. The system SHALL distinguish between recoverable and non-recoverable errors. The system SHALL log errors for debugging purposes.

**Implementation**:
- `src/lib/sanity-errors.ts` defines error types
- Retry logic in fetch wrapper

#### Scenario: Network Error with Retry
```typescript
// Network fails on first attempt, succeeds on retry
const result = await client.fetch(query)
// Result: Query succeeds after 1 retry with exponential backoff
```

#### Scenario: Invalid GROQ Query
```typescript
try {
  await client.fetch('INVALID GROQ')
} catch (error) {
  expect(error.message).toContain('GROQ syntax')
}
// Result: SanityValidationError thrown immediately (not retried)
```

---

### Requirement: System SHALL generate TypeScript types from Sanity schema

The system SHALL auto-generate types from Sanity content schema. The system SHALL include all document types and fields in generated types. The system SHALL export types for use in application code. The system SHALL support updating types when schema changes.

**Implementation**:
- `sanity typegen` command generates `src/lib/sanity-types.ts`
- Types used throughout application for type safety

#### Scenario: Type-Safe Query Results
```typescript
import type { Post } from '@/lib/sanity-types'

const post: Post = await client.fetch(
  '*[_type == "post" && slug.current == $slug][0]',
  { slug: 'hello' }
)

// TypeScript knows post.title is string, post.author is Author ref, etc.
post.title.toUpperCase() // ✓ Valid
post.unknownField // ✗ TypeScript error
```

---

### Requirement: System SHALL provide Solid.js reactive hooks for Sanity content

The system SHALL provide `createSanityQuery()` for fetching lists and queries. The system SHALL provide `createSanityDocument()` for fetching single documents. The hooks SHALL be Suspense-compatible for server-side rendering. The hooks SHALL support reactive updates when query parameters change.

**Implementation**:
- `src/hooks/use-sanity-content.ts` exports Solid.js primitives
- Uses `createResource()` or `createAsync()` internally

#### Scenario: Reactive Query with Dependencies
```typescript
const [search, setSearch] = createSignal('hello')

const [results] = createSanityQuery(() => {
  return `*[_type == "post" && title match "${search()}"]`
})

// When search() changes, query re-executes automatically
setSearch('world')
// Result: results updates with new search results
```

#### Scenario: Suspense with Fallback
```typescript
<Suspense fallback={<div>Loading...</div>}>
  <PostList />
</Suspense>

// Inside PostList:
const [posts] = createSanityQuery('*[_type == "post"]')
// Result: Suspense shows fallback until posts load
```

---

### Requirement: System SHALL provide server-side content fetching for SSR and route loaders

The system SHALL provide `fetchSanityContent()` function for server-only usage. The system SHALL integrate with TanStack Start route loaders. Server-fetched data SHALL have no client-side overhead. The system SHALL support hydration of pre-fetched data.

**Implementation**:
- `src/lib/sanity-server.ts` exports `fetchSanityContent()`
- Used in TanStack Start `server$()` functions and loaders

#### Scenario: Route Loader with SSR
```typescript
// In route loader
export const loader = async () => {
  const posts = await fetchSanityContent('*[_type == "post"]')
  return { posts }
}

// Component receives pre-hydrated data
export default function PostsPage(props: { posts: Post[] }) {
  // posts already loaded server-side
  return <div>{props.posts.map(p => <p>{p.title}</p>)}</div>
}
```

---

## MODIFIED Requirements

### Requirement: Environment configuration SHALL support Sanity credentials

The system SHALL support `VITE_SANITY_PROJECT_ID` in `.env.local`. The system SHALL support `VITE_SANITY_DATASET` in `.env.local`. The system SHALL support optional server-side `SANITY_API_TOKEN`. Environment variables SHALL be documented in `.env.example`.

**Implementation**:
- Update `.env.example` with Sanity variables
- Document in setup guide

#### Scenario: Environment Variables Present
```env
VITE_SANITY_PROJECT_ID=abc123xyz
VITE_SANITY_DATASET=production
```

---

## Performance Characteristics

- **Query Latency**: < 200ms for cached queries, < 500ms for API calls (depends on network)
- **Cache Hit Rate**: 80%+ for repeated queries
- **Bundle Size Impact**: ~30KB (gzipped) for `@sanity/client`
- **Memory**: Negligible with configurable cache limits

## Security Considerations

- API tokens stored in environment variables only (never hardcoded)
- Read-only tokens safe for client-side usage
- GROQ queries parameterized to prevent injection
- CORS configured in Sanity project for allowed domains

## Testing Strategy

- Unit tests mock `@sanity/client` for isolated testing
- Integration tests use Sanity sandbox environment
- E2E tests verify full content flow
- Type tests verify generated types match schema

## Backward Compatibility

This is a new capability - no backward compatibility concerns.
Existing hardcoded content can coexist with Sanity-fetched content during migration.

## Future Extensions

- Real-time sync via Sanity Listening API
- Webhook integration for cache invalidation
- Batch query optimization
- Visual editing integration
- Content preview mode
