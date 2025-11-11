import { SanityValidationError } from '@/lib/sanity-errors';
import type { Post, Project, Tag, Employment } from './sanity-types';

/**
 * Type definitions for query results
 */
export type PostQueryResult = Post;
export type ProjectQueryResult = Project;
export type TagQueryResult = Tag;
export type EmploymentQueryResult = Employment;

/**
 * Common projection for post documents
 */
const POST_PROJECTION = `
  _id,
  _type,
  _createdAt,
  _updatedAt,
  title,
  slug,
  excerpt,
  content,
  publishedAt,
  author->,
  tags[]->
`;

/**
 * Common projection for project documents
 */
const PROJECT_PROJECTION = `
  _id,
  _type,
  _createdAt,
  _updatedAt,
  title,
  slug,
  description,
  content,
  url,
  github,
  image,
  tags[]->
`;

/**
 * Common projection for tag documents
 */
const TAG_PROJECTION = `
  _id,
  _type,
  _createdAt,
  _updatedAt,
  name,
  slug,
  description
`;

/**
 * Common projection for employment documents
 */
const EMPLOYMENT_PROJECTION = `
  _id,
  _type,
  _createdAt,
  _updatedAt,
  company,
  slug,
  position,
  description,
  startDate,
  endDate,
  current,
  skills[]->
`;

/**
 * List all published posts ordered by publish date (newest first)
 *
 * @example
 * ```typescript
 * const posts = await client.fetch(listPosts());
 * ```
 *
 * @returns GROQ query string
 */
export function listPosts(): string {
  return `*[_type == "post" && !(_id in path("drafts.**"))] | order(publishedAt desc) {
    ${POST_PROJECTION}
  }`;
}

/**
 * Get a single post by slug
 *
 * @example
 * ```typescript
 * const post = await client.fetch(getPostBySlug(), { slug: 'hello-world' });
 * ```
 *
 * @returns GROQ query string
 */
export function getPostBySlug(): string {
  return `*[_type == "post" && slug.current == $slug && !(_id in path("drafts.**"))][0] {
    ${POST_PROJECTION}
  }`;
}

/**
 * List all projects ordered by creation date (newest first)
 *
 * @example
 * ```typescript
 * const projects = await client.fetch(listProjects());
 * ```
 *
 * @returns GROQ query string
 */
export function listProjects(): string {
  return `*[_type == "project" && !(_id in path("drafts.**"))] | order(_createdAt desc) {
    ${PROJECT_PROJECTION}
  }`;
}

/**
 * Get a single project by slug
 *
 * @example
 * ```typescript
 * const project = await client.fetch(getProjectBySlug(), { slug: 'my-project' });
 * ```
 *
 * @returns GROQ query string
 */
export function getProjectBySlug(): string {
  return `*[_type == "project" && slug.current == $slug && !(_id in path("drafts.**"))][0] {
    ${PROJECT_PROJECTION}
  }`;
}

/**
 * List all tags ordered by name
 *
 * @example
 * ```typescript
 * const tags = await client.fetch(listTags());
 * ```
 *
 * @returns GROQ query string
 */
export function listTags(): string {
  return `*[_type == "tag" && !(_id in path("drafts.**"))] | order(name asc) {
    ${TAG_PROJECTION}
  }`;
}

/**
 * Get a single tag by slug
 *
 * @example
 * ```typescript
 * const tag = await client.fetch(getTagBySlug(), { slug: 'typescript' });
 * ```
 *
 * @returns GROQ query string
 */
export function getTagBySlug(): string {
  return `*[_type == "tag" && slug.current == $slug && !(_id in path("drafts.**"))][0] {
    ${TAG_PROJECTION}
  }`;
}

/**
 * List all employment records ordered by start date (newest first)
 *
 * @example
 * ```typescript
 * const employments = await client.fetch(listEmployments());
 * ```
 *
 * @returns GROQ query string
 */
export function listEmployments(): string {
  return `*[_type == "employment" && !(_id in path("drafts.**"))] | order(startDate desc) {
    ${EMPLOYMENT_PROJECTION}
  }`;
}

