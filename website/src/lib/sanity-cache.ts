import { logger } from '@/lib/sanity-errors';

/**
 * Cache entry with TTL support
 */
interface CacheEntry<T> {
  data: T;
  expiresAt: number;
  key: string;
}

/**
 * Cache statistics for monitoring
 */
interface CacheStats {
  hits: number;
  misses: number;
  sets: number;
  invalidations: number;
  size: number;
}

/**
 * TTL configurations for different content types (in seconds)
 */
export const DEFAULT_TTLS = {
  post: 3600, // 1 hour - posts don't change often
  project: 3600, // 1 hour - projects are relatively stable
  tag: 7200, // 2 hours - tags change rarely
  employment: 7200, // 2 hours - employment history is stable
  list: 1800, // 30 minutes - lists might update more frequently
  count: 900, // 15 minutes - counts can change
  search: 600, // 10 minutes - search results should be fresh
  default: 3600, // 1 hour - default fallback
} as const;

/**
 * LRU (Least Recently Used) Cache implementation for Sanity queries
 *
 * This cache automatically evicts the least recently used entries when the
 * maximum size is reached. All entries have TTL-based expiration.
 */
export class SanityCache {
  private cache = new Map<string, CacheEntry<unknown>>();
  private accessOrder = new Map<string, number>(); // Track access time for LRU
  private stats: CacheStats = {
    hits: 0,
    misses: 0,
    sets: 0,
    invalidations: 0,
    size: 0,
  };

  constructor(
    private maxSize = 100, // Maximum number of cache entries
    private defaultTtl = DEFAULT_TTLS.default // Default TTL in seconds
  ) {}

  /**
   * Generate cache key from query and parameters
   */
  private generateKey(query: string, params?: Record<string, unknown>): string {
    const paramsStr = params ? JSON.stringify(params, Object.keys(params).sort()) : '';
    return `${query}::${paramsStr}`;
  }

  /**
   * Check if entry is expired
   */
  private isExpired(entry: CacheEntry<unknown>): boolean {
    return Date.now() > entry.expiresAt;
  }

  /**
   * Evict least recently used entry
   */
  private evictLRU(): void {
    if (this.accessOrder.size === 0) return;

    // Find the entry with the oldest access time
    let oldestKey: string | null = null;
    let oldestTime = Infinity;

    for (const [key, time] of this.accessOrder.entries()) {
      if (time < oldestTime) {
        oldestTime = time;
        oldestKey = key;
      }
    }

    if (oldestKey) {
      this.cache.delete(oldestKey);
      this.accessOrder.delete(oldestKey);
      logger.debug('Cache: Evicted LRU entry', { key: oldestKey });
    }
  }

  /**
   * Update access time for LRU tracking
   */
  private updateAccessTime(key: string): void {
    this.accessOrder.set(key, Date.now());
  }

  /**
   * Get value from cache
   *
   * @param query - GROQ query string
   * @param params - Query parameters
   * @returns Cached value or undefined if not found/expired
   */
  get<T>(query: string, params?: Record<string, unknown>): T | undefined {
    const key = this.generateKey(query, params);
    const entry = this.cache.get(key) as CacheEntry<T> | undefined;

    if (!entry) {
      this.stats.misses++;
      logger.debug('Cache: Miss', { key });
      return undefined;
    }

    if (this.isExpired(entry)) {
      this.cache.delete(key);
      this.accessOrder.delete(key);
      this.stats.misses++;
      logger.debug('Cache: Expired', { key });
      return undefined;
    }

    this.updateAccessTime(key);
    this.stats.hits++;
    logger.debug('Cache: Hit', { key });
    return entry.data;
  }

  /**
   * Set value in cache with TTL
   *
   * @param query - GROQ query string
   * @param params - Query parameters
   * @param data - Data to cache
   * @param ttl - Time to live in seconds (optional, uses default if not provided)
   */
  set<T>(
    query: string,
    params: Record<string, unknown> | undefined,
    data: T,
    ttl?: number
  ): void {
    const key = this.generateKey(query, params);

    // Evict LRU if at capacity and this is a new key
    if (this.cache.size >= this.maxSize && !this.cache.has(key)) {
      this.evictLRU();
    }

    const ttlSeconds = ttl ?? this.defaultTtl;
    const expiresAt = Date.now() + ttlSeconds * 1000;

    this.cache.set(key, {
      data,
      expiresAt,
      key,
    });

    this.updateAccessTime(key);
    this.stats.sets++;
    this.stats.size = this.cache.size;

    logger.debug('Cache: Set', { key, ttl: ttlSeconds, expiresAt: new Date(expiresAt) });
  }

  /**
   * Invalidate cache entries matching a pattern
   *
   * @param pattern - String or RegExp to match against cache keys
   * @returns Number of invalidated entries
   *
   * @example
   * ```typescript
   * // Invalidate all post queries
   * cache.invalidate(/post/);
   *
   * // Invalidate specific query
   * cache.invalidate('*[_type == "post"]');
   * ```
   */
  invalidate(pattern: string | RegExp): number {
    const regex = typeof pattern === 'string' ? new RegExp(pattern, 'i') : pattern;
    let count = 0;

    for (const [key] of this.cache.entries()) {
      if (regex.test(key)) {
        this.cache.delete(key);
        this.accessOrder.delete(key);
        count++;
      }
    }

    this.stats.invalidations += count;
    this.stats.size = this.cache.size;

    logger.debug('Cache: Invalidated entries', { pattern: pattern.toString(), count });
    return count;
  }

