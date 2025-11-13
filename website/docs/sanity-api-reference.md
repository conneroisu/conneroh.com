# Sanity CMS API Reference

Complete API documentation for all Sanity integration functions and hooks.

## Table of Contents

- [Client Configuration](#client-configuration)
- [Server Functions](#server-functions)
- [Query Builders](#query-builders)
- [Client-Side Hooks](#client-side-hooks)
- [Cache Management](#cache-management)
- [Error Handling](#error-handling)
- [Type Definitions](#type-definitions)

## Client Configuration

### `createSanityClient()`

Create a custom Sanity client instance with specific configuration.

**Import**: `import { createSanityClient } from '@/lib/sanity-client';`

**Signature**:
```typescript
function createSanityClient(options?: SanityClientOptions): SanityClient
```

**Parameters**:
- `options` (optional): Configuration object
  - `projectId?: string` - Sanity project ID (defaults to `VITE_SANITY_PROJECT_ID`)
  - `dataset?: string` - Dataset name (defaults to `VITE_SANITY_DATASET`)
  - `apiVersion?: string` - API version (defaults to `'2025-01-11'`)
  - `token?: string` - API authentication token
  - `useCdn?: boolean` - Use CDN for read operations (defaults to `true` without token, `false` with token)
  - `perspective?: 'published' | 'previewDrafts' | 'raw'` - Content perspective (defaults to `'published'`)

**Returns**: Configured `SanityClient` instance

**Example**:
```typescript
// Default client
const client = createSanityClient();

// Custom client for preview mode
const previewClient = createSanityClient({
  useCdn: false,
  perspective: 'previewDrafts',
  token: process.env.SANITY_API_TOKEN,
});
```

**Throws**:
- Error if `projectId` is not provided or found in environment
- Error if `dataset` is not provided or found in environment

---

### `getSanityClient()`

Get or create the singleton browser client for public read operations.

**Import**: `import { getSanityClient } from '@/lib/sanity-client';`

**Signature**:
```typescript
function getSanityClient(): SanityClient
```

**Returns**: Singleton `SanityClient` instance configured for CDN-enabled public reads

**Example**:
```typescript
const client = getSanityClient();
const posts = await client.fetch('*[_type == "post"]');
```

**Notes**:
- Uses CDN for optimal performance
- Caches client instance (singleton pattern)
- Suitable for client-side queries only

---

### `createServerClient()`

Create an authenticated server-side client for mutations and private operations.

**Import**: `import { createServerClient } from '@/lib/sanity-client';`

**Signature**:
```typescript
function createServerClient(): SanityClient
```

**Returns**: Authenticated `SanityClient` instance with write permissions

**Example**:
```typescript
const client = createServerClient();
await client.create({
  _type: 'post',
  title: 'New Post',
  publishedAt: new Date().toISOString(),
});
```

**Throws**:
- Error if `SANITY_API_TOKEN` environment variable is not set

**Notes**:
- Requires `SANITY_API_TOKEN` to be set
- Uses `previewDrafts` perspective by default
- CDN is disabled for server client
- Should only be used in server contexts (never expose to browser)

---

## Server Functions

### `fetchSanityContent()`

Server-side function to fetch content with retry logic and error handling.

**Import**: `import { fetchSanityContent } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function fetchSanityContent<T = unknown>(
  query: string,
  params?: SanityQueryParams,
  options?: ServerFetchOptions
): Promise<T>
```

**Parameters**:
- `query: string` - GROQ query string
- `params?: Record<string, unknown>` - Query parameters (optional)
- `options?: ServerFetchOptions` - Fetch configuration (optional)
  - `client?: SanityClient` - Custom client instance
  - `retry?: boolean | RetryOptions` - Enable retry logic (default: `true`)
  - `logErrors?: boolean` - Log errors to console (default: `true`)

**Returns**: Promise resolving to typed query result

**Example**:
```typescript
const posts = await fetchSanityContent<Post[]>(
  '*[_type == "post"]',
  {},
  { retry: { maxAttempts: 3 } }
);
```

**Throws**: `SanityError` on fetch failure after all retries

---

### `fetchDocumentById()`

Fetch a single document by its ID.

**Import**: `import { fetchDocumentById } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function fetchDocumentById<T = unknown>(
  type: string,
  id: string,
  options?: ServerFetchOptions
): Promise<T | null>
```

**Parameters**:
- `type: string` - Document type (e.g., `'post'`, `'project'`)
- `id: string` - Document ID
- `options?: ServerFetchOptions` - Fetch configuration

**Returns**: Promise resolving to document or `null` if not found

**Example**:
```typescript
const post = await fetchDocumentById<Post>('post', 'abc123');
```

---

### `fetchDocumentBySlug()`

Fetch a single document by its slug.

**Import**: `import { fetchDocumentBySlug } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function fetchDocumentBySlug<T = unknown>(
  type: string,
  slug: string,
  options?: ServerFetchOptions
): Promise<T | null>
```

**Parameters**:
- `type: string` - Document type
- `slug: string` - Document slug
- `options?: ServerFetchOptions` - Fetch configuration

**Returns**: Promise resolving to document or `null` if not found

**Example**:
```typescript
const post = await fetchDocumentBySlug<Post>('post', 'hello-world');
```

---

### `fetchPaginatedDocuments()`

Fetch documents with pagination support.

**Import**: `import { fetchPaginatedDocuments } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function fetchPaginatedDocuments<T = unknown>(
  query: string,
  page: number,
  pageSize: number,
  options?: ServerFetchOptions
): Promise<T[]>
```

**Parameters**:
- `query: string` - Base GROQ query (without pagination slice)
- `page: number` - Page number (0-indexed)
- `pageSize: number` - Items per page
- `options?: ServerFetchOptions` - Fetch configuration

**Returns**: Promise resolving to array of documents

**Example**:
```typescript
const posts = await fetchPaginatedDocuments<Post>(
  '*[_type == "post"] | order(publishedAt desc)',
  0,  // page
  10  // pageSize
);
```

---

### `fetchDocumentCount()`

Fetch total count of documents matching a query.

**Import**: `import { fetchDocumentCount } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function fetchDocumentCount(
  query: string,
  params?: SanityQueryParams,
  options?: ServerFetchOptions
): Promise<number>
```

**Parameters**:
- `query: string` - GROQ query string
- `params?: Record<string, unknown>` - Query parameters
- `options?: ServerFetchOptions` - Fetch configuration

**Returns**: Promise resolving to document count

**Example**:
```typescript
const totalPosts = await fetchDocumentCount('*[_type == "post"]');
```

---

### `fetchDocumentsByIds()`

Fetch multiple documents by their IDs.

**Import**: `import { fetchDocumentsByIds } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function fetchDocumentsByIds<T = unknown>(
  type: string,
  ids: string[],
  options?: ServerFetchOptions
): Promise<T[]>
```

**Parameters**:
- `type: string` - Document type
- `ids: string[]` - Array of document IDs
- `options?: ServerFetchOptions` - Fetch configuration

**Returns**: Promise resolving to array of documents (empty array if no IDs provided)

**Example**:
```typescript
const posts = await fetchDocumentsByIds<Post>('post', ['id1', 'id2', 'id3']);
```

---

### `fetchAllDocuments()`

Fetch all documents of a specific type.

**Import**: `import { fetchAllDocuments } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function fetchAllDocuments<T = unknown>(
  type: string,
  orderBy?: string,
  options?: ServerFetchOptions
): Promise<T[]>
```

**Parameters**:
- `type: string` - Document type
- `orderBy?: string` - Optional order clause (e.g., `'_createdAt desc'`)
- `options?: ServerFetchOptions` - Fetch configuration

**Returns**: Promise resolving to array of all documents

**Example**:
```typescript
const allTags = await fetchAllDocuments<Tag>('tag', 'name asc');
```

**Warning**: Use with caution on large datasets. Consider pagination for content types with many documents.

---

### `documentExists()`

Check if a document exists by field value.

**Import**: `import { documentExists } from '@/lib/sanity-server';`

**Signature**:
```typescript
async function documentExists(
  type: string,
  field: string,
  value: string,
  options?: ServerFetchOptions
): Promise<boolean>
```

**Parameters**:
- `type: string` - Document type
- `field: string` - Field to check (e.g., `'slug.current'`, `'_id'`)
- `value: string` - Value to match
- `options?: ServerFetchOptions` - Fetch configuration

**Returns**: Promise resolving to `true` if document exists, `false` otherwise

**Example**:
```typescript
const slugTaken = await documentExists('post', 'slug.current', 'my-post');
```

---

## Query Builders

### `listPosts()`

Generate GROQ query for listing all published posts.

**Import**: `import { listPosts } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function listPosts(): string
```

**Returns**: GROQ query string

**Example**:
```typescript
const query = listPosts();
// Result: '*[_type == "post" && !(_id in path("drafts.**"))] | order(publishedAt desc) { ... }'

const posts = await client.fetch<Post[]>(query);
```

---

### `getPostBySlug()`

Generate GROQ query for fetching a post by slug.

**Import**: `import { getPostBySlug } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getPostBySlug(): string
```

**Returns**: GROQ query string with `$slug` parameter

**Example**:
```typescript
const query = getPostBySlug();
const post = await client.fetch<Post>(query, { slug: 'hello-world' });
```

---

### `listProjects()`

Generate GROQ query for listing all projects.

**Import**: `import { listProjects } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function listProjects(): string
```

**Returns**: GROQ query string

**Example**:
```typescript
const projects = await client.fetch<Project[]>(listProjects());
```

---

### `getProjectBySlug()`

Generate GROQ query for fetching a project by slug.

**Import**: `import { getProjectBySlug } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getProjectBySlug(): string
```

**Returns**: GROQ query string with `$slug` parameter

**Example**:
```typescript
const project = await client.fetch<Project>(
  getProjectBySlug(),
  { slug: 'my-project' }
);
```

---

### `listTags()`

Generate GROQ query for listing all tags.

**Import**: `import { listTags } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function listTags(): string
```

**Returns**: GROQ query string ordered by name

**Example**:
```typescript
const tags = await client.fetch<Tag[]>(listTags());
```

---

### `getTagBySlug()`

Generate GROQ query for fetching a tag by slug.

**Import**: `import { getTagBySlug } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getTagBySlug(): string
```

**Returns**: GROQ query string with `$slug` parameter

**Example**:
```typescript
const tag = await client.fetch<Tag>(getTagBySlug(), { slug: 'typescript' });
```

---

### `listEmployments()`

Generate GROQ query for listing employment records.

**Import**: `import { listEmployments } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function listEmployments(): string
```

**Returns**: GROQ query string ordered by start date (newest first)

**Example**:
```typescript
const employments = await client.fetch<Employment[]>(listEmployments());
```

---

### `getEmploymentBySlug()`

Generate GROQ query for fetching an employment record by slug.

**Import**: `import { getEmploymentBySlug } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getEmploymentBySlug(): string
```

**Returns**: GROQ query string with `$slug` parameter

**Example**:
```typescript
const employment = await client.fetch<Employment>(
  getEmploymentBySlug(),
  { slug: 'acme-corp' }
);
```

---

### `getPostsByTag()`

Generate GROQ query for fetching posts by tag slug.

**Import**: `import { getPostsByTag } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getPostsByTag(): string
```

**Returns**: GROQ query string with `$tagSlug` parameter

**Example**:
```typescript
const posts = await client.fetch<Post[]>(
  getPostsByTag(),
  { tagSlug: 'typescript' }
);
```

---

### `getProjectsByTag()`

Generate GROQ query for fetching projects by tag slug.

**Import**: `import { getProjectsByTag } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getProjectsByTag(): string
```

**Returns**: GROQ query string with `$tagSlug` parameter

**Example**:
```typescript
const projects = await client.fetch<Project[]>(
  getProjectsByTag(),
  { tagSlug: 'react' }
);
```

---

### `searchPosts()`

Generate GROQ query for searching posts by title or excerpt.

**Import**: `import { searchPosts } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function searchPosts(): string
```

**Returns**: GROQ query string with `$searchTerm` parameter

**Example**:
```typescript
const results = await client.fetch<Post[]>(
  searchPosts(),
  { searchTerm: 'solid.js' }
);
```

---

### `withPagination()`

Add pagination slice to a GROQ query.

**Import**: `import { withPagination } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function withPagination(
  query: string,
  offset: number,
  limit: number
): string
```

**Parameters**:
- `query: string` - Base GROQ query
- `offset: number` - Starting index (0-based, must be >= 0)
- `limit: number` - Number of items to fetch (must be > 0 and <= 100)

**Returns**: Paginated query string

**Example**:
```typescript
const query = withPagination(listPosts(), 0, 10);
const posts = await client.fetch<Post[]>(query, { offset: 0, endIndex: 10 });
```

**Throws**: `SanityValidationError` if parameters are invalid

---

### `getPaginationParams()`

Generate pagination parameters for use with `withPagination`.

**Import**: `import { getPaginationParams } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getPaginationParams(
  page: number,
  pageSize: number
): { offset: number; endIndex: number }
```

**Parameters**:
- `page: number` - Page number (0-indexed)
- `pageSize: number` - Items per page

**Returns**: Object with `offset` and `endIndex` for query parameters

**Example**:
```typescript
const params = getPaginationParams(2, 10);
// Result: { offset: 20, endIndex: 30 }

const query = withPagination(listPosts(), params.offset, 10);
const posts = await client.fetch<Post[]>(query, params);
```

---

### `withDateRange()`

Add date range filter to a GROQ query.

**Import**: `import { withDateRange } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function withDateRange(
  query: string,
  field: string,
  startDate: string,
  endDate: string
): string
```

**Parameters**:
- `query: string` - Base GROQ query
- `field: string` - Date field to filter on (e.g., `'publishedAt'`)
- `startDate: string` - Start date (ISO string)
- `endDate: string` - End date (ISO string)

**Returns**: Filtered query string

**Example**:
```typescript
const query = withDateRange(
  listPosts(),
  'publishedAt',
  '2024-01-01',
  '2024-12-31'
);
const posts = await client.fetch<Post[]>(query);
```

---

### `countPosts()`

Generate GROQ query for counting posts.

**Import**: `import { countPosts } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function countPosts(): string
```

**Returns**: GROQ count query string

**Example**:
```typescript
const count = await client.fetch<number>(countPosts());
```

---

### `countProjects()`

Generate GROQ query for counting projects.

**Import**: `import { countProjects } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function countProjects(): string
```

**Returns**: GROQ count query string

**Example**:
```typescript
const count = await client.fetch<number>(countProjects());
```

---

### `countPostsByTag()`

Generate GROQ query for counting posts by tag.

**Import**: `import { countPostsByTag } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function countPostsByTag(): string
```

**Returns**: GROQ count query string with `$tagSlug` parameter

**Example**:
```typescript
const count = await client.fetch<number>(
  countPostsByTag(),
  { tagSlug: 'typescript' }
);
```

---

### `getRelatedPosts()`

Generate GROQ query for fetching related posts by shared tags.

**Import**: `import { getRelatedPosts } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function getRelatedPosts(): string
```

**Returns**: GROQ query string with `$postId` and `$limit` parameters

**Example**:
```typescript
const related = await client.fetch<Post[]>(
  getRelatedPosts(),
  { postId: 'abc123', limit: 3 }
);
```

---

### Validation Functions

#### `validateSlug()`

Validate slug parameter format.

**Import**: `import { validateSlug } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function validateSlug(slug: unknown): asserts slug is string
```

**Throws**: `SanityValidationError` if slug is invalid

**Example**:
```typescript
validateSlug(userInput); // Throws if invalid
// userInput is now typed as string
```

---

#### `validateId()`

Validate document ID parameter.

**Import**: `import { validateId } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function validateId(id: unknown): asserts id is string
```

**Throws**: `SanityValidationError` if ID is invalid

---

#### `validateSearchTerm()`

Validate search term parameter (max 200 characters).

**Import**: `import { validateSearchTerm } from '@/lib/sanity-queries';`

**Signature**:
```typescript
function validateSearchTerm(searchTerm: unknown): asserts searchTerm is string
```

**Throws**: `SanityValidationError` if search term is invalid

---

## Client-Side Hooks

### `createSanityQuery()`

Create a reactive Solid.js resource for GROQ queries.

**Import**: `import { createSanityQuery } from '@/hooks/use-sanity-content';`

**Signature**:
```typescript
function createSanityQuery<T = unknown>(
  query: string | (() => string),
  params?: SanityQueryParams | (() => SanityQueryParams),
  options?: SanityQueryOptions<T>
): ResourceReturn<T | undefined>
```

**Parameters**:
- `query` - GROQ query string or accessor function
- `params` - Query parameters or accessor function
- `options` - Configuration options
  - `initialValue?: T` - Initial data while loading
  - `client?: SanityClient` - Custom client instance
  - `deferStream?: boolean` - Disable automatic refetching
  - `onError?: (error: Error) => void` - Error handler callback

**Returns**: Solid.js resource tuple `[data, { refetch, mutate }]`

**Example**:
```typescript
const [posts, { refetch }] = createSanityQuery<Post[]>(
  listPosts(),
  {},
  { initialValue: [] }
);

// Access data
const postData = posts();

// Refetch manually
refetch();
```

---

### `createSanityDocument()`

Create a reactive resource for fetching a single document by ID.

**Import**: `import { createSanityDocument } from '@/hooks/use-sanity-content';`

**Signature**:
```typescript
function createSanityDocument<T = unknown>(
  type: string | (() => string),
  id: string | (() => string),
  options?: SanityQueryOptions<T>
): ResourceReturn<T | undefined>
```

**Parameters**:
- `type` - Document type or accessor function
- `id` - Document ID or accessor function
- `options` - Configuration options

**Returns**: Solid.js resource tuple

**Example**:
```typescript
const [post] = createSanityDocument<Post>(
  'post',
  () => postId()
);
```

---

### `createSanityDocumentBySlug()`

Create a reactive resource for fetching a document by slug.

**Import**: `import { createSanityDocumentBySlug } from '@/hooks/use-sanity-content';`

**Signature**:
```typescript
function createSanityDocumentBySlug<T = unknown>(
  type: string | (() => string),
  slug: string | (() => string),
  options?: SanityQueryOptions<T>
): ResourceReturn<T | undefined>
```

**Parameters**:
- `type` - Document type or accessor function
- `slug` - Document slug or accessor function
- `options` - Configuration options

**Returns**: Solid.js resource tuple

**Example**:
```typescript
const params = useParams();
const [post] = createSanityDocumentBySlug<Post>(
  'post',
  () => params.slug
);
```

---

### `createSanityPaginatedQuery()`

Create a reactive resource for paginated queries.

**Import**: `import { createSanityPaginatedQuery } from '@/hooks/use-sanity-content';`

**Signature**:
```typescript
function createSanityPaginatedQuery<T = unknown>(
  query: string | (() => string),
  page: number | (() => number),
  pageSize: number,
  options?: SanityQueryOptions<T[]>
): ResourceReturn<T[] | undefined>
```

**Parameters**:
- `query` - Base GROQ query (without pagination)
- `page` - Page number or accessor (0-indexed)
- `pageSize` - Items per page
- `options` - Configuration options

**Returns**: Solid.js resource tuple

**Example**:
```typescript
const [page, setPage] = createSignal(0);
const [posts] = createSanityPaginatedQuery<Post>(
  listPosts(),
  page,
  10
);
```

---

### `createSanityCount()`

Create a reactive resource for counting documents.

**Import**: `import { createSanityCount } from '@/hooks/use-sanity-content';`

**Signature**:
```typescript
function createSanityCount(
  query: string | (() => string),
  params?: SanityQueryParams | (() => SanityQueryParams),
  options?: SanityQueryOptions<number>
): ResourceReturn<number | undefined>
```

**Parameters**:
- `query` - GROQ query string or accessor
- `params` - Query parameters or accessor
- `options` - Configuration options

**Returns**: Solid.js resource tuple with count

**Example**:
```typescript
const [count] = createSanityCount('*[_type == "post"]');
```

---

### `hasData()`

Type guard to check if a resource has loaded data.

**Import**: `import { hasData } from '@/hooks/use-sanity-content';`

**Signature**:
```typescript
function hasData<T>(
  resource: Resource<T | undefined>
): resource is Resource<T>
```

**Parameters**:
- `resource` - Solid.js resource to check

**Returns**: `true` if resource has data, narrows type to non-undefined

**Example**:
```typescript
const [posts] = createSanityQuery<Post[]>(listPosts(), {});

if (hasData(posts)) {
  // TypeScript knows posts() is Post[] here
  console.log(posts().length);
}
```

---

## Cache Management

### `SanityCache`

LRU cache implementation for Sanity queries.

**Import**: `import { SanityCache } from '@/lib/sanity-cache';`

**Constructor**:
```typescript
new SanityCache(maxSize?: number, defaultTtl?: number)
```

**Parameters**:
- `maxSize?: number` - Maximum cache entries (default: 100)
- `defaultTtl?: number` - Default TTL in seconds (default: 3600)

**Methods**:

#### `get<T>(query, params?): T | undefined`

Get cached value.

**Example**:
```typescript
const cached = cache.get<Post[]>('*[_type == "post"]');
```

---

#### `set<T>(query, params, data, ttl?): void`

Set cached value with TTL.

**Example**:
```typescript
cache.set('*[_type == "post"]', {}, posts, 3600);
```

---

#### `invalidate(pattern): number`

Invalidate entries matching pattern.

**Example**:
```typescript
// Invalidate all post queries
cache.invalidate(/post/);

// Invalidate specific query
cache.invalidate('*[_type == "post"]');
```

---

#### `clear(): void`

Clear all cache entries.

---

#### `getStats(): CacheStats`

Get cache statistics.

**Returns**:
```typescript
{
  hits: number;
  misses: number;
  sets: number;
  invalidations: number;
  size: number;
}
```

---

#### `getHitRate(): number`

Get cache hit rate percentage.

---

#### `cleanup(): number`

Remove expired entries manually.

---

### `getGlobalCache()`

Get or create the global singleton cache instance.

**Import**: `import { getGlobalCache } from '@/lib/sanity-cache';`

**Signature**:
```typescript
function getGlobalCache(maxSize?: number, defaultTtl?: number): SanityCache
```

**Example**:
```typescript
const cache = getGlobalCache();
```

---

### `cachedFetch()`

Cached fetch wrapper for Sanity queries.

**Import**: `import { cachedFetch } from '@/lib/sanity-cache';`

**Signature**:
```typescript
async function cachedFetch<T>(
  client: { fetch: (query: string, params?: Record<string, unknown>) => Promise<T> },
  query: string,
  params?: Record<string, unknown>,
  options?: {
    cache?: SanityCache;
    ttl?: number;
    skipCache?: boolean;
  }
): Promise<T>
```

**Example**:
```typescript
const posts = await cachedFetch(
  client,
  '*[_type == "post"]',
  {},
  { ttl: 3600 }
);
```

---

### `getTtlForQuery()`

Determine appropriate TTL based on query content.

**Import**: `import { getTtlForQuery } from '@/lib/sanity-cache';`

**Signature**:
```typescript
function getTtlForQuery(query: string): number
```

**Returns**: Recommended TTL in seconds

**Example**:
```typescript
const ttl = getTtlForQuery('*[_type == "post"]');
// Returns: 3600 (1 hour for posts)
```

---

### `invalidationHelpers`

Content type-specific cache invalidation helpers.

**Import**: `import { invalidationHelpers } from '@/lib/sanity-cache';`

**Methods**:
- `posts(cache?): number` - Invalidate all post queries
- `projects(cache?): number` - Invalidate all project queries
- `tags(cache?): number` - Invalidate all tag queries
- `employments(cache?): number` - Invalidate all employment queries
- `bySlug(slug, cache?): number` - Invalidate queries containing slug
- `lists(cache?): number` - Invalidate all list queries

**Example**:
```typescript
// Invalidate all post-related queries
invalidationHelpers.posts();

// Invalidate specific slug
invalidationHelpers.bySlug('hello-world');
```

---

## Error Handling

### Error Classes

#### `SanityError`

Base error class for Sanity operations.

**Import**: `import { SanityError } from '@/lib/sanity-errors';`

**Properties**:
- `name: string` - Error name
- `message: string` - Error message
- `code: string` - Error code
- `statusCode?: number` - HTTP status code
- `details?: unknown` - Additional error details

---

#### `SanityValidationError`

Validation error (invalid parameters, etc.).

**Import**: `import { SanityValidationError } from '@/lib/sanity-errors';`

---

#### `SanityNetworkError`

Network-related error (timeout, connection failure).

**Import**: `import { SanityNetworkError } from '@/lib/sanity-errors';`

---

#### `SanityAuthError`

Authentication/authorization error.

**Import**: `import { SanityAuthError } from '@/lib/sanity-errors';`

---

#### `SanityRateLimitError`

Rate limit exceeded error.

**Import**: `import { SanityRateLimitError } from '@/lib/sanity-errors';`

---

### `parseSanityError()`

Normalize errors from Sanity API into typed error classes.

**Import**: `import { parseSanityError } from '@/lib/sanity-errors';`

**Signature**:
```typescript
function parseSanityError(error: unknown): SanityError
```

**Example**:
```typescript
try {
  await client.fetch(query);
} catch (error) {
  const sanityError = parseSanityError(error);
  console.error(`Error ${sanityError.code}:`, sanityError.message);
}
```

---

### `withRetry()`

Execute function with exponential backoff retry logic.

**Import**: `import { withRetry } from '@/lib/sanity-errors';`

**Signature**:
```typescript
async function withRetry<T>(
  fn: () => Promise<T>,
  options?: RetryOptions
): Promise<T>
```

**Parameters**:
- `fn` - Async function to execute
- `options` - Retry configuration
  - `maxAttempts?: number` - Max retry attempts (default: 3)
  - `baseDelay?: number` - Base delay in ms (default: 1000)
  - `maxDelay?: number` - Max delay in ms (default: 10000)
  - `shouldRetry?: (error: Error) => boolean` - Custom retry predicate

**Example**:
```typescript
const posts = await withRetry(
  () => client.fetch('*[_type == "post"]'),
  { maxAttempts: 3, baseDelay: 1000 }
);
```

---

### `safeFetch()`

Safe fetch wrapper with fallback support.

**Import**: `import { safeFetch } from '@/lib/sanity-errors';`

**Signature**:
```typescript
async function safeFetch<T>(
  fn: () => Promise<T>,
  options?: {
    fallback?: T;
    logError?: boolean;
  }
): Promise<T>
```

**Example**:
```typescript
const posts = await safeFetch(
  () => client.fetch('*[_type == "post"]'),
  { fallback: [] }
);
```

---

## Type Definitions

### Document Types

Generated from Sanity schema:

- `Post` - Blog post document
- `Project` - Portfolio project document
- `Tag` - Tag document
- `Employment` - Employment history document
- `BlockContent` - Portable Text content blocks
- `Slug` - Slug field type
- `SanityImageAsset` - Image asset with metadata
- `SanityImageCrop` - Image crop data
- `SanityImageHotspot` - Image hotspot data
- `AllSanitySchemaTypes` - Union of all schema types

**Import**: `import type { Post, Project, Tag, Employment } from '@/lib/sanity-client';`

---

### Query Types

- `SanityQueryParams` - Query parameters object
- `SanityQueryResult<T>` - Typed query result

---

### Options Types

- `SanityClientOptions` - Client configuration options
- `ServerFetchOptions` - Server fetch configuration
- `SanityQueryOptions<T>` - Solid.js hook configuration
- `RetryOptions` - Retry logic configuration

---

## Related Documentation

- [Setup Guide](./sanity-setup.md) - Initial configuration
- [Usage Examples](./sanity-usage-examples.md) - Practical examples
- [Type Generation](./SANITY_TYPEGEN.md) - TypeScript type generation
- [Sanity Documentation](https://www.sanity.io/docs) - Official Sanity docs