/**
 * Get a single employment record by slug
 *
 * @example
 * ```typescript
 * const employment = await client.fetch(getEmploymentBySlug(), { slug: 'acme-corp' });
 * ```
 *
 * @returns GROQ query string
 */
export function getEmploymentBySlug(): string {
  return `*[_type == "employment" && slug.current == $slug && !(_id in path("drafts.**"))][0] {
    ${EMPLOYMENT_PROJECTION}
  }`;
}

/**
 * Get posts by tag slug
 *
 * @example
 * ```typescript
 * const posts = await client.fetch(getPostsByTag(), { tagSlug: 'typescript' });
 * ```
 *
 * @returns GROQ query string
 */
export function getPostsByTag(): string {
  return `*[_type == "post" && !(_id in path("drafts.**")) && references(*[_type == "tag" && slug.current == $tagSlug][0]._id)] | order(publishedAt desc) {
    ${POST_PROJECTION}
  }`;
}

/**
 * Get projects by tag slug
 *
 * @example
 * ```typescript
 * const projects = await client.fetch(getProjectsByTag(), { tagSlug: 'react' });
 * ```
 *
 * @returns GROQ query string
 */
export function getProjectsByTag(): string {
  return `*[_type == "project" && !(_id in path("drafts.**")) && references(*[_type == "tag" && slug.current == $tagSlug][0]._id)] | order(_createdAt desc) {
    ${PROJECT_PROJECTION}
  }`;
}

/**
 * Search posts by title or excerpt
 *
 * @example
 * ```typescript
 * const results = await client.fetch(searchPosts(), { searchTerm: 'solid.js' });
 * ```
 *
 * @returns GROQ query string
 */
export function searchPosts(): string {
  return `*[_type == "post" && !(_id in path("drafts.**")) && (title match $searchTerm || excerpt match $searchTerm)] | order(publishedAt desc) {
    ${POST_PROJECTION}
  }`;
}

/**
 * Pagination utilities
 */

/**
 * Validate pagination parameters
 *
 * @param offset - Starting index
 * @param limit - Number of items
 * @throws SanityValidationError if parameters are invalid
 */
function validatePagination(offset: number, limit: number): void {
  if (offset < 0) {
    throw new SanityValidationError('Offset must be non-negative', { offset });
  }
  if (limit <= 0) {
    throw new SanityValidationError('Limit must be positive', { limit });
  }
  if (limit > 100) {
    throw new SanityValidationError('Limit cannot exceed 100', { limit });
  }
}

/**
 * Add pagination to a GROQ query
 *
 * This utility appends array slice syntax to paginate results.
 * Validates that offset >= 0 and 0 < limit <= 100.
 *
 * @example
 * ```typescript
 * const query = withPagination(listPosts(), 0, 10);
 * const posts = await client.fetch(query);
 * ```
 *
 * @param query - Base GROQ query string
 * @param offset - Starting index (0-based)
 * @param limit - Number of items to fetch
 * @returns Paginated query string
 * @throws SanityValidationError if pagination parameters are invalid
 */
export function withPagination(query: string, offset: number, limit: number): string {
  validatePagination(offset, limit);

  const endIndex = offset + limit;
  return `${query}[$offset...$endIndex]`;
}

/**
 * Create pagination parameters object for use with withPagination
 *
 * @example
 * ```typescript
 * const query = withPagination(listPosts(), page * pageSize, pageSize);
 * const params = getPaginationParams(page, pageSize);
 * const posts = await client.fetch(query, params);
 * ```
 *
 * @param page - Page number (0-indexed)
 * @param pageSize - Items per page
 * @returns Object with offset and endIndex for query params
 */
export function getPaginationParams(
  page: number,
  pageSize: number
): {
  offset: number;
  endIndex: number;
} {
  validatePagination(page * pageSize, pageSize);
  const offset = page * pageSize;
  const endIndex = offset + pageSize;
  return {
    offset,
    endIndex,
  };
}

/**
 * Filter queries
 */

/**
 * Add date range filter to a query
 *
 * @example
 * ```typescript
 * const query = withDateRange(
 *   listPosts(),
 *   'publishedAt',
 *   '2024-01-01',
 *   '2024-12-31'
 * );
 * ```
 *
 * @param query - Base GROQ query
 * @param field - Date field to filter on
 * @param startDate - Start date (ISO string)
 * @param endDate - End date (ISO string)
 * @returns Filtered query string
 */
