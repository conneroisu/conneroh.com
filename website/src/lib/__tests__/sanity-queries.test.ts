import { describe, it, expect } from 'vitest';
import {
  listPosts,
  getPostBySlug,
  listProjects,
  getProjectBySlug,
  listTags,
  getTagBySlug,
  listEmployments,
  getEmploymentBySlug,
  getPostsByTag,
  getProjectsByTag,
  searchPosts,
  withPagination,
  getPaginationParams,
  withDateRange,
  countPosts,
  countProjects,
  countPostsByTag,
  getRelatedPosts,
  validateSlug,
  validateId,
  validateSearchTerm,
} from '../sanity-queries';
import { SanityValidationError } from '../sanity-errors';

describe('sanity-queries', () => {
  describe('list queries', () => {
    it('should generate listPosts query', () => {
      const query = listPosts();
      expect(query).toContain('*[_type == "post"');
      expect(query).toContain('order(publishedAt desc)');
      expect(query).toContain('!(_id in path("drafts.**"))'); // Should exclude drafts
    });

    it('should generate listProjects query', () => {
      const query = listProjects();
      expect(query).toContain('*[_type == "project"');
      expect(query).toContain('order(_createdAt desc)');
    });

    it('should generate listTags query', () => {
      const query = listTags();
      expect(query).toContain('*[_type == "tag"');
      expect(query).toContain('order(name asc)');
    });

    it('should generate listEmployments query', () => {
      const query = listEmployments();
      expect(query).toContain('*[_type == "employment"');
      expect(query).toContain('order(startDate desc)');
    });
  });

  describe('single document queries', () => {
    it('should generate getPostBySlug query', () => {
      const query = getPostBySlug();
      expect(query).toContain('*[_type == "post"');
      expect(query).toContain('slug.current == $slug');
      expect(query).toContain('[0]');
    });

    it('should generate getProjectBySlug query', () => {
      const query = getProjectBySlug();
      expect(query).toContain('*[_type == "project"');
      expect(query).toContain('slug.current == $slug');
    });

    it('should generate getTagBySlug query', () => {
      const query = getTagBySlug();
      expect(query).toContain('*[_type == "tag"');
      expect(query).toContain('slug.current == $slug');
    });

    it('should generate getEmploymentBySlug query', () => {
      const query = getEmploymentBySlug();
      expect(query).toContain('*[_type == "employment"');
      expect(query).toContain('slug.current == $slug');
    });
  });

  describe('filtered queries', () => {
    it('should generate getPostsByTag query', () => {
      const query = getPostsByTag();
      expect(query).toContain('*[_type == "post"');
      expect(query).toContain('references');
      expect(query).toContain('$tagSlug');
    });

    it('should generate getProjectsByTag query', () => {
      const query = getProjectsByTag();
      expect(query).toContain('*[_type == "project"');
      expect(query).toContain('references');
      expect(query).toContain('$tagSlug');
    });

    it('should generate searchPosts query', () => {
      const query = searchPosts();
      expect(query).toContain('*[_type == "post"');
      expect(query).toContain('title match $searchTerm');
      expect(query).toContain('excerpt match $searchTerm');
    });
  });

  describe('pagination utilities', () => {
    it('should add pagination to query', () => {
      const baseQuery = '*[_type == "post"]';
      const paginatedQuery = withPagination(baseQuery, 0, 10);

      expect(paginatedQuery).toBe('*[_type == "post"][$offset...$endIndex]');
    });

    it('should validate offset is non-negative', () => {
      expect(() => withPagination('*[_type == "post"]', -1, 10)).toThrow(
        SanityValidationError
      );
    });

    it('should validate limit is positive', () => {
      expect(() => withPagination('*[_type == "post"]', 0, 0)).toThrow(
        SanityValidationError
      );
      expect(() => withPagination('*[_type == "post"]', 0, -5)).toThrow(
        SanityValidationError
      );
    });

    it('should validate limit does not exceed 100', () => {
      expect(() => withPagination('*[_type == "post"]', 0, 101)).toThrow(
        SanityValidationError
      );
    });

    it('should generate correct pagination params', () => {
      const params = getPaginationParams(0, 10);
      expect(params).toEqual({ offset: 0, endIndex: 10 });

      const params2 = getPaginationParams(2, 10);
      expect(params2).toEqual({ offset: 20, endIndex: 30 });
    });
  });

  describe('date range filter', () => {
    it('should add date range filter to query', () => {
      const baseQuery = '*[_type == "post"]';
      const filteredQuery = withDateRange(
        baseQuery,
        'publishedAt',
        '2024-01-01',
        '2024-12-31'
      );

      expect(filteredQuery).toContain('publishedAt >= "2024-01-01"');
      expect(filteredQuery).toContain('publishedAt <= "2024-12-31"');
    });

    it('should throw error for invalid query format', () => {
      expect(() => withDateRange('invalid query', 'date', '2024-01-01', '2024-12-31')).toThrow(
        SanityValidationError
      );
    });
  });

  describe('count queries', () => {
    it('should generate countPosts query', () => {
      const query = countPosts();
      expect(query).toBe('count(*[_type == "post" && !(_id in path("drafts.**"))])');
    });

    it('should generate countProjects query', () => {
      const query = countProjects();
      expect(query).toBe('count(*[_type == "project" && !(_id in path("drafts.**"))])');
    });

    it('should generate countPostsByTag query', () => {
      const query = countPostsByTag();
      expect(query).toContain('count(*[_type == "post"');
      expect(query).toContain('references');
      expect(query).toContain('$tagSlug');
    });
  });

  describe('related content queries', () => {
    it('should generate getRelatedPosts query', () => {
      const query = getRelatedPosts();
      expect(query).toContain('*[_type == "post"');
      expect(query).toContain('_id != $postId');
      expect(query).toContain('count((tags[]->_id)');
      expect(query).toContain('[0...$limit]');
    });
  });

  describe('validation helpers', () => {
    describe('validateSlug', () => {
      it('should accept valid slugs', () => {
        expect(() => validateSlug('hello-world')).not.toThrow();
        expect(() => validateSlug('test_slug')).not.toThrow();
        expect(() => validateSlug('slug-123')).not.toThrow();
      });

      it('should reject non-string values', () => {
        expect(() => validateSlug(123)).toThrow(SanityValidationError);
        expect(() => validateSlug(null)).toThrow(SanityValidationError);
        expect(() => validateSlug(undefined)).toThrow(SanityValidationError);
      });

      it('should reject empty strings', () => {
        expect(() => validateSlug('')).toThrow(SanityValidationError);
      });

      it('should reject invalid characters', () => {
        expect(() => validateSlug('hello world')).toThrow(SanityValidationError);
        expect(() => validateSlug('hello/world')).toThrow(SanityValidationError);
        expect(() => validateSlug('hello@world')).toThrow(SanityValidationError);
      });
    });

    describe('validateId', () => {
      it('should accept valid IDs', () => {
        expect(() => validateId('abc123')).not.toThrow();
        expect(() => validateId('document-id-123')).not.toThrow();
      });

      it('should reject non-string values', () => {
        expect(() => validateId(123)).toThrow(SanityValidationError);
        expect(() => validateId(null)).toThrow(SanityValidationError);
      });

      it('should reject empty strings', () => {
        expect(() => validateId('')).toThrow(SanityValidationError);
      });
    });

    describe('validateSearchTerm', () => {
      it('should accept valid search terms', () => {
        expect(() => validateSearchTerm('test')).not.toThrow();
        expect(() => validateSearchTerm('hello world')).not.toThrow();
      });

      it('should reject non-string values', () => {
        expect(() => validateSearchTerm(123)).toThrow(SanityValidationError);
      });

      it('should reject empty strings', () => {
        expect(() => validateSearchTerm('')).toThrow(SanityValidationError);
      });

      it('should reject terms that are too long', () => {
        const longTerm = 'a'.repeat(201);
        expect(() => validateSearchTerm(longTerm)).toThrow(SanityValidationError);
      });
    });
  });

  describe('draft filtering', () => {
    it('should exclude drafts from all queries', () => {
      expect(listPosts()).toContain('!(_id in path("drafts.**"))');
      expect(listProjects()).toContain('!(_id in path("drafts.**"))');
      expect(listTags()).toContain('!(_id in path("drafts.**"))');
      expect(listEmployments()).toContain('!(_id in path("drafts.**"))');
      expect(getPostBySlug()).toContain('!(_id in path("drafts.**"))');
    });
  });
});