  /**
   * Invalidate all cache entries
   */
  clear(): void {
    const size = this.cache.size;
    this.cache.clear();
    this.accessOrder.clear();
    this.stats.invalidations += size;
    this.stats.size = 0;
    logger.debug('Cache: Cleared all entries', { count: size });
  }

  /**
   * Get cache statistics
   */
  getStats(): Readonly<CacheStats> {
    return { ...this.stats };
  }

  /**
   * Get hit rate percentage
   */
  getHitRate(): number {
    const total = this.stats.hits + this.stats.misses;
    return total === 0 ? 0 : (this.stats.hits / total) * 100;
  }

  /**
   * Remove expired entries (manual cleanup)
   */
  cleanup(): number {
    let count = 0;

    for (const [key, entry] of this.cache.entries()) {
      if (this.isExpired(entry)) {
        this.cache.delete(key);
        this.accessOrder.delete(key);
        count++;
      }
    }

    this.stats.size = this.cache.size;
    logger.debug('Cache: Cleanup completed', { expired: count });
    return count;
  }
}

/**
 * Global singleton cache instance
 */
let _globalCache: SanityCache | null = null;

/**
 * Get or create the global cache instance
 *
 * @param maxSize - Maximum cache size (only used on first call)
 * @param defaultTtl - Default TTL in seconds (only used on first call)
 * @returns Global cache instance
 */
export function getGlobalCache(maxSize = 100, defaultTtl = DEFAULT_TTLS.default): SanityCache {
  if (!_globalCache) {
    _globalCache = new SanityCache(maxSize, defaultTtl);
    logger.debug('Cache: Created global instance', { maxSize, defaultTtl });
  }
  return _globalCache;
}

/**
 * Reset the global cache instance (useful for testing)
 */
export function resetGlobalCache(): void {
  _globalCache = null;
  logger.debug('Cache: Reset global instance');
}

/**
 * Cached fetch wrapper for Sanity queries
 *
 * @example
 * ```typescript
 * const posts = await cachedFetch(
 *   client,
 *   '*[_type == "post"]',
 *   {},
 *   { ttl: 3600 }
 * );
 * ```
 *
 * @param client - Sanity client instance
 * @param query - GROQ query string
 * @param params - Query parameters
 * @param options - Cache options
 * @returns Promise resolving to query result
 */
export async function cachedFetch<T>(
  client: { fetch: (query: string, params?: Record<string, unknown>) => Promise<T> },
  query: string,
  params?: Record<string, unknown>,
  options: {
    cache?: SanityCache;
    ttl?: number;
    skipCache?: boolean;
  } = {}
): Promise<T> {
  const { cache = getGlobalCache(), ttl, skipCache = false } = options;

  // Check cache first
  if (!skipCache) {
    const cached = cache.get<T>(query, params);
    if (cached !== undefined) {
      return cached;
    }
  }

  // Fetch from Sanity
  const result = await client.fetch(query, params);

  // Cache the result
  if (!skipCache) {
    cache.set(query, params, result, ttl);
  }

  return result;
}

/**
 * Helper to determine appropriate TTL based on query content
 *
 * @param query - GROQ query string
 * @returns Recommended TTL in seconds
 *
 * @example
 * ```typescript
 * const ttl = getTtlForQuery('*[_type == "post"]');
 * cache.set(query, params, data, ttl);
 * ```
 */
export function getTtlForQuery(query: string): number {
  const lowerQuery = query.toLowerCase();

  if (lowerQuery.includes('count(')) {
    return DEFAULT_TTLS.count;
  }

  if (lowerQuery.includes('match')) {
    return DEFAULT_TTLS.search;
  }

  if (lowerQuery.includes('"post"')) {
    return DEFAULT_TTLS.post;
  }

  if (lowerQuery.includes('"project"')) {
    return DEFAULT_TTLS.project;
  }

  if (lowerQuery.includes('"tag"')) {
    return DEFAULT_TTLS.tag;
  }

  if (lowerQuery.includes('"employment"')) {
    return DEFAULT_TTLS.employment;
  }

  // If it's a list query (array results)
  if (lowerQuery.includes('[') || lowerQuery.includes('order')) {
    return DEFAULT_TTLS.list;
  }

  return DEFAULT_TTLS.default;
}

/**
 * Create content type-specific invalidation helpers
 */
export const invalidationHelpers = {
  /**
   * Invalidate all post-related queries
   */
  posts: (cache: SanityCache = getGlobalCache()): number => {
    return cache.invalidate(/post/i);
  },

  /**
   * Invalidate all project-related queries
   */
  projects: (cache: SanityCache = getGlobalCache()): number => {
    return cache.invalidate(/project/i);
  },

  /**
   * Invalidate all tag-related queries
   */
  tags: (cache: SanityCache = getGlobalCache()): number => {
    return cache.invalidate(/tag/i);
  },

  /**
   * Invalidate all employment-related queries
   */
  employments: (cache: SanityCache = getGlobalCache()): number => {
    return cache.invalidate(/employment/i);
  },

  /**
   * Invalidate queries containing a specific slug
   */
  bySlug: (slug: string, cache: SanityCache = getGlobalCache()): number => {
    return cache.invalidate(new RegExp(slug, 'i'));
  },

  /**
   * Invalidate all list queries
   */
  lists: (cache: SanityCache = getGlobalCache()): number => {
    return cache.invalidate(/\*\[_type/i);
  },
};
