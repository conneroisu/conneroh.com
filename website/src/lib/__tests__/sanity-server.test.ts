import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  fetchSanityContent,
  fetchDocumentById,
  fetchDocumentBySlug,
  fetchPaginatedDocuments,
  fetchDocumentCount,
  fetchDocumentsByIds,
  fetchAllDocuments,
  documentExists,
} from '../sanity-server';
import type { SanityClient } from '@sanity/client';

describe('sanity-server', () => {
  // Create fresh mock client for each test suite
  let mockClient: Partial<SanityClient>;

  beforeEach(() => {
    mockClient = {
      fetch: vi.fn(),
    };
  });

  describe('fetchSanityContent', () => {
    it('should fetch content with query and params', async () => {
      const mockData = [{ id: 1, title: 'Test' }];
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockData);

      const result = await fetchSanityContent(
        '*[_type == $type]',
        { type: 'post' },
        { client: mockClient as SanityClient }
      );

      expect(result).toEqual(mockData);
      expect(mockClient.fetch).toHaveBeenCalledWith('*[_type == $type]', { type: 'post' });
    });

    it('should handle errors', async () => {
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Network error'));

      await expect(
        fetchSanityContent('*[_type == "post"]', {}, { client: mockClient as SanityClient, retry: false, logErrors: false })
      ).rejects.toThrow();
    });

    it('should retry on failure when retry is enabled', async () => {
      (mockClient.fetch as ReturnType<typeof vi.fn>)
        .mockRejectedValueOnce(new Error('Network error'))
        .mockResolvedValueOnce([{ id: 1 }]);

      const result = await fetchSanityContent(
        '*[_type == "post"]',
        {},
        { client: mockClient as SanityClient, retry: { maxAttempts: 2, baseDelay: 10 }, logErrors: false }
      );

      expect(result).toEqual([{ id: 1 }]);
      expect(mockClient.fetch).toHaveBeenCalledTimes(2);
    });
  });

  describe('fetchDocumentById', () => {
    it('should fetch a document by type and id', async () => {
      const mockDoc = { _id: '123', _type: 'post', title: 'Test' };
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockDoc);

      const result = await fetchDocumentById('post', '123', { client: mockClient as SanityClient });

      expect(result).toEqual(mockDoc);
      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == $type && _id == $id][0]',
        { type: 'post', id: '123' }
      );
    });

    it('should return null if document not found', async () => {
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(null);

      const result = await fetchDocumentById('post', '999', { client: mockClient as SanityClient });

      expect(result).toBeNull();
    });
  });

  describe('fetchDocumentBySlug', () => {
    it('should fetch a document by type and slug', async () => {
      const mockDoc = { slug: { current: 'test-post' }, title: 'Test' };
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockDoc);

      const result = await fetchDocumentBySlug('post', 'test-post', {
        client: mockClient as SanityClient,
      });

      expect(result).toEqual(mockDoc);
      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == $type && slug.current == $slug][0]',
        { type: 'post', slug: 'test-post' }
      );
    });
  });

  describe('fetchPaginatedDocuments', () => {
    it('should fetch paginated results', async () => {
      const mockDocs = [{ id: 1 }, { id: 2 }];
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockDocs);

      const result = await fetchPaginatedDocuments(
        '*[_type == "post"]',
        0,
        10,
        { client: mockClient as SanityClient }
      );

      expect(result).toEqual(mockDocs);
      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == "post"][$offset...$limit]',
        { offset: 0, limit: 10 }
      );
    });

    it('should calculate correct offset for page 2', async () => {
      const mockDocs = [{ id: 11 }, { id: 12 }];
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockDocs);

      const result = await fetchPaginatedDocuments(
        '*[_type == "post"]',
        2,
        10,
        { client: mockClient as SanityClient }
      );

      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == "post"][$offset...$limit]',
        { offset: 20, limit: 30 }
      );
    });
  });

  describe('fetchDocumentCount', () => {
    it('should fetch count of documents', async () => {
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(42);

      const result = await fetchDocumentCount('*[_type == "post"]', {}, {
        client: mockClient as SanityClient,
      });

      expect(result).toBe(42);
      expect(mockClient.fetch).toHaveBeenCalledWith('count(*[_type == "post"])', {});
    });
  });

  describe('fetchDocumentsByIds', () => {
    it('should fetch multiple documents by ids', async () => {
      const mockDocs = [{ _id: '1' }, { _id: '2' }];
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockDocs);

      const result = await fetchDocumentsByIds('post', ['1', '2'], {
        client: mockClient as SanityClient,
      });

      expect(result).toEqual(mockDocs);
      expect(mockClient.fetch).toHaveBeenCalledWith('*[_type == $type && _id in $ids]', {
        type: 'post',
        ids: ['1', '2'],
      });
    });

    it('should return empty array for empty ids', async () => {
      const result = await fetchDocumentsByIds('post', [], { client: mockClient as SanityClient });

      expect(result).toEqual([]);
      expect(mockClient.fetch).not.toHaveBeenCalled();
    });
  });

  describe('fetchAllDocuments', () => {
    it('should fetch all documents of a type', async () => {
      const mockDocs = [{ id: 1 }, { id: 2 }, { id: 3 }];
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockDocs);

      const result = await fetchAllDocuments('tag', undefined, {
        client: mockClient as SanityClient,
      });

      expect(result).toEqual(mockDocs);
      expect(mockClient.fetch).toHaveBeenCalledWith('*[_type == $type]', { type: 'tag' });
    });

    it('should support orderBy parameter', async () => {
      const mockDocs = [{ id: 3 }, { id: 2 }, { id: 1 }];
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockDocs);

      const result = await fetchAllDocuments('post', '_createdAt desc', {
        client: mockClient as SanityClient,
      });

      expect(mockClient.fetch).toHaveBeenCalledWith(
        '*[_type == $type] | order(_createdAt desc)',
        { type: 'post' }
      );
    });
  });

  describe('documentExists', () => {
    it('should return true if document exists', async () => {
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(true);

      const result = await documentExists('post', 'slug.current', 'test-post', {
        client: mockClient as SanityClient,
      });

      expect(result).toBe(true);
      expect(mockClient.fetch).toHaveBeenCalledWith(
        'count(*[_type == $type && slug.current == $value]) > 0',
        { type: 'post', value: 'test-post' }
      );
    });

    it('should return false if document does not exist', async () => {
      (mockClient.fetch as ReturnType<typeof vi.fn>).mockResolvedValue(false);

      const result = await documentExists('post', '_id', '999', {
        client: mockClient as SanityClient,
      });

      expect(result).toBe(false);
    });
  });
});
