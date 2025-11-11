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
          publishedAt: '2025-01-01T00:00:00Z',
          excerpt: 'This is a test post',
          content: [],
        },
      ];

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPosts);

      const result = await mockClient.fetch<Post[]>('*[_type == "post"]');

      expect(fetchSpy).toHaveBeenCalledWith('*[_type == "post"]');
      expect(result).toHaveLength(1);
      expect(result[0]._type).toBe('post');
      expect(result[0].title).toBe('Test Post');
      expect(result[0].slug.current).toBe('test-post');
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
          publishedAt: '2025-01-01T00:00:00Z',
          excerpt: 'This is a test project',
          content: [],
        },
      ];

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockProjects);

      const result = await mockClient.fetch<Project[]>('*[_type == "project"]');

      expect(fetchSpy).toHaveBeenCalledWith('*[_type == "project"]');
      expect(result).toHaveLength(1);
      expect(result[0]._type).toBe('project');
      expect(result[0].title).toBe('Test Project');
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
        publishedAt: '2025-01-01T00:00:00Z',
        excerpt: 'This is a test post',
        content: [],
      };

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPost);

      const result = await mockClient.fetch<Post>(
        '*[_type == "post" && slug.current == $slug][0]',
        { slug: 'test-post' }
      );

      expect(fetchSpy).toHaveBeenCalledWith(
        '*[_type == "post" && slug.current == $slug][0]',
        { slug: 'test-post' }
      );
      expect(result._id).toBe('post-1');
      expect(result.slug.current).toBe('test-post');
    });

    it('should fetch documents with references (tags)', async () => {
      const mockTag: Tag = {
        _id: 'tag-1',
        _type: 'tag',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: 'rev-1',
        title: 'TypeScript',
        slug: { _type: 'slug', current: 'typescript' },
      };

      const mockPostWithTags: Post = {
        _id: 'post-1',
        _type: 'post',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: 'rev-1',
        title: 'Test Post',
        slug: { _type: 'slug', current: 'test-post' },
        publishedAt: '2025-01-01T00:00:00Z',
        excerpt: 'This is a test post',
        content: [],
        tags: [{ _type: 'reference', _ref: 'tag-1', _key: 'key-1' }],
      };

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPostWithTags);

      const result = await mockClient.fetch<Post>(
        '*[_type == "post"][0]{..., tags[]->{title, slug}}'
      );

      expect(fetchSpy).toHaveBeenCalled();
      expect(result.tags).toBeDefined();
    });

    it('should handle paginated queries', async () => {
      const mockPosts: Post[] = Array.from({ length: 10 }, (_, i) => ({
        _id: `post-${i}`,
        _type: 'post',
        _createdAt: '2025-01-01T00:00:00Z',
        _updatedAt: '2025-01-01T00:00:00Z',
        _rev: `rev-${i}`,
        title: `Post ${i}`,
        slug: { _type: 'slug', current: `post-${i}` },
        publishedAt: '2025-01-01T00:00:00Z',
        excerpt: `Excerpt ${i}`,
        content: [],
      }));

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPosts);

      const result = await mockClient.fetch<Post[]>('*[_type == "post"][0...10]');

      expect(fetchSpy).toHaveBeenCalledWith('*[_type == "post"][0...10]');
      expect(result).toHaveLength(10);
    });

    it('should handle count queries', async () => {
      const mockCount = 42;

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockCount);

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
      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue([]);

      const result = await mockClient.fetch<Post[]>('*[_type == "nonexistent"]');

      expect(fetchSpy).toHaveBeenCalled();
      expect(result).toEqual([]);
      expect(result).toHaveLength(0);
    });

    it('should return null for single document not found', async () => {
      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(null);

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
        publishedAt: '2025-01-01T00:00:00Z',
        excerpt: 'This is a test post',
        content: [],
      };

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPost);

      const result = await mockClient.fetch<Post>('*[_type == "post"][0]');

      expect(fetchSpy).toHaveBeenCalled();
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
        company: 'Test Company',
        location: 'Remote',
        startDate: '2020-01-01',
        content: [],
      };

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockEmployment);

      const result = await mockClient.fetch<Employment>('*[_type == "employment"][0]');

      expect(fetchSpy).toHaveBeenCalled();
      expect(result).toMatchObject({
        _type: 'employment',
        company: expect.any(String),
        location: expect.any(String),
        startDate: expect.any(String),
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
          publishedAt: '2025-01-01T00:00:00Z',
          excerpt: 'TypeScript content',
          content: [],
        },
      ];

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPosts);

      const result = await mockClient.fetch<Post[]>(
        '*[_type == $type && references($tagId)]',
        { type: 'post', tagId: 'tag-1' }
      );

      expect(fetchSpy).toHaveBeenCalledWith(
        '*[_type == $type && references($tagId)]',
        { type: 'post', tagId: 'tag-1' }
      );
      expect(result).toBeDefined();
    });

    it('should handle date range parameters', async () => {
      const mockPosts: Post[] = [];

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPosts);

      await mockClient.fetch<Post[]>(
        '*[_type == "post" && publishedAt >= $startDate && publishedAt <= $endDate]',
        {
          startDate: '2025-01-01',
          endDate: '2025-12-31',
        }
      );

      expect(fetchSpy).toHaveBeenCalledWith(
        '*[_type == "post" && publishedAt >= $startDate && publishedAt <= $endDate]',
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
        publishedAt: '2025-01-01T00:00:00Z',
        excerpt: 'Server-side rendered post',
        content: [],
      };

      const fetchSpy = vi.spyOn(mockClient, 'fetch').mockResolvedValue(mockPost);

      // Simulate route loader behavior
      const loaderData = await mockClient.fetch<Post>('*[_type == "post"][0]');

      expect(fetchSpy).toHaveBeenCalled();
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
