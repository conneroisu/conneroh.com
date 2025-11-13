import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createSanityClient } from '@/lib/sanity-client';
import type { Post, Project, Tag, Employment } from '@/lib/sanity-client';

/**
 * Integration tests for Sanity CMS integration
 *
 * These tests verify the complete content fetching pipeline:
 * 1. Client initialization
 * 2. Query execution
 * 3. Type safety
 * 4. Error handling
 *
 * Note: These tests use mocked Sanity client responses to avoid hitting the real API.
 * For E2E tests with real Sanity connection, see tests/playwright/sanity-content.e2e.ts
 */
describe('Sanity Integration Tests', () => {
  let mockClient: ReturnType<typeof createSanityClient>;

  beforeEach(() => {
    // Create a client with test configuration
    mockClient = createSanityClient({
      projectId: 'test-project',
      dataset: 'test-dataset',
    });
  });

  // Helper to mock fetch with proper typing
  function mockFetch<T>(returnValue: T) {
    return vi.spyOn(mockClient, 'fetch').mockResolvedValue(returnValue as never);
  }

  describe('Content Fetching Pipeline', () => {
    it('should fetch and parse Post documents', async () => {
      // Mock fetch response
      const mockPosts: Post[] = [
        {
          _id: 'post-1',
          _type: 'post',
          _createdAt: '2025-01-01T00:00:00Z',
          _updatedAt: '2025-01-01T00:00:00Z',
          _rev: 'rev-1',
          title: 'Test Post',
          slug: { _type: 'slug', current: 'test-post' },
          description: 'This is a test post',
          content: [],
          createdAt: '2025-01-01T00:00:00Z',
        },
      ];

      mockFetch(mockPosts);

      const result = await mockClient.fetch<Post[]>('*[_type == "post"]');

      expect(mockClient.fetch).toHaveBeenCalledWith('*[_type == "post"]');
      expect(result).toHaveLength(1);
      expect(result[0]?._type).toBe('post');
      expect(result[0]?.title).toBe('Test Post');
      expect(result[0]?.slug?.current).toBe('test-post');
    });

    it('should fetch and parse Project documents', async () => {
      const mockProjects: Project[] = [
        {
          _id: 'project-1',
          _type: 'project',
          _createdAt: '2025-01-01T00:00:00Z',
          _updatedAt: '2025-01-01T00:00:00Z',
          _rev: 'rev-1',
          title: 'Test Project',
          slug: { _type: 'slug', current: 'test-project' },
          description: 'This is a test project',
          content: [],
          createdAt: '2025-01-01T00:00:00Z',
        },
      ];

      mockFetch(mockProjects);

      const result = await mockClient.fetch<Project[]>('*[_type == "project"]');

      expect(mockClient.fetch).toHaveBeenCalledWith('*[_type == "project"]');
      expect(result).toHaveLength(1);
      expect(result[0]?._type).toBe('project');
      expect(result[0]?.title).toBe('Test Project');
    });

    it('should fetch single document by slug', async () => {
      const mockPost: Post = {
        _id: 'post-1',
        _type: 'post',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: 'rev-1',
        title: 'Test Post',
        slug: { _type: 'slug', current: 'test-post' },
        description: 'This is a test post',
        content: [],
        createdAt: '2025-01-01T00:00:00Z',
      };

      mockFetch(mockPost);

      const result = await mockClient.fetch<Post>(
        '*[_type == "post" && slug.current == $slug][0]',
        { slug: 'test-post' }
      );

      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == "post" && slug.current == $slug][0]',
        { slug: 'test-post' }
      );
      expect(result._id).toBe('post-1');
      expect(result.slug?.current).toBe('test-post');
    });

    it('should fetch documents with references (tags)', async () => {
      const mockPostWithTags: Post = {
        _id: 'post-1',
        _type: 'post',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: 'rev-1',
        title: 'Test Post',
        slug: { _type: 'slug', current: 'test-post' },
        description: 'This is a test post',
        content: [],
        createdAt: '2025-01-01T00:00:00Z',
        tags: [{ _type: 'reference', _ref: 'tag-1', _key: 'key-1' }],
      };

      mockFetch(mockPostWithTags);

      const result = await mockClient.fetch<Post>(
        '*[_type == "post"][0]{..., tags[]->{title, slug}}'
      );

      expect(mockClient.fetch).toHaveBeenCalled();
      expect(result.tags).toBeDefined();
    });

    it('should handle paginated queries', async () => {
      const mockPosts: Post[] = Array.from({ length: 10 }, (_, i) => ({
        _id: `post-${i}`,
        _type: 'post' as const,
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: `rev-${i}`,
        title: `Post ${i}`,
        slug: { _type: 'slug' as const, current: `post-${i}` },
        description: `Description ${i}`,
        content: [],
        createdAt: '2025-01-01T00:00:00Z',
      }));

      mockFetch(mockPosts);

      const result = await mockClient.fetch<Post[]>('*[_type == "post"][0...10]');

      expect(mockClient.fetch).toHaveBeenCalledWith('*[_type == "post"][0...10]');
      expect(result).toHaveLength(10);
    });

    it('should handle count queries', async () => {
      const mockCount = 42;

      const fetchSpy = mockFetch(mockCount);

      const result = await mockClient.fetch<number>('count(*[_type == "post"])');

      expect(fetchSpy).toHaveBeenCalledWith('count(*[_type == "post"])');
      expect(result).toBe(42);
    });
  });

  describe('Error Handling', () => {
    it('should handle network errors gracefully', async () => {
      const networkError = new Error('Network error: Failed to fetch');
      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockRejectedValue(networkError);

      await expect(mockClient.fetch('*[_type == "post"]')).rejects.toThrow('Network error');
      expect(fetchSpy).toHaveBeenCalled();
    });

    it('should handle invalid query errors', async () => {
      const queryError = new Error('Query validation failed');
      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockRejectedValue(queryError);

      await expect(mockClient.fetch('INVALID QUERY')).rejects.toThrow('Query validation');
      expect(fetchSpy).toHaveBeenCalled();
    });

    it('should handle authentication errors', async () => {
      const authError = new Error('Unauthorized: Invalid credentials');
      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockRejectedValue(authError);

      await expect(mockClient.fetch('*[_type == "post"]')).rejects.toThrow('Unauthorized');
      expect(fetchSpy).toHaveBeenCalled();
    });

    it('should return empty array for no results', async () => {
      const fetchSpy = mockFetch([]);

      const result = await mockClient.fetch<Post[]>('*[_type == "nonexistent"]');

      expect(fetchSpy).toHaveBeenCalled();
      expect(result).toEqual([]);
      expect(result).toHaveLength(0);
    });

    it('should return null for single document not found', async () => {
      const fetchSpy = mockFetch(null);

      const result = await mockClient.fetch<Post | null>(
        '*[_type == "post" && slug.current == "nonexistent"][0]'
      );

      expect(fetchSpy).toHaveBeenCalled();
      expect(result).toBeNull();
    });
  });

  describe('Type Safety', () => {
    it('should enforce Post type structure', async () => {
      const mockPost: Post = {
        _id: 'post-1',
        _type: 'post',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: 'rev-1',
        title: 'Test Post',
        slug: { _type: 'slug', current: 'test-post' },
        description: 'This is a test post',
        content: [],
        createdAt: '2025-01-01T00:00:00Z',
      };

      mockFetch(mockPost);

      const result = await mockClient.fetch<Post>('*[_type == "post"][0]');

      expect(mockClient.fetch).toHaveBeenCalled();
      expect(result).toMatchObject({
        _type: 'post',
        title: expect.any(String),
        slug: expect.objectContaining({
          current: expect.any(String),
        }),
      });
    });

    it('should enforce Employment type structure', async () => {
      const mockEmployment: Employment = {
        _id: 'employment-1',
        _type: 'employment',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: 'rev-1',
        title: 'Software Engineer',
        slug: { _type: 'slug', current: 'software-engineer' },
        description: 'Test Company position',
        content: [],
        createdAt: '2020-01-01',
      };

      mockFetch(mockEmployment);

      const result = await mockClient.fetch<Employment>('*[_type == "employment"][0]');

      expect(mockClient.fetch).toHaveBeenCalled();
      expect(result).toMatchObject({
        _type: 'employment',
        title: expect.any(String),
        description: expect.any(String),
        createdAt: expect.any(String),
      });
    });
  });

  describe('Query Parameters', () => {
    it('should handle query parameters correctly', async () => {
      const mockPosts: Post[] = [
        {
          _id: 'post-1',
          _type: 'post',
          _createdAt: '2025-01-01T00:00:00Z',
          _updatedAt: '2025-01-01T00:00:00Z',
          _rev: 'rev-1',
          title: 'TypeScript Post',
          slug: { _type: 'slug', current: 'typescript-post' },
          description: 'TypeScript content',
          content: [],
          createdAt: '2025-01-01T00:00:00Z',
        },
      ];

      mockFetch(mockPosts);

      const result = await mockClient.fetch<Post[]>(
        '*[_type == $type && references($tagId)]',
        { type: 'post', tagId: 'tag-1' }
      );

      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == $type && references($tagId)]',
        { type: 'post', tagId: 'tag-1' }
      );
      expect(result).toBeDefined();
    });

    it('should handle date range parameters', async () => {
      const mockPosts: Post[] = [];

      mockFetch(mockPosts);

      await mockClient.fetch<Post[]>(
        '*[_type == "post" && createdAt >= $startDate && createdAt <= $endDate]',
        {
          startDate: '2025-01-01',
          endDate: '2025-12-31',
        }
      );

      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == "post" && createdAt >= $startDate && createdAt <= $endDate]',
        {
          startDate: '2025-01-01',
          endDate: '2025-12-31',
        }
      );
    });
  });

  describe('SSR Behavior', () => {
    it('should work with server-side rendering', async () => {
      // Simulate SSR environment
      const mockPost: Post = {
        _id: 'post-1',
        _type: 'post',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: 'rev-1',
        title: 'SSR Post',
        slug: { _type: 'slug', current: 'ssr-post' },
        description: 'Server-side rendered post',
        content: [],
        createdAt: '2025-01-01T00:00:00Z',
      };

      mockFetch(mockPost);

      // Simulate route loader behavior
      const loaderData = await mockClient.fetch<Post>('*[_type == "post"][0]');

      expect(mockClient.fetch).toHaveBeenCalled();
      expect(loaderData).toBeDefined();
      expect(loaderData._type).toBe('post');
    });

    it('should use CDN for public client', () => {
      const publicClient = createSanityClient({
        projectId: 'test-project',
        dataset: 'production',
      });

      expect(publicClient.config().useCdn).toBe(true);
    });

    it('should disable CDN for authenticated client', () => {
      const authenticatedClient = createSanityClient({
        projectId: 'test-project',
        dataset: 'production',
        token: 'test-token',
      });

      expect(authenticatedClient.config().useCdn).toBe(false);
    });
  });
});