export function withDateRange(
  query: string,
  field: string,
  startDate: string,
  endDate: string
): string {
  // Extract the filter portion (everything before the first ']')
  const filterEnd = query.indexOf(']');
  if (filterEnd === -1) {
    throw new SanityValidationError('Invalid query format: no closing bracket found', { query });
  }

  const beforeFilter = query.substring(0, filterEnd);
  const afterFilter = query.substring(filterEnd);

  return `${beforeFilter} && ${field} >= "${startDate}" && ${field} <= "${endDate}"${afterFilter}`;
}

/**
 * Count queries
 */

/**
 * Get count of all posts
 *
 * @example
 * ```typescript
 * const count = await client.fetch(countPosts());
 * ```
 *
 * @returns GROQ count query string
 */
export function countPosts(): string {
  return `count(*[_type == "post" && !(_id in path("drafts.**"))])`;
}

/**
 * Get count of all projects
 *
 * @example
 * ```typescript
 * const count = await client.fetch(countProjects());
 * ```
 *
 * @returns GROQ count query string
 */
export function countProjects(): string {
  return `count(*[_type == "project" && !(_id in path("drafts.**"))])`;
}

/**
 * Get count of posts by tag
 *
 * @example
 * ```typescript
 * const count = await client.fetch(countPostsByTag(), { tagSlug: 'typescript' });
 * ```
 *
 * @returns GROQ count query string
 */
export function countPostsByTag(): string {
  return `count(*[_type == "post" && !(_id in path("drafts.**")) && references(*[_type == "tag" && slug.current == $tagSlug][0]._id)])`;
}

/**
 * Related content queries
 */

/**
 * Get related posts by shared tags
 *
 * @example
 * ```typescript
 * const related = await client.fetch(getRelatedPosts(), { postId: '123', limit: 3 });
 * ```
 *
 * @returns GROQ query string
 */
export function getRelatedPosts(): string {
  return `*[_type == "post" && !(_id in path("drafts.**")) && _id != $postId && count((tags[]->_id)[@ in *[_id == $postId][0].tags[]->_id]) > 0] | order(publishedAt desc) [0...$limit] {
    ${POST_PROJECTION}
  }`;
}

/**
 * Query parameter validation
 */

/**
 * Validate slug parameter
 *
 * @param slug - Slug value to validate
 * @throws SanityValidationError if slug is invalid
 */
export function validateSlug(slug: unknown): asserts slug is string {
  if (typeof slug !== 'string') {
    throw new SanityValidationError('Slug must be a string', { slug });
  }
  if (slug.length === 0) {
    throw new SanityValidationError('Slug cannot be empty', { slug });
  }
  // Basic slug format validation (alphanumeric, hyphens, underscores)
  if (!/^[a-z0-9_-]+$/i.test(slug)) {
    throw new SanityValidationError(
      'Slug must contain only alphanumeric characters, hyphens, and underscores',
      { slug }
    );
  }
}

/**
 * Validate document ID parameter
 *
 * @param id - ID value to validate
 * @throws SanityValidationError if ID is invalid
 */
export function validateId(id: unknown): asserts id is string {
  if (typeof id !== 'string') {
    throw new SanityValidationError('Document ID must be a string', { id });
  }
  if (id.length === 0) {
    throw new SanityValidationError('Document ID cannot be empty', { id });
  }
}

/**
 * Validate search term parameter
 *
 * @param searchTerm - Search term to validate
 * @throws SanityValidationError if search term is invalid
 */
export function validateSearchTerm(searchTerm: unknown): asserts searchTerm is string {
  if (typeof searchTerm !== 'string') {
    throw new SanityValidationError('Search term must be a string', { searchTerm });
  }
  if (searchTerm.length === 0) {
    throw new SanityValidationError('Search term cannot be empty', { searchTerm });
  }
  if (searchTerm.length > 200) {
    throw new SanityValidationError('Search term is too long (max 200 characters)', {
      searchTerm,
    });
  }
}
