import { createClient, type ClientConfig, type SanityClient } from '@sanity/client';
import { env } from '../env';

/**
 * Re-export generated Sanity types for convenience
 */
export type {
  Post,
  Project,
  Tag,
  Employment,
  BlockContent,
  Slug,
  SanityImageAsset,
  SanityImageCrop,
  SanityImageHotspot,
  AllSanitySchemaTypes,
} from './sanity-types';

/**
 * Sanity client configuration options
 */
export interface SanityClientOptions {
  /**
   * Project ID from Sanity dashboard
   */
  projectId?: string;
  /**
   * Dataset name (e.g., 'production', 'development')
   */
  dataset?: string;
  /**
   * API version in format YYYY-MM-DD
   */
  apiVersion?: string;
  /**
   * API token for authenticated requests (server-side only)
   */
  token?: string;
  /**
   * Use CDN for faster read operations (default: true for browser, false with token)
   */
  useCdn?: boolean;
  /**
   * Enable perspective for draft content
   */
  perspective?: 'published' | 'previewDrafts' | 'raw';
}

/**
 * Create a configured Sanity client instance
 *
 * @example
 * ```typescript
 * // Browser client (public, read-only)
 * const client = createSanityClient();
 *
 * // Server client (authenticated, with token)
 * const serverClient = createSanityClient({
 *   token: process.env.SANITY_API_TOKEN,
 *   useCdn: false
 * });
 * ```
 *
 * @param options - Client configuration options
 * @returns Configured Sanity client instance
 * @throws Error if required environment variables are missing
 */
export function createSanityClient(options: SanityClientOptions = {}): SanityClient {
  const projectId = options.projectId || env.VITE_SANITY_PROJECT_ID;
  const dataset = options.dataset || env.VITE_SANITY_DATASET;

  if (!projectId) {
    throw new Error(
      'Sanity project ID is required. Set VITE_SANITY_PROJECT_ID in environment variables or pass projectId option.'
    );
  }

  if (!dataset) {
    throw new Error(
      'Sanity dataset is required. Set VITE_SANITY_DATASET in environment variables or pass dataset option.'
    );
  }

  const config: ClientConfig = {
    projectId,
    dataset,
    apiVersion: options.apiVersion || '2025-01-11',
    useCdn: options.useCdn !== undefined ? options.useCdn : !options.token,
    perspective: options.perspective || 'published',
  };

  if (options.token) {
    config.token = options.token;
  }

  return createClient(config);
}

/**
 * Get or create the default browser client for public read operations
 * Uses CDN for optimal performance
 *
 * @returns Singleton Sanity client instance
 */
let _sanityClient: SanityClient | null = null;

export function getSanityClient(): SanityClient {
  if (!_sanityClient) {
    _sanityClient = createSanityClient();
  }
  return _sanityClient;
}

/**
 * Create a server-side client with authentication
 * Requires SANITY_API_TOKEN environment variable
 *
 * @example
 * ```typescript
 * import { createServerClient } from '@/lib/sanity-client';
 *
 * const client = createServerClient();
 * await client.create({ _type: 'post', title: 'Hello World' });
 * ```
 *
 * @returns Authenticated Sanity client for server-side operations
 * @throws Error if SANITY_API_TOKEN is not set
 */
export function createServerClient(): SanityClient {
  const token = env.SANITY_API_TOKEN;

  if (!token) {
    throw new Error(
      'SANITY_API_TOKEN is required for server-side operations. Add it to your .env.local file.'
    );
  }

  return createSanityClient({
    token,
    useCdn: false,
    perspective: 'previewDrafts',
  });
}

/**
 * Type helper for Sanity query results
 */
export type SanityQueryResult<T> = T;

/**
 * Type helper for Sanity query parameters
 */
export type SanityQueryParams = Record<string, unknown>;
