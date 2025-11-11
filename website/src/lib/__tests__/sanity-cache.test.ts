import { describe, it, expect, beforeEach, vi } from 'vitest';
import {
  SanityCache,
  getGlobalCache,
  resetGlobalCache,
  cachedFetch,
  getTtlForQuery,
  invalidationHelpers,
  DEFAULT_TTLS,
} from '../sanity-cache';

describe('sanity-cache', () => {
  describe('SanityCache', () => {
    let cache: SanityCache;

    beforeEach(() => {
      cache = new SanityCache(5, 1); // Small size and short TTL for testing
    });

    describe('get and set', () => {
      it('should store and retrieve values', () => {
        cache.set('*[_type == "post"]', {}, [{ id: 1 }]);
        const result = cache.get('*[_type == "post"]', {});

        expect(result).toEqual([{ id: 1 }]);
      });

      it('should return undefined for missing keys', () => {
        const result = cache.get('*[_type == "missing"]', {});
        expect(result).toBeUndefined();
      });

      it('should handle different parameter combinations', () => {
        cache.set('*[_type == $type]', { type: 'post' }, [{ id: 1 }]);
        cache.set('*[_type == $type]', { type: 'project' }, [{ id: 2 }]);

        expect(cache.get('*[_type == $type]', { type: 'post' })).toEqual([{ id: 1 }]);
        expect(cache.get('*[_type == $type]', { type: 'project' })).toEqual([{ id: 2 }]);
      });

      it('should expire entries after TTL', async () => {
        cache.set('*[_type == "post"]', {}, [{ id: 1 }], 0.1); // 100ms TTL

        expect(cache.get('*[_type == "post"]', {})).toEqual([{ id: 1 }]);

        // Wait for expiration
        await new Promise((resolve) => setTimeout(resolve, 150));

        expect(cache.get('*[_type == "post"]', {})).toBeUndefined();
      });
    });

    describe('LRU eviction', () => {
      it('should evict least recently used entry when at capacity', () => {
        // Fill cache to capacity
        for (let i = 0; i < 5; i++) {
          cache.set(`query-${i}`, {}, { data: i });
        }

        // Access query-1 and query-2 to make them recently used
        cache.get('query-1', {});
        cache.get('query-2', {});

        // Add one more item (should evict query-0, the oldest unaccessed)
        cache.set('query-5', {}, { data: 5 });

        expect(cache.get('query-0', {})).toBeUndefined(); // Evicted (oldest)
        expect(cache.get('query-1', {})).toBeDefined(); // Still in cache (recently accessed)
        expect(cache.get('query-5', {})).toBeDefined(); // New entry
      });
    });

    describe('invalidate', () => {
      beforeEach(() => {
        cache.set('*[_type == "post"]', {}, [{ id: 1 }]);
        cache.set('*[_type == "project"]', {}, [{ id: 2 }]);
        cache.set('*[_type == "tag"]', {}, [{ id: 3 }]);
      });

      it('should invalidate entries matching string pattern', () => {
        const count = cache.invalidate('post');

        expect(count).toBe(1);
        expect(cache.get('*[_type == "post"]', {})).toBeUndefined();
        expect(cache.get('*[_type == "project"]', {})).toBeDefined();
      });

      it('should invalidate entries matching regex pattern', () => {
        const count = cache.invalidate(/post|project/);

        expect(count).toBe(2);
        expect(cache.get('*[_type == "post"]', {})).toBeUndefined();
        expect(cache.get('*[_type == "project"]', {})).toBeUndefined();
        expect(cache.get('*[_type == "tag"]', {})).toBeDefined();
      });

      it('should be case insensitive by default', () => {
        const count = cache.invalidate('POST');

        expect(count).toBe(1);
        expect(cache.get('*[_type == "post"]', {})).toBeUndefined();
      });
    });

    describe('clear', () => {
      it('should remove all entries', () => {
        cache.set('query-1', {}, { data: 1 });
        cache.set('query-2', {}, { data: 2 });

        cache.clear();

        expect(cache.get('query-1', {})).toBeUndefined();
        expect(cache.get('query-2', {})).toBeUndefined();
        expect(cache.getStats().size).toBe(0);
      });
    });

    describe('cleanup', () => {
      it('should remove only expired entries', async () => {
        cache.set('fresh', {}, { data: 1 }, 10); // 10s TTL
        cache.set('stale', {}, { data: 2 }, 0.05); // 50ms TTL

        await new Promise((resolve) => setTimeout(resolve, 100));

        const removed = cache.cleanup();

        expect(removed).toBe(1);
        expect(cache.get('fresh', {})).toBeDefined();
        expect(cache.get('stale', {})).toBeUndefined();
      });
    });

    describe('statistics', () => {
      it('should track hits and misses', () => {
        cache.set('query', {}, { data: 1 });

        cache.get('query', {}); // Hit
        cache.get('missing', {}); // Miss

        const stats = cache.getStats();
        expect(stats.hits).toBe(1);
        expect(stats.misses).toBe(1);
      });

      it('should track sets and invalidations', () => {
        cache.set('query-1', {}, { data: 1 });
        cache.set('query-2', {}, { data: 2 });
        cache.invalidate('query-1');

        const stats = cache.getStats();
        expect(stats.sets).toBe(2);
        expect(stats.invalidations).toBe(1);
      });

      it('should calculate hit rate', () => {
        cache.set('query', {}, { data: 1 });

        cache.get('query', {}); // Hit
        cache.get('query', {}); // Hit
        cache.get('missing', {}); // Miss

        expect(cache.getHitRate()).toBeCloseTo(66.67, 1);
      });

      it('should return 0 hit rate when no requests', () => {
        expect(cache.getHitRate()).toBe(0);
      });
    });
  });

  describe('global cache', () => {
    beforeEach(() => {
      resetGlobalCache();
    });

    it('should return singleton instance', () => {
      const cache1 = getGlobalCache();
      const cache2 = getGlobalCache();

      expect(cache1).toBe(cache2);
    });

    it('should reset global instance', () => {
      const cache1 = getGlobalCache();
      cache1.set('test', {}, { data: 1 });

      resetGlobalCache();

      const cache2 = getGlobalCache();
      expect(cache2).not.toBe(cache1);
      expect(cache2.get('test', {})).toBeUndefined();
    });
  });

  describe('cachedFetch', () => {
    const mockClient = {
      fetch: vi.fn(),
    };

    beforeEach(() => {
      resetGlobalCache();
      vi.clearAllMocks();
    });

    it('should fetch and cache result', async () => {
      mockClient.fetch.mockResolvedValue([{ id: 1 }]);

      const result1 = await cachedFetch(mockClient, '*[_type == "post"]', {});
      const result2 = await cachedFetch(mockClient, '*[_type == "post"]', {});

      expect(result1).toEqual([{ id: 1 }]);
      expect(result2).toEqual([{ id: 1 }]);
      expect(mockClient.fetch).toHaveBeenCalledTimes(1); // Second call uses cache
    });

    it('should skip cache when skipCache is true', async () => {
      mockClient.fetch.mockResolvedValue([{ id: 1 }]);

      const result1 = await cachedFetch(mockClient, '*[_type == "post"]', {}, { skipCache: true });
      const result2 = await cachedFetch(mockClient, '*[_type == "post"]', {}, { skipCache: true });

      expect(mockClient.fetch).toHaveBeenCalledTimes(2);
    });

    it('should use custom TTL', async () => {
      mockClient.fetch.mockResolvedValue([{ id: 1 }]);
      const customCache = new SanityCache();

      await cachedFetch(mockClient, '*[_type == "post"]', {}, { cache: customCache, ttl: 3600 });

      // Verify the entry exists (we can't easily verify TTL directly, but this ensures it was set)
      const cached = customCache.get('*[_type == "post"]', {});
      expect(cached).toEqual([{ id: 1 }]);
    });
  });

  describe('getTtlForQuery', () => {
    it('should return count TTL for count queries', () => {
      expect(getTtlForQuery('count(*[_type == "post"])')).toBe(DEFAULT_TTLS.count);
    });

    it('should return search TTL for match queries', () => {
      expect(getTtlForQuery('*[title match $term]')).toBe(DEFAULT_TTLS.search);
    });

    it('should return post TTL for post queries', () => {
      expect(getTtlForQuery('*[_type == "post"]')).toBe(DEFAULT_TTLS.post);
    });

    it('should return project TTL for project queries', () => {
      expect(getTtlForQuery('*[_type == "project"]')).toBe(DEFAULT_TTLS.project);
    });

    it('should return tag TTL for tag queries', () => {
      expect(getTtlForQuery('*[_type == "tag"]')).toBe(DEFAULT_TTLS.tag);
    });

    it('should return employment TTL for employment queries', () => {
      expect(getTtlForQuery('*[_type == "employment"]')).toBe(DEFAULT_TTLS.employment);
    });

    it('should return list TTL for list queries', () => {
      expect(getTtlForQuery('*[_type == "unknown"] | order(name)')).toBe(DEFAULT_TTLS.list);
    });

    it('should return default TTL for unknown queries', () => {
      expect(getTtlForQuery('some custom query')).toBe(DEFAULT_TTLS.default);
    });
  });

  describe('invalidationHelpers', () => {
    let cache: SanityCache;

    beforeEach(() => {
      resetGlobalCache();
      cache = getGlobalCache();

      // Populate cache with test data
      cache.set('*[_type == "post"]', {}, [{ id: 1 }]);
      cache.set('*[_type == "project"]', {}, [{ id: 2 }]);
      cache.set('*[_type == "tag"]', {}, [{ id: 3 }]);
      cache.set('*[_type == "employment"]', {}, [{ id: 4 }]);
    });

    it('should invalidate all post queries', () => {
      const count = invalidationHelpers.posts();
      expect(count).toBeGreaterThan(0);
      expect(cache.get('*[_type == "post"]', {})).toBeUndefined();
    });

    it('should invalidate all project queries', () => {
      const count = invalidationHelpers.projects();
      expect(count).toBeGreaterThan(0);
      expect(cache.get('*[_type == "project"]', {})).toBeUndefined();
    });

    it('should invalidate all tag queries', () => {
      const count = invalidationHelpers.tags();
      expect(count).toBeGreaterThan(0);
      expect(cache.get('*[_type == "tag"]', {})).toBeUndefined();
    });

    it('should invalidate all employment queries', () => {
      const count = invalidationHelpers.employments();
      expect(count).toBeGreaterThan(0);
      expect(cache.get('*[_type == "employment"]', {})).toBeUndefined();
    });

    it('should invalidate by slug', () => {
      cache.set('*[slug.current == "test-slug"]', {}, { data: 1 });
      const count = invalidationHelpers.bySlug('test-slug');
      expect(count).toBeGreaterThan(0);
    });

    it('should invalidate all list queries', () => {
      const count = invalidationHelpers.lists();
      expect(count).toBe(4); // All our test queries are list queries
    });
  });
});
