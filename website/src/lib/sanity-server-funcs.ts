import { createServerFn } from '@tanstack/solid-start';
import { fetchSanityContent } from './sanity-server';
import {
  listPosts,
  getPostBySlug as getPostBySlugQuery,
  listProjects,
  listEmployments,
  getPostsByTag as getPostsByTagQuery,
  getRelatedPosts as getRelatedPostsQuery,
  getTagBySlug as getTagBySlugQuery,
  getProjectsByTag as getProjectsByTagQuery,
} from './sanity-queries';
import type { Post, Project, Employment, Tag } from './sanity-types';

/**
 * Server function to fetch all published posts
 *
 * @returns Promise<Post[]> - Array of published posts ordered by publish date (newest first)
 *
 * @example
 * ```typescript
 * import { getPosts } from '@/lib/sanity-server-funcs';
 *
 * const posts = await getPosts();
 * ```
 */
export const getPosts = createServerFn({
  method: 'GET',
}).handler(async (): Promise<Post[]> => {
  return await fetchSanityContent<Post[]>(listPosts());
});

/**
 * Server function to fetch a single post by slug
 *
 * @param slug - Post slug
 * @returns Promise<Post | null> - Post document or null if not found
 *
 * @example
 * ```typescript
 * import { getPostBySlug } from '@/lib/sanity-server-funcs';
 *
 * const post = await getPostBySlug('hello-world');
 * ```
 */
export const getPostBySlug = createServerFn({
  method: 'GET',
}).handler(async (slug: string): Promise<Post | null> => {
  const slugValue = String(slug);
  return await fetchSanityContent<Post | null>(getPostBySlugQuery(), { slug: slugValue });
});

/**
 * Server function to fetch all projects
 *
 * @returns Promise<Project[]> - Array of projects ordered by creation date (newest first)
 *
 * @example
 * ```typescript
 * import { getProjects } from '@/lib/sanity-server-funcs';
 *
 * const projects = await getProjects();
 * ```
 */
export const getProjects = createServerFn({
  method: 'GET',
}).handler(async (): Promise<Project[]> => {
  return await fetchSanityContent<Project[]>(listProjects());
});

/**
 * Server function to fetch all employment records
 *
 * @returns Promise<Employment[]> - Array of employment records ordered by start date (newest first)
 *
 * @example
 * ```typescript
 * import { getExperience } from '@/lib/sanity-server-funcs';
 *
 * const experience = await getExperience();
 * ```
 */
export const getExperience = createServerFn({
  method: 'GET',
}).handler(async (): Promise<Employment[]> => {
  return await fetchSanityContent<Employment[]>(listEmployments());
});

/**
 * Server function to fetch posts filtered by tag slug
 *
 * @param tag - Tag slug to filter by
 * @returns Promise<Post[]> - Array of posts with the specified tag, ordered by publish date (newest first)
 *
 * @example
 * ```typescript
 * import { getPostsByTag } from '@/lib/sanity-server-funcs';
 *
 * const posts = await getPostsByTag('typescript');
 * ```
 */
export const getPostsByTag = createServerFn({
  method: 'GET',
}).handler(async (tag: string): Promise<Post[]> => {
  const tagValue = String(tag);
  return await fetchSanityContent<Post[]>(getPostsByTagQuery(), { tagSlug: tagValue });
});

/**
 * Server function to fetch posts related to a given post
 *
 * @param postId - Post ID to find related posts for
 * @param limit - Maximum number of related posts to return (default: 3)
 * @returns Promise<Post[]> - Array of related posts based on shared tags
 *
 * @example
 * ```typescript
 * import { getRelatedPosts } from '@/lib/sanity-server-funcs';
 *
 * const relatedPosts = await getRelatedPosts('post-id-123', 3);
 * ```
 */
export const getRelatedPosts = createServerFn({
  method: 'GET',
}).handler(async (postId: string, limit: number = 3): Promise<Post[]> => {
  const postIdValue = String(postId);
  const limitValue = Number(limit);
  return await fetchSanityContent<Post[]>(getRelatedPostsQuery(), { postId: postIdValue, limit: limitValue });
});

/**
 * Server function to fetch a single tag by slug
 *
 * @param slug - Tag slug
 * @returns Promise<Tag | null> - Tag document or null if not found
 *
 * @example
 * ```typescript
 * import { getTagBySlug } from '@/lib/sanity-server-funcs';
 *
 * const tag = await getTagBySlug('typescript');
 * ```
 */
export const getTagBySlug = createServerFn({
  method: 'GET',
}).handler(async (slug: string): Promise<Tag | null> => {
  const slugValue = String(slug);
  return await fetchSanityContent<Tag | null>(getTagBySlugQuery(), { slug: slugValue });
});

/**
 * Server function to fetch projects filtered by tag
 *
 * @param tag - Tag slug to filter by
 * @returns Promise<Project[]> - Array of projects with the specified tag, ordered by creation date (newest first)
 *
 * @example
 * ```typescript
 * import { getProjectsByTag } from '@/lib/sanity-server-funcs';
 *
 * const projects = await getProjectsByTag('react');
 * ```
 */
export const getProjectsByTag = createServerFn({
  method: 'GET',
}).handler(async (tag: string): Promise<Project[]> => {
  const tagValue = String(tag);
  return await fetchSanityContent<Project[]>(getProjectsByTagQuery(), { tagSlug: tagValue });
});
