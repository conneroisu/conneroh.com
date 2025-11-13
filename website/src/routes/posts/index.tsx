import { createFileRoute } from '@tanstack/solid-router';
import { For, Show, Suspense } from 'solid-js';
import { getPosts } from '@/lib/sanity-server-funcs';
import type { Post } from '@/lib/sanity-types';

export const Route = createFileRoute('/posts/')({
  ssr: true,
  component: PostsPage,
  loader: async () => {
    // Fetch all posts (pagination can be added later)
    const posts = await getPosts();

    return {
      posts,
      page: 0,
      pageSize: 10,
    };
  },
});

function PostsPage() {
  const data = Route.useLoaderData();

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'No date';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <div class="min-h-screen bg-gray-900 py-12 px-4 sm:px-6 lg:px-8" data-testid="posts-page">
      <div class="max-w-6xl mx-auto">
        <header class="mb-12">
          <h1 class="text-5xl font-bold text-white mb-4">Blog Posts</h1>
          <p class="text-xl text-gray-300">
            Explore my thoughts, tutorials, and technical deep-dives
          </p>
        </header>

        <Show
          when={data().posts && data().posts.length > 0}
          fallback={
            <div class="text-center py-12 bg-gray-800 rounded-lg" data-testid="posts-empty">
              <p class="text-gray-400 text-lg">No posts found</p>
              <p class="text-gray-500 mt-2">Check back later for new content!</p>
            </div>
          }
        >
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8" data-testid="posts-list">
            <For each={data().posts}>
              {(post) => (
                <article
                  class="bg-gray-800 rounded-lg border border-gray-700 hover:border-green-500 transition-all duration-300 overflow-hidden"
                  data-testid={`post-item-${post.slug?.current}`}
                >
                  <a
                    href={`/posts/${post.slug?.current}`}
                    class="block p-6 hover:bg-gray-750 transition-colors"
                  >
                    <h2 class="text-2xl font-bold text-white mb-3 hover:text-green-400 transition-colors">
                      {post.title || 'Untitled Post'}
                    </h2>

                    <Show when={post.description}>
                      <p class="text-gray-300 mb-4 line-clamp-3">
                        {post.description}
                      </p>
                    </Show>

                    <div class="flex flex-wrap items-center gap-4 text-sm text-gray-400">
                      <time
                        datetime={post.createdAt}
                        class="flex items-center gap-1"
                      >
                        <svg
                          class="w-4 h-4"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                          />
                        </svg>
                        {formatDate(post.createdAt)}
                      </time>

                      <Show when={post.tags && post.tags.length > 0}>
                        <div class="flex items-center gap-2">
                          <svg
                            class="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                            />
                          </svg>
                          <div class="flex gap-2">
                            <For each={post.tags?.slice(0, 3)}>
                              {(tag) => (
                                <Show when={typeof tag !== 'string' && tag._ref}>
                                  <span class="text-green-400">
                                    #{tag._ref}
                                  </span>
                                </Show>
                              )}
                            </For>
                          </div>
                        </div>
                      </Show>
                    </div>
                  </a>
                </article>
              )}
            </For>
          </div>

          {/* Pagination Controls */}
          <nav class="mt-12 flex justify-center items-center gap-2 text-sm" data-testid="posts-pagination">
            <Show when={data().page > 0}>
              <a
                href={`/posts?page=${data().page - 1}`}
                class="px-4 py-2 bg-gray-700 text-gray-300 rounded hover:bg-gray-600 transition-colors"
                data-testid="posts-prev-page"
              >
                Previous
              </a>
            </Show>

            <span class="px-4 py-2 bg-gray-800 text-gray-300 rounded">
              Page {data().page + 1}
            </span>

            <Show when={data().posts.length === data().pageSize}>
              <a
                href={`/posts?page=${data().page + 1}`}
                class="px-4 py-2 bg-gray-700 text-gray-300 rounded hover:bg-gray-600 transition-colors"
                data-testid="posts-next-page"
              >
                Next
              </a>
            </Show>
          </nav>
        </Show>
      </div>
    </div>
  );
}
