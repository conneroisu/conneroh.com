import { createResource, createMemo, type ResourceReturn, type Resource } from 'solid-js';
import { getSanityClient } from '@/lib/sanity-client';
import type { SanityClient } from '@sanity/client';
import type { SanityQueryParams, Post, Project, Tag, Employment } from '@/lib/sanity-client';
import { parseSanityError, logger } from '@/lib/sanity-errors';

/**
 * Re-export generated types for convenience in components
 */
export type { Post, Project, Tag, Employment };

/**
 * Options for Sanity query resources
 */
export interface SanityQueryOptions<T> {
  /**
   * Initial data to show while loading (optional)
   */
  initialValue?: T;
  /**
   * Custom Sanity client instance (optional, defaults to public client)
   */
  client?: SanityClient;
  /**
   * Disable automatic refetching (default: false)
   */
  deferStream?: boolean;
  /**
   * Error handler callback
   */
  onError?: (error: Error) => void;
}

/**
 * Create a reactive Solid.js resource for Sanity GROQ queries
 *
 * This hook creates a Solid.js resource that fetches data from Sanity CMS
 * and automatically handles loading, error states, and reactivity.
 *
 * @example
 * ```typescript
 * const [posts, { refetch, mutate }] = createSanityQuery(
 *   '*[_type == "post"] | order(_createdAt desc)',
 *   {}
 * );
 *
 * return (
 *   <Show when={!posts.loading} fallback={<div>Loading...</div>}>
 *     <For each={posts()}>{post => <PostCard post={post} />}</For>
 *   </Show>
 * );
 * ```
 *
 * @param query - GROQ query string or accessor function
 * @param params - Query parameters or accessor function
 * @param options - Configuration options
 * @returns Solid.js resource tuple [data, controls]
 */
export function createSanityQuery<T = unknown>(
  query: string | (() => string),
  params?: SanityQueryParams | (() => SanityQueryParams),
  options: SanityQueryOptions<T> = {}
): ResourceReturn<T | undefined> {
  const { client = getSanityClient(), deferStream = false, onError } = options;

  // Create a reactive source that combines query and params
  const source = createMemo(() => {
    const q = typeof query === 'function' ? query() : query;
    const p = typeof params === 'function' ? params() : params || {};
    return { query: q, params: p };
  });

  const [data, { refetch, mutate }] = createResource(
    source,
    async ({ query: q, params: p }) => {
      try {
        logger.debug('Fetching Sanity query', { query: q, params: p });
        const result = await client.fetch<T>(q, p);
        logger.debug('Sanity query succeeded', { query: q });
        return result;
      } catch (error) {
        const sanityError = parseSanityError(error);
        logger.error('Sanity query failed', sanityError);
        if (onError) {
          onError(sanityError);
        }
        throw sanityError;
      }
    },
    {
      initialValue: options.initialValue,
      deferStream,
    }
  );

  return [data, { refetch, mutate }] as ResourceReturn<T | undefined>;
}

/**
 * Create a reactive resource for fetching a single Sanity document by ID
 *
 * @example
 * ```typescript
 * const [post, { refetch }] = createSanityDocument(
 *   'post',
 *   () => postId()
 * );
 *
 * return (
 *   <Show when={post()} fallback={<div>Loading...</div>}>
 *     {p => <article>{p.title}</article>}
 *   </Show>
 * );
 * ```
 *
 * @param type - Document type (e.g., 'post', 'project')
 * @param id - Document ID or accessor function
 * @param options - Configuration options
 * @returns Solid.js resource tuple [data, controls]
 */
export function createSanityDocument<T = unknown>(
  type: string | (() => string),
  id: string | (() => string),
  options: SanityQueryOptions<T> = {}
): ResourceReturn<T | undefined> {
  const query = () => `*[_type == $type && _id == $id][0]`;

  const params = createMemo(() => {
    const t = typeof type === 'function' ? type() : type;
    const i = typeof id === 'function' ? id() : id;
    return { type: t, id: i };
  });

  return createSanityQuery<T>(query, params, options);
}

/**
 * Create a reactive resource for fetching a single document by slug
 *
 * @example
 * ```typescript
 * const [post, { refetch }] = createSanityDocumentBySlug(
 *   'post',
 *   () => params.slug
 * );
 * ```
 *
 * @param type - Document type
 * @param slug - Document slug or accessor function
 * @param options - Configuration options
 * @returns Solid.js resource tuple [data, controls]
 */
export function createSanityDocumentBySlug<T = unknown>(
  type: string | (() => string),
  slug: string | (() => string),
  options: SanityQueryOptions<T> = {}
): ResourceReturn<T | undefined> {
  const query = createMemo(() => {
    return `*[_type == $type && slug.current == $slug][0]`;
  });

  const params = createMemo(() => {
    const t = typeof type === 'function' ? type() : type;
    const s = typeof slug === 'function' ? slug() : slug;
    return { type: t, slug: s };
  });

  return createSanityQuery<T>(query, params, options);
}

/**
 * Create a paginated query resource
 *
 * @example
 * ```typescript
 * const [page, setPage] = createSignal(0);
 * const [posts, { refetch }] = createSanityPaginatedQuery(
 *   '*[_type == "post"] | order(_createdAt desc)',
 *   page,
 *   10
 * );
 * ```
 *
 * @param query - Base GROQ query (without pagination)
 * @param page - Page number or accessor (0-indexed)
 * @param pageSize - Items per page
 * @param options - Configuration options
 * @returns Solid.js resource tuple [data, controls]
 */
export function createSanityPaginatedQuery<T = unknown>(
  query: string | (() => string),
  page: number | (() => number),
  pageSize: number,
  options: SanityQueryOptions<T[]> = {}
): ResourceReturn<T[] | undefined> {
  const paginatedQuery = createMemo(() => {
    const q = typeof query === 'function' ? query() : query;
    return `${q}[$offset...$limit]`;
  });

  const params = createMemo(() => {
    const p = typeof page === 'function' ? page() : page;
    const offset = p * pageSize;
    const limit = offset + pageSize;
    return { offset, limit };
  });

  return createSanityQuery<T[]>(paginatedQuery, params, options);
}

/**
 * Create a resource that fetches total count for a query
 *
 * @example
 * ```typescript
 * const [count] = createSanityCount('*[_type == "post"]');
 * ```
 *
 * @param query - GROQ query string
 * @param params - Query parameters
 * @param options - Configuration options
 * @returns Solid.js resource tuple [count, controls]
 */
export function createSanityCount(
  query: string | (() => string),
  params?: SanityQueryParams | (() => SanityQueryParams),
  options: SanityQueryOptions<number> = {}
): ResourceReturn<number | undefined> {
  const countQuery = createMemo(() => {
    const q = typeof query === 'function' ? query() : query;
    return `count(${q})`;
  });

  return createSanityQuery<number>(countQuery, params, options);
}

/**
 * Type guard to check if a resource has data
 *
 * @example
 * ```typescript
 * const [posts] = createSanityQuery('*[_type == "post"]');
 *
 * if (hasData(posts)) {
 *   // TypeScript knows posts() is defined here
 *   console.log(posts().length);
 * }
 * ```
 */
export function hasData<T>(resource: Resource<T | undefined>): resource is Resource<T> {
  return resource() !== undefined;
}
