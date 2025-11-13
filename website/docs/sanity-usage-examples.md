# Sanity CMS Usage Examples

Practical, copy-paste examples for common Sanity CMS operations in TanStack Start + Solid.js.

## Table of Contents

- [Client-Side Queries](#client-side-queries)
- [Server-Side Queries](#server-side-queries)
- [Dynamic Routing with Slugs](#dynamic-routing-with-slugs)
- [Loading States and Suspense](#loading-states-and-suspense)
- [Pagination](#pagination)
- [Relationships and References](#relationships-and-references)
- [Error Handling](#error-handling)
- [Caching](#caching)
- [Advanced Patterns](#advanced-patterns)

## Client-Side Queries

### Fetching a List of Posts

**Use Case**: Display all published blog posts on the client side with reactive updates.

```typescript
// src/routes/posts/index.tsx
import { For, Show } from 'solid-js';
import { createSanityQuery } from '@/hooks/use-sanity-content';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

export default function PostsPage() {
  const [posts] = createSanityQuery<Post[]>(listPosts(), {});

  return (
    <div class="posts-page">
      <h1>Blog Posts</h1>

      <Show when={!posts.loading} fallback={<div>Loading posts...</div>}>
        <Show when={posts()} fallback={<div>No posts found</div>}>
          {(postList) => (
            <div class="posts-grid">
              <For each={postList()}>
                {(post) => (
                  <article class="post-card">
                    <h2>{post.title}</h2>
                    <p>{post.excerpt}</p>
                    <a href={`/posts/${post.slug?.current}`}>
                      Read more →
                    </a>
                  </article>
                )}
              </For>
            </div>
          )}
        </Show>
      </Show>
    </div>
  );
}
```

### Fetching Projects

**Use Case**: Display portfolio projects with tags.

```typescript
// src/routes/projects/index.tsx
import { For, Show } from 'solid-js';
import { createSanityQuery } from '@/hooks/use-sanity-content';
import { listProjects } from '@/lib/sanity-queries';
import type { Project } from '@/lib/sanity-client';

export default function ProjectsPage() {
  const [projects] = createSanityQuery<Project[]>(listProjects(), {});

  return (
    <div class="projects-page">
      <h1>Projects</h1>

      <Show when={projects()}>
        {(projectList) => (
          <div class="projects-grid">
            <For each={projectList()}>
              {(project) => (
                <div class="project-card">
                  <h3>{project.title}</h3>
                  <p>{project.description}</p>

                  <Show when={project.tags && project.tags.length > 0}>
                    <div class="tags">
                      <For each={project.tags}>
                        {(tag) => <span class="tag">{tag.name}</span>}
                      </For>
                    </div>
                  </Show>

                  <div class="links">
                    <Show when={project.url}>
                      <a href={project.url} target="_blank">
                        View Project
                      </a>
                    </Show>
                    <Show when={project.github}>
                      <a href={project.github} target="_blank">
                        View Source
                      </a>
                    </Show>
                  </div>
                </div>
              )}
            </For>
          </div>
        )}
      </Show>
    </div>
  );
}
```

## Server-Side Queries

### Server Function with Type Safety

**Use Case**: Fetch posts on the server for better performance and SEO.

```typescript
// src/routes/posts/index.tsx
import { createAsync } from '@tanstack/solid-start';
import { createServerFn } from '@tanstack/solid-start/server';
import { fetchSanityContent } from '@/lib/sanity-server';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

// Server function - runs on server only
const getPostsFn = createServerFn('GET', async () => {
  const posts = await fetchSanityContent<Post[]>(listPosts(), {});
  return posts;
});

export default function PostsPage() {
  // Data loaded on server, streamed to client
  const posts = createAsync(() => getPostsFn());

  return (
    <div class="posts-page">
      <h1>Blog Posts</h1>

      <Show when={posts()}>
        {(postList) => (
          <For each={postList()}>
            {(post) => (
              <article>
                <h2>{post.title}</h2>
                <p>{post.excerpt}</p>
              </article>
            )}
          </For>
        )}
      </Show>
    </div>
  );
}
```

### Route Loader Pattern

**Use Case**: Load data in route loader for instant navigation.

```typescript
// src/routes/posts/index.tsx
import { useLoaderData } from '@tanstack/solid-start';
import { createServerFn } from '@tanstack/solid-start/server';
import { fetchSanityContent } from '@/lib/sanity-server';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

// Loader function
export const loader = createServerFn('GET', async () => {
  return fetchSanityContent<Post[]>(listPosts(), {});
});

export default function PostsPage() {
  const posts = useLoaderData<typeof loader>();

  return (
    <div class="posts-page">
      <For each={posts()}>
        {(post) => (
          <article>
            <h2>{post.title}</h2>
          </article>
        )}
      </For>
    </div>
  );
}
```

## Dynamic Routing with Slugs

### Single Post Page

**Use Case**: Display individual post by slug with full content.

```typescript
// src/routes/posts/[slug].tsx
import { useParams } from '@tanstack/solid-start';
import { createAsync } from '@tanstack/solid-start';
import { createServerFn } from '@tanstack/solid-start/server';
import { fetchDocumentBySlug } from '@/lib/sanity-server';
import type { Post } from '@/lib/sanity-client';
import { Show } from 'solid-js';

// Server function to fetch post by slug
const getPostBySlugFn = createServerFn('GET', async (slug: string) => {
  const post = await fetchDocumentBySlug<Post>('post', slug);

  if (!post) {
    throw new Response('Post not found', { status: 404 });
  }

  return post;
});

export default function PostPage() {
  const params = useParams();
  const post = createAsync(() => getPostBySlugFn(params.slug));

  return (
    <div class="post-page">
      <Show when={post()}>
        {(p) => (
          <article>
            <header>
              <h1>{p().title}</h1>
              <Show when={p().publishedAt}>
                <time>{new Date(p().publishedAt!).toLocaleDateString()}</time>
              </Show>
            </header>

            <Show when={p().excerpt}>
              <p class="excerpt">{p().excerpt}</p>
            </Show>

            <Show when={p().tags && p().tags!.length > 0}>
              <div class="tags">
                <For each={p().tags}>
                  {(tag) => (
                    <a href={`/tags/${tag.slug?.current}`} class="tag">
                      {tag.name}
                    </a>
                  )}
                </For>
              </div>
            </Show>

            {/* Render content using Portable Text renderer */}
            <div class="content">
              {/* TODO: Add Portable Text renderer */}
              {JSON.stringify(p().content)}
            </div>
          </article>
        )}
      </Show>
    </div>
  );
}
```

### Client-Side Dynamic Route

**Use Case**: Fetch post data on client side with slug from route params.

```typescript
// src/routes/posts/[slug].tsx
import { useParams } from '@tanstack/solid-start';
import { createSanityDocumentBySlug } from '@/hooks/use-sanity-content';
import type { Post } from '@/lib/sanity-client';
import { Show } from 'solid-js';

export default function PostPage() {
  const params = useParams();

  // Reactive query based on slug
  const [post] = createSanityDocumentBySlug<Post>(
    'post',
    () => params.slug
  );

  return (
    <div class="post-page">
      <Show
        when={!post.loading}
        fallback={<div>Loading post...</div>}
      >
        <Show when={post()} fallback={<div>Post not found</div>}>
          {(p) => (
            <article>
              <h1>{p().title}</h1>
              <p>{p().excerpt}</p>
            </article>
          )}
        </Show>
      </Show>
    </div>
  );
}
```

## Loading States and Suspense

### Using Solid.js Suspense

**Use Case**: Show loading fallback while content loads.

```typescript
// src/routes/posts/index.tsx
import { Suspense, For } from 'solid-js';
import { createSanityQuery } from '@/hooks/use-sanity-content';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

function PostList() {
  const [posts] = createSanityQuery<Post[]>(listPosts(), {});

  return (
    <For each={posts()}>
      {(post) => <PostCard post={post} />}
    </For>
  );
}

export default function PostsPage() {
  return (
    <div class="posts-page">
      <h1>Blog Posts</h1>

      <Suspense fallback={
        <div class="loading-skeleton">
          <div class="skeleton-card" />
          <div class="skeleton-card" />
          <div class="skeleton-card" />
        </div>
      }>
        <PostList />
      </Suspense>
    </div>
  );
}
```

### Custom Loading States

**Use Case**: Show different UI based on loading/error states.

```typescript
// src/routes/posts/index.tsx
import { Show, Match, Switch } from 'solid-js';
import { createSanityQuery } from '@/hooks/use-sanity-content';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

export default function PostsPage() {
  const [posts] = createSanityQuery<Post[]>(listPosts(), {});

  return (
    <div class="posts-page">
      <h1>Blog Posts</h1>

      <Switch>
        <Match when={posts.loading}>
          <div class="loading">
            <Spinner />
            <p>Loading posts...</p>
          </div>
        </Match>

        <Match when={posts.error}>
          <div class="error">
            <p>Failed to load posts</p>
            <button onClick={() => posts.refetch()}>
              Retry
            </button>
          </div>
        </Match>

        <Match when={posts()}>
          <PostGrid posts={posts()!} />
        </Match>
      </Switch>
    </div>
  );
}
```

## Pagination

### Client-Side Pagination

**Use Case**: Paginate large lists on the client side.

```typescript
// src/routes/posts/index.tsx
import { createSignal, For, Show } from 'solid-js';
import { createSanityPaginatedQuery } from '@/hooks/use-sanity-content';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

export default function PostsPage() {
  const [page, setPage] = createSignal(0);
  const pageSize = 10;

  const [posts] = createSanityPaginatedQuery<Post>(
    listPosts(),
    page,
    pageSize
  );

  return (
    <div class="posts-page">
      <h1>Blog Posts</h1>

      <Show when={posts()}>
        {(postList) => (
          <>
            <div class="posts-grid">
              <For each={postList()}>
                {(post) => <PostCard post={post} />}
              </For>
            </div>

            <div class="pagination">
              <button
                onClick={() => setPage(p => Math.max(0, p - 1))}
                disabled={page() === 0}
              >
                Previous
              </button>

              <span>Page {page() + 1}</span>

              <button
                onClick={() => setPage(p => p + 1)}
                disabled={postList().length < pageSize}
              >
                Next
              </button>
            </div>
          </>
        )}
      </Show>
    </div>
  );
}
```

### Server-Side Pagination with Count

**Use Case**: Paginate with total count for better UX.

```typescript
// src/routes/posts/index.tsx
import { createAsync } from '@tanstack/solid-start';
import { createServerFn } from '@tanstack/solid-start/server';
import { useSearchParams } from '@tanstack/solid-start';
import { fetchPaginatedDocuments, fetchDocumentCount } from '@/lib/sanity-server';
import { listPosts, countPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

const getPostsPageFn = createServerFn('GET', async (page: number) => {
  const pageSize = 10;

  const [posts, total] = await Promise.all([
    fetchPaginatedDocuments<Post>(listPosts(), page, pageSize),
    fetchDocumentCount(countPosts()),
  ]);

  return { posts, total, page, pageSize };
});

export default function PostsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = () => parseInt(searchParams.page || '0', 10);

  const data = createAsync(() => getPostsPageFn(page()));

  return (
    <div class="posts-page">
      <Show when={data()}>
        {(d) => (
          <>
            <For each={d().posts}>
              {(post) => <PostCard post={post} />}
            </For>

            <Pagination
              currentPage={d().page}
              totalItems={d().total}
              pageSize={d().pageSize}
              onPageChange={(newPage) => setSearchParams({ page: newPage.toString() })}
            />
          </>
        )}
      </Show>
    </div>
  );
}
```

## Relationships and References

### Fetching Posts with Tags

**Use Case**: Display posts with their referenced tags.

```typescript
// src/routes/posts/index.tsx
import { For, Show } from 'solid-js';
import { createSanityQuery } from '@/hooks/use-sanity-content';
import type { Post } from '@/lib/sanity-client';

// GROQ query with tag references expanded
const postsWithTagsQuery = `
  *[_type == "post" && !(_id in path("drafts.**"))] | order(publishedAt desc) {
    _id,
    title,
    slug,
    excerpt,
    publishedAt,
    tags[]-> {
      _id,
      name,
      slug
    }
  }
`;

export default function PostsPage() {
  const [posts] = createSanityQuery<Post[]>(postsWithTagsQuery, {});

  return (
    <div class="posts-page">
      <Show when={posts()}>
        {(postList) => (
          <For each={postList()}>
            {(post) => (
              <article class="post-card">
                <h2>{post.title}</h2>
                <p>{post.excerpt}</p>

                <Show when={post.tags && post.tags.length > 0}>
                  <div class="tags">
                    <For each={post.tags}>
                      {(tag) => (
                        <a href={`/tags/${tag.slug?.current}`}>
                          {tag.name}
                        </a>
                      )}
                    </For>
                  </div>
                </Show>
              </article>
            )}
          </For>
        )}
      </Show>
    </div>
  );
}
```

### Posts by Tag

**Use Case**: Display all posts for a specific tag.

```typescript
// src/routes/tags/[slug].tsx
import { useParams } from '@tanstack/solid-start';
import { createAsync } from '@tanstack/solid-start';
import { createServerFn } from '@tanstack/solid-start/server';
import { fetchSanityContent } from '@/lib/sanity-server';
import { getPostsByTag } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';
import { For, Show } from 'solid-js';

const getPostsByTagFn = createServerFn('GET', async (tagSlug: string) => {
  return fetchSanityContent<Post[]>(
    getPostsByTag(),
    { tagSlug }
  );
});

export default function TagPage() {
  const params = useParams();
  const posts = createAsync(() => getPostsByTagFn(params.slug));

  return (
    <div class="tag-page">
      <h1>Posts tagged with "{params.slug}"</h1>

      <Show when={posts()}>
        {(postList) => (
          <Show when={postList().length > 0} fallback={<p>No posts found</p>}>
            <For each={postList()}>
              {(post) => <PostCard post={post} />}
            </For>
          </Show>
        )}
      </Show>
    </div>
  );
}
```

### Related Posts by Shared Tags

**Use Case**: Show related posts based on common tags.

```typescript
// src/routes/posts/[slug].tsx
import { createAsync } from '@tanstack/solid-start';
import { createServerFn } from '@tanstack/solid-start/server';
import { fetchSanityContent } from '@/lib/sanity-server';
import { getRelatedPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';
import { For, Show } from 'solid-js';

const getRelatedPostsFn = createServerFn('GET', async (postId: string) => {
  return fetchSanityContent<Post[]>(
    getRelatedPosts(),
    { postId, limit: 3 }
  );
});

function RelatedPosts(props: { postId: string }) {
  const related = createAsync(() => getRelatedPostsFn(props.postId));

  return (
    <aside class="related-posts">
      <h3>Related Posts</h3>

      <Show when={related() && related()!.length > 0}>
        <For each={related()}>
          {(post) => (
            <a href={`/posts/${post.slug?.current}`}>
              {post.title}
            </a>
          )}
        </For>
      </Show>
    </aside>
  );
}
```

## Error Handling

### Graceful Error Display

**Use Case**: Show user-friendly errors with retry option.

```typescript
// src/routes/posts/index.tsx
import { Show, ErrorBoundary } from 'solid-js';
import { createSanityQuery } from '@/hooks/use-sanity-content';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

export default function PostsPage() {
  const [posts, { refetch }] = createSanityQuery<Post[]>(
    listPosts(),
    {},
    {
      onError: (error) => {
        console.error('Failed to fetch posts:', error);
      },
    }
  );

  return (
    <div class="posts-page">
      <h1>Blog Posts</h1>

      <ErrorBoundary
        fallback={(err, reset) => (
          <div class="error-container">
            <h2>Something went wrong</h2>
            <p>{err.message}</p>
            <button onClick={() => {
              reset();
              refetch();
            }}>
              Try Again
            </button>
          </div>
        )}
      >
        <Show when={!posts.error} fallback={
          <div class="error">
            <p>Failed to load posts</p>
            <button onClick={() => refetch()}>Retry</button>
          </div>
        }>
          <PostList posts={posts()} />
        </Show>
      </ErrorBoundary>
    </div>
  );
}
```

## Caching

### Using Cached Fetch

**Use Case**: Cache frequently accessed content to reduce API calls.

```typescript
// src/lib/post-api.ts
import { getSanityClient } from '@/lib/sanity-client';
import { cachedFetch, getTtlForQuery } from '@/lib/sanity-cache';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

export async function getCachedPosts(): Promise<Post[]> {
  const client = getSanityClient();
  const query = listPosts();

  return cachedFetch<Post[]>(
    client,
    query,
    {},
    {
      ttl: getTtlForQuery(query), // Auto-determine TTL
    }
  );
}
```

### Manual Cache Invalidation

**Use Case**: Invalidate cache when content is updated.

```typescript
// src/lib/post-mutations.ts
import { createServerClient } from '@/lib/sanity-client';
import { invalidationHelpers, getGlobalCache } from '@/lib/sanity-cache';

export async function createPost(title: string, content: string) {
  const client = createServerClient();

  const newPost = await client.create({
    _type: 'post',
    title,
    content,
    publishedAt: new Date().toISOString(),
  });

  // Invalidate post-related cache entries
  invalidationHelpers.posts();

  return newPost;
}

export async function updatePost(id: string, updates: Partial<Post>) {
  const client = createServerClient();

  const updated = await client.patch(id).set(updates).commit();

  // Invalidate specific post and all lists
  const cache = getGlobalCache();
  cache.invalidate(new RegExp(id));
  invalidationHelpers.posts();

  return updated;
}
```

## Advanced Patterns

### Search with Debouncing

**Use Case**: Search posts with debounced input.

```typescript
// src/routes/search.tsx
import { createSignal, createEffect, on } from 'solid-js';
import { createSanityQuery } from '@/hooks/use-sanity-content';
import { searchPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

export default function SearchPage() {
  const [searchTerm, setSearchTerm] = createSignal('');
  const [debouncedTerm, setDebouncedTerm] = createSignal('');

  // Debounce search term
  createEffect(on(searchTerm, (term) => {
    const timer = setTimeout(() => {
      setDebouncedTerm(term);
    }, 300);

    return () => clearTimeout(timer);
  }));

  const [results] = createSanityQuery<Post[]>(
    searchPosts(),
    () => ({ searchTerm: `${debouncedTerm()}*` }),
    { deferStream: true }
  );

  return (
    <div class="search-page">
      <input
        type="search"
        value={searchTerm()}
        onInput={(e) => setSearchTerm(e.currentTarget.value)}
        placeholder="Search posts..."
      />

      <Show when={debouncedTerm() && results()}>
        {(resultList) => (
          <div class="results">
            <p>{resultList().length} results found</p>
            <For each={resultList()}>
              {(post) => <PostCard post={post} />}
            </For>
          </div>
        )}
      </Show>
    </div>
  );
}
```

### Infinite Scroll

**Use Case**: Load more posts as user scrolls.

```typescript
// src/routes/posts/index.tsx
import { createSignal, onMount, For, Show } from 'solid-js';
import { fetchSanityContent } from '@/lib/sanity-server';
import { listPosts } from '@/lib/sanity-queries';
import type { Post } from '@/lib/sanity-client';

export default function PostsPage() {
  const [posts, setPosts] = createSignal<Post[]>([]);
  const [page, setPage] = createSignal(0);
  const [loading, setLoading] = createSignal(false);
  const [hasMore, setHasMore] = createSignal(true);
  const pageSize = 10;

  const loadMore = async () => {
    if (loading() || !hasMore()) return;

    setLoading(true);
    try {
      const query = `${listPosts()}[${page() * pageSize}...${(page() + 1) * pageSize}]`;
      const newPosts = await fetchSanityContent<Post[]>(query, {});

      if (newPosts.length < pageSize) {
        setHasMore(false);
      }

      setPosts(prev => [...prev, ...newPosts]);
      setPage(p => p + 1);
    } finally {
      setLoading(false);
    }
  };

  onMount(() => {
    loadMore();

    // Intersection Observer for infinite scroll
    const observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        loadMore();
      }
    });

    const sentinel = document.querySelector('#scroll-sentinel');
    if (sentinel) observer.observe(sentinel);

    return () => observer.disconnect();
  });

  return (
    <div class="posts-page">
      <For each={posts()}>
        {(post) => <PostCard post={post} />}
      </For>

      <Show when={hasMore()}>
        <div id="scroll-sentinel" />
      </Show>

      <Show when={loading()}>
        <div class="loading">Loading more...</div>
      </Show>
    </div>
  );
}
```

## Next Steps

- Explore [API Reference](./sanity-api-reference.md) for complete function documentation
- Review [Setup Guide](./sanity-setup.md) for configuration details
- Check [GROQ Documentation](https://www.sanity.io/docs/groq) for advanced queries
- Experiment with [Portable Text](https://www.sanity.io/docs/presenting-block-text) for rich content rendering
