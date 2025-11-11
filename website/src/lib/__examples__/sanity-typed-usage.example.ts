/**
 * Example usage of Sanity CMS with generated TypeScript types
 *
 * This file demonstrates how to use the type-safe Sanity integration
 * across different parts of the application.
 */

import { createSanityQuery, createSanityDocumentBySlug } from '@/hooks/use-sanity-content';
import type { Post, Project, Tag, Employment } from '@/lib/sanity-client';
import { listPosts, getPostBySlug, listProjects } from '@/lib/sanity-queries';
import { fetchSanityContent, fetchDocumentBySlug } from '@/lib/sanity-server';
import { getSanityClient } from '@/lib/sanity-client';

/**
 * Example 1: Using typed queries in Solid.js components (client-side)
 */
export function PostListComponent() {
  // createSanityQuery now returns properly typed data
  const [posts, { refetch }] = createSanityQuery<Post[]>(
    listPosts(),
    {}
  );

  return () => {
    const postData = posts();
    if (!postData) return <div>Loading...</div>;

    // TypeScript knows the shape of Post
    return (
      <div>
        {postData.map((post) => (
          <article key={post._id}>
            <h2>{post.title}</h2>
            <p>{post.description}</p>
            {/* TypeScript autocompletes all Post fields */}
            <time>{post.createdAt}</time>
          </article>
        ))}
      </div>
    );
  };
}

/**
 * Example 2: Using typed document queries by slug
 */
export function PostDetailComponent(props: { slug: string }) {
  const [post] = createSanityDocumentBySlug<Post>(
    'post',
    () => props.slug
  );

  return () => {
    const postData = post();
    if (!postData) return <div>Loading...</div>;

    // Full type safety on Post fields
    return (
      <article>
        <h1>{postData.title}</h1>
        <div>{postData.description}</div>
        {/* BlockContent is also typed */}
        {postData.content && <div>Content blocks: {postData.content.length}</div>}
      </article>
    );
  };
}

/**
 * Example 3: Server-side typed queries (TanStack Start server functions)
 */
export async function serverGetPosts(): Promise<Post[]> {
  // fetchSanityContent returns typed results
  const posts = await fetchSanityContent<Post[]>(listPosts());

  // TypeScript knows posts is Post[]
  return posts.filter(post => post.title !== undefined);
}

/**
 * Example 4: Server-side typed document fetching
 */
export async function serverGetPostBySlug(slug: string): Promise<Post | null> {
  const post = await fetchDocumentBySlug<Post>('post', slug);

  // TypeScript knows post is Post | null
  if (post && post.title) {
    console.log(`Fetched post: ${post.title}`);
  }

  return post;
}

/**
 * Example 5: Using typed queries with the client directly
 */
export async function directClientQuery() {
  const client = getSanityClient();

  // Explicitly type the query result
  const projects = await client.fetch<Project[]>(listProjects());

  // TypeScript knows the shape of Project
  projects.forEach(project => {
    console.log(project.title, project.description);
  });
}

/**
 * Example 6: Working with related content (references)
 */
export async function getPostWithTags(slug: string): Promise<void> {
  const post = await fetchDocumentBySlug<Post>('post', slug);

  if (!post || !post.tags) return;

  // Tags are references in the Post type
  // You'd need to expand them in your query projection
  // or fetch them separately
  const client = getSanityClient();

  const expandedPost = await client.fetch<Post & { tags: Tag[] }>(
    `*[_type == "post" && slug.current == $slug][0] {
      ...,
      tags[]->
    }`,
    { slug }
  );

  if (expandedPost?.tags) {
    // Now tags are fully typed Tag objects
    expandedPost.tags.forEach(tag => {
      console.log(tag.title, tag.description);
    });
  }
}

/**
 * Example 7: Type-safe employment records
 */
export async function getEmploymentHistory(): Promise<Employment[]> {
  const employments = await fetchSanityContent<Employment[]>(
    `*[_type == "employment"] | order(createdAt desc)`
  );

  // Filter out entries without required fields
  return employments.filter(emp =>
    emp.title !== undefined &&
    emp.description !== undefined
  );
}

/**
 * Example 8: Using union types with AllSanitySchemaTypes
 */
import type { AllSanitySchemaTypes } from '@/lib/sanity-client';

export async function getAnyDocument(id: string): Promise<AllSanitySchemaTypes | null> {
  const client = getSanityClient();
  const doc = await client.fetch<AllSanitySchemaTypes | null>(
    `*[_id == $id][0]`,
    { id }
  );

  // Type narrowing based on _type
  if (doc?._type === 'post') {
    // TypeScript knows this is a Post
    console.log(doc.title);
  } else if (doc?._type === 'project') {
    // TypeScript knows this is a Project
    console.log(doc.description);
  }

  return doc;
}
