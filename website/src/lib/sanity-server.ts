import { createServerClient } from '@/lib/sanity-client';
import type { SanityClient } from '@sanity/client';
import type { SanityQueryParams, Post, Project, Tag, Employment } from '@/lib/sanity-client';
import { parseSanityError, logger, withRetry, type RetryOptions } from '@/lib/sanity-errors';

/**
 * Type definitions for server-side query results
 */
export type { Post, Project, Tag, Employment };

/**
 * Options for server-side Sanity fetching
 */
export interface ServerFetchOptions {
  /**
   * Custom Sanity client instance (optional, defaults to server client with token)
   */
  client?: SanityClient;
  /**
   * Enable retry logic (default: true)
   */
  retry?: boolean | RetryOptions;
  /**
   * Cache the result (default: true)
   */
  cache?: boolean;
  /**
   * Cache TTL in seconds (default: 3600 = 1 hour)
   */
  cacheTtl?: number;
  /**
   * Log errors (default: true)
   */
  logErrors?: boolean;
}

/**
 * Server-side function to fetch content from Sanity CMS
 *
 * This function is designed for use in TanStack Start server functions and route loaders.
 * It uses the server client with API token for authenticated requests.
 *
 * @example
 * ```typescript
 * import { server$ } from '@tanstack/solid-start/server';
 * import { fetchSanityContent } from '@/lib/sanity-server';
 *
 * export const getPosts = server$(async () => {
 *   return fetchSanityContent<Post[]>('*[_type == "post"] | order(_createdAt desc)');
 * });
 * ```
 *
 * @param query - GROQ query string
 * @param params - Query parameters (optional)
 * @param options - Fetch configuration options
 * @returns Promise resolving to query result
 * @throws SanityError if the fetch fails
 */
export async function fetchSanityContent<T = unknown>(
  query: string,
  params?: SanityQueryParams,
  options: ServerFetchOptions = {}
): Promise<T> {
  const {
    client = createServerClient(),
    retry = true,
    logErrors = true,
  } = options;

  const fetcher = async (): Promise<T> => {
    try {
      logger.debug('Server: Fetching Sanity content', { query, params });
      const result = await client.fetch<T>(query, params || {});
      logger.debug('Server: Sanity fetch succeeded', { query });
      return result;
    } catch (error) {
      const sanityError = parseSanityError(error);
      if (logErrors) {
        logger.error('Server: Sanity fetch failed', sanityError);
      }
      throw sanityError;
    }
  };

  if (retry) {
    const retryOptions = typeof retry === 'boolean' ? {} : retry;
    return withRetry(fetcher, retryOptions);
  }

  return fetcher();
}

/**
 * Fetch a single document by ID
 *
 * @example
 * ```typescript
 * export const getPost = server$(async (id: string) => {
 *   return fetchDocumentById<Post>('post', id);
 * });
 * ```
 *
 * @param type - Document type (e.g., 'post', 'project')
 * @param id - Document ID
 * @param options - Fetch configuration options
 * @returns Promise resolving to document or null
 */
export async function fetchDocumentById<T = unknown>(
  type: string,
  id: string,
  options: ServerFetchOptions = {}
): Promise<T | null> {
  const query = `*[_type == $type && _id == $id][0]`;
  const params = { type, id };
  return fetchSanityContent<T | null>(query, params, options);
}

/**
 * Fetch a single document by slug
 *
 * @example
 * ```typescript
 * export const getPostBySlug = server$(async (slug: string) => {
 *   return fetchDocumentBySlug<Post>('post', slug);
 * });
 * ```
 *
 * @param type - Document type
 * @param slug - Document slug
 * @param options - Fetch configuration options
 * @returns Promise resolving to document or null
 */
export async function fetchDocumentBySlug<T = unknown>(
  type: string,
  slug: string,
  options: ServerFetchOptions = {}
): Promise<T | null> {
  const query = `*[_type == $type && slug.current == $slug][0]`;
  const params = { type, slug };
  return fetchSanityContent<T | null>(query, params, options);
}

/**
 * Fetch documents with pagination
 *
 * @example
 * ```typescript
 * export const getPostsPage = server$(async (page: number) => {
 *   return fetchPaginatedDocuments<Post>(
 *     '*[_type == "post"] | order(_createdAt desc)',
 *     page,
 *     10
 *   );
 * });
 * ```
 *
 * @param query - Base GROQ query (without pagination)
 * @param page - Page number (0-indexed)
 * @param pageSize - Number of items per page
 * @param options - Fetch configuration options
 * @returns Promise resolving to array of documents
 */
