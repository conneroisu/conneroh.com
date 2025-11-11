# Design: Sanity App SDK Integration

## Architecture Overview

```
┌─────────────────────────────────────────────────┐
│         TanStack Start + Solid.js App           │
├─────────────────────────────────────────────────┤
│                                                   │
│  ┌─────────────────────────────────────────┐   │
│  │   Server Routes & API Endpoints         │   │
│  │   (TanStack Start server functions)    │   │
│  └────────────┬────────────────────────────┘   │
│               │                                  │
│  ┌────────────▼────────────────────────────┐   │
│  │  Sanity Client Wrapper (@sanity/client)│   │
│  │  • GROQ Queries                         │   │
│  │  • Document Fetching                   │   │
│  │  • Real-time Sync (via Listening API)  │   │
│  └────────────┬────────────────────────────┘   │
│               │                                  │
│  ┌────────────▼────────────────────────────┐   │
│  │    Solid.js Reactive Components         │   │
│  │  • Signals for content state            │   │
│  │  • Fine-grained reactivity              │   │
│  │  • Suspense integration                 │   │
│  └────────────┬────────────────────────────┘   │
│               │                                  │
│  ┌────────────▼────────────────────────────┐   │
│  │    UI Components (Solid.js)             │   │
│  │  • Pages                                │   │
│  │  • Blog                                 │   │
│  │  • Dynamic content                      │   │
│  └────────────────────────────────────────┘   │
│                                                   │
└─────────────────────────────────────────────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │   Sanity Backend      │
              │ (Hosted CMS Platform) │
              │  • GROQ API           │
              │  • Content Lake       │
              │  • Webhooks           │
              │  • Assets             │
              └───────────────────────┘
```

## Key Components

### 1. Sanity Client Wrapper (`src/lib/sanity-client.ts`)

Provides a thin wrapper around `@sanity/client` with:
- Project ID and dataset configuration from environment variables
- GROQ query builder utilities
- Error handling and retry logic
- Caching strategy (if needed)

```typescript
// Usage pattern
const client = createSanityClient()
const posts = await client.fetch(
  '*[_type == "post"] | order(_createdAt desc)',
  { limit: 10 }
)
```

### 2. Solid.js Content Hooks (`src/hooks/use-sanity-content.ts`)

Creates Solid.js primitives for content fetching:
- `createSanityQuery()` - reactive query wrapper
- `createSanityDocument()` - single document fetching
- `createSanityLiveQuery()` - real-time content sync
- Suspense-compatible for SSR

```typescript
// Usage pattern
const [posts] = createSanityQuery(
  '*[_type == "post"] | order(_createdAt desc)',
  { limit: 10 }
)
```

### 3. Type Generation (`src/lib/sanity-types.ts`)

Auto-generated TypeScript types from Sanity schema:
- Run `sanity typegen` during build
- Generates complete type definitions
- Zero-cost abstraction (types-only)

### 4. Environment Configuration

```
.env.local (not committed):
VITE_SANITY_PROJECT_ID=abc123xyz
VITE_SANITY_DATASET=production
SANITY_API_TOKEN=<optional-read-token>
```

### 5. SSR Compatibility

- Use TanStack Start `server$()` for server-side content fetching
- Hydrate data through route loaders
- Solid.js `Suspense` for fallbacks during hydration

## Content Flow

### Server-Side Rendering (SSR)

```
TanStack Start Route Loader
    ↓
server$() function
    ↓
Sanity Client query
    ↓
Data returned to component
    ↓
Solid.js Suspense hydration
    ↓
HTML rendered to client
```

### Client-Side Updates

```
User interaction
    ↓
Solid.js createSanityQuery signal
    ↓
Sanity Client API call
    ↓
Response received
    ↓
Signal updated → fine-grained reactivity
    ↓
Only affected components re-render
```

## Integration Points

### API Routes (`src/routes/api/`)
- Query handler: `/api/content` - fetch Sanity content
- Webhook receiver (future): `/api/webhooks/sanity` - handle content changes

### Routes (`src/routes/`)
- Dynamic routes using content slugs
- SSR data loading via route loaders
- Suspense boundaries for loading states

### Types (`src/lib/`)
- Generated Sanity types
- Custom content type definitions
- Query response types

## Error Handling Strategy

1. **Validation Errors**: Catch GROQ syntax errors, log, return fallback data
2. **Network Errors**: Retry with exponential backoff (3 attempts max)
3. **Auth Errors**: Check API token validity, log and alert
4. **Type Errors**: TypeScript strict mode catches at compile time
5. **Runtime Errors**: Suspense fallback UI for failed fetches

## Caching Strategy

| Data | Strategy | TTL |
|------|----------|-----|
| Published content | Browser + CDN | 1 hour |
| User drafts | Browser only | Session |
| Media assets | CDN | Long-lived |
| Queries with `isDraft: true` | No caching | N/A |

## Security Considerations

1. **API Tokens**:
   - Use read-only tokens for client-side queries (CORS enabled)
   - Use full-access tokens on server only
   - Store in environment variables, never hardcode

2. **GROQ Queries**:
   - Validate/sanitize user input in dynamic queries
   - Use parameterized queries to prevent injection
   - Limit query complexity with `maxDocuments` where applicable

3. **Authentication**:
   - Leverage Sanity's token-based auth
   - Implement role-based access if needed
   - Sign API requests if using private tokens

4. **CORS**:
   - Configure Sanity project to allow client domain
   - Restrict queries to public content by default

## Testing Strategy

1. **Unit Tests**: Mock Sanity client, test query builders
2. **Integration Tests**: Test with live Sanity sandbox environment
3. **E2E Tests**: Playwright tests for full content flows
4. **Type Tests**: Verify generated types match schema

## Performance Considerations

1. **Query Optimization**:
   - Project only needed fields in GROQ
   - Use pagination for large datasets
   - Cache queries with stable parameters

2. **Bundle Size**:
   - `@sanity/client` is ~30KB gzipped
   - TypeGen output is tree-shakeable
   - Minimal impact on build output

3. **Network**:
   - Use Sanity's API CDN for read queries
   - Batch queries where possible
   - Implement request deduplication

## Migration Path

If migrating from hardcoded content:
1. Export existing content as JSON
2. Import into Sanity via API
3. Update components to fetch from Sanity
4. Remove hardcoded content files
5. Test all content-dependent routes

## Future Enhancements

- Real-time collaboration UI (multiple editors)
- Webhook-triggered incremental static regeneration
- Sanity Studio plugin for custom workflows
- Visual editing integration
- Content preview mode