export async function fetchPaginatedDocuments<T = unknown>(
  query: string,
  page: number,
  pageSize: number,
  options: ServerFetchOptions = {}
): Promise<T[]> {
  const offset = page * pageSize;
  const limit = offset + pageSize;
  const paginatedQuery = `${query}[$offset...$limit]`;
  const params = { offset, limit };
  return fetchSanityContent<T[]>(paginatedQuery, params, options);
}

/**
 * Fetch total count of documents matching a query
 *
 * @example
 * ```typescript
 * export const getPostCount = server$(async () => {
 *   return fetchDocumentCount('*[_type == "post"]');
 * });
 * ```
 *
 * @param query - GROQ query string
 * @param params - Query parameters (optional)
 * @param options - Fetch configuration options
 * @returns Promise resolving to count
 */
export async function fetchDocumentCount(
  query: string,
  params?: SanityQueryParams,
  options: ServerFetchOptions = {}
): Promise<number> {
  const countQuery = `count(${query})`;
  return fetchSanityContent<number>(countQuery, params, options);
}

/**
 * Fetch multiple documents by IDs
 *
 * @example
 * ```typescript
 * export const getPostsByIds = server$(async (ids: string[]) => {
 *   return fetchDocumentsByIds<Post>('post', ids);
 * });
 * ```
 *
 * @param type - Document type
 * @param ids - Array of document IDs
 * @param options - Fetch configuration options
 * @returns Promise resolving to array of documents
 */
export async function fetchDocumentsByIds<T = unknown>(
  type: string,
  ids: string[],
  options: ServerFetchOptions = {}
): Promise<T[]> {
  if (ids.length === 0) {
    return [];
  }

  const query = `*[_type == $type && _id in $ids]`;
  const params = { type, ids };
  return fetchSanityContent<T[]>(query, params, options);
}

/**
 * Fetch all documents of a specific type
 *
 * WARNING: Use with caution on types with many documents. Consider using pagination instead.
 *
 * @example
 * ```typescript
 * export const getAllTags = server$(async () => {
 *   return fetchAllDocuments<Tag>('tag');
 * });
 * ```
 *
 * @param type - Document type
 * @param orderBy - Optional order clause (e.g., '_createdAt desc')
 * @param options - Fetch configuration options
 * @returns Promise resolving to array of documents
 */
export async function fetchAllDocuments<T = unknown>(
  type: string,
  orderBy?: string,
  options: ServerFetchOptions = {}
): Promise<T[]> {
  const query = orderBy
    ? `*[_type == $type] | order(${orderBy})`
    : `*[_type == $type]`;
  const params = { type };
  return fetchSanityContent<T[]>(query, params, options);
}

/**
 * Check if a document exists
 *
 * @example
 * ```typescript
 * export const postExists = server$(async (slug: string) => {
 *   return documentExists('post', 'slug.current', slug);
 * });
 * ```
 *
 * @param type - Document type
 * @param field - Field to check (e.g., 'slug.current', '_id')
 * @param value - Value to match
 * @param options - Fetch configuration options
 * @returns Promise resolving to boolean
 */
export async function documentExists(
  type: string,
  field: string,
  value: string,
  options: ServerFetchOptions = {}
): Promise<boolean> {
  const query = `count(*[_type == $type && ${field} == $value]) > 0`;
  const params = { type, value };
  return fetchSanityContent<boolean>(query, params, options);
}

/**
 * Helper type for TanStack Start server functions
 *
 * @example
 * ```typescript
 * export const getPosts: ServerFunction<Post[]> = server$(async () => {
 *   return fetchSanityContent<Post[]>('*[_type == "post"]');
 * });
 * ```
 */
export type ServerFunction<T> = () => Promise<T>;

/**
 * Helper type for TanStack Start route loaders
 *
 * @example
 * ```typescript
 * export const Route = createRoute({
 *   component: PostPage,
 *   loader: async ({ params }) => {
 *     return fetchDocumentBySlug<Post>('post', params.slug);
 *   },
 * });
 * ```
 */
export type RouteLoader<T, TParams = Record<string, string>> = (context: {
  params: TParams;
}) => Promise<T>;
