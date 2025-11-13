import { createFileRoute, Navigate } from '@tanstack/solid-router';
import { Show, For } from 'solid-js';
import { getPostBySlug, getRelatedPosts } from '@/lib/sanity-server-funcs';
import type { Post } from '@/lib/sanity-types';
import { PortableText } from '@/components/sanity/PortableText';

export const Route = createFileRoute('/posts/$slug')({
  ssr: true,
  component: PostDetailPage,
  loader: async ({ params }) => {
    const post = await getPostBySlug(params.slug);

    if (!post) {
      return { post: null, relatedPosts: [] };
    }

    // Fetch related posts (limit to 3)
    const relatedPosts = await getRelatedPosts(post._id, 3);

    return {
      post,
      relatedPosts,
    };
  },
});

function PostDetailPage() {
  const data = Route.useLoaderData();

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'No date';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const calculateReadTime = (content: any) => {
    // Simple read time calculation - approximately 200 words per minute
    if (!content || !Array.isArray(content)) return '5 min read';

    const wordCount = content.reduce((count, block) => {
      if (block._type === 'block' && block.children) {
        const text = block.children.map((child: any) => child.text || '').join(' ');
        return count + text.split(/\s+/).length;
      }
      return count;
    }, 0);

    const minutes = Math.ceil(wordCount / 200);
    return `${minutes} min read`;
  };

  return (
    <Show
      when={data().post}
      fallback={<Navigate to="/posts" />}
    >
      <article class="min-h-screen bg-gray-900" data-testid="post-detail">
        {/* Main Container */}
        <div class="max-w-4xl mx-auto px-4 py-12">
          {/* Banner Image (if exists) */}
          <Show when={data().post!.bannerImage}>
            <img
              src={data().post!.bannerImage}
              alt={data().post!.title || 'Post banner'}
              class="w-full h-64 md:h-96 object-cover rounded-lg mb-8"
            />
          </Show>

          {/* Header Card */}
          <div class="bg-gray-800 rounded-lg p-6 mb-8 shadow-lg">
            <h1 class="text-4xl font-bold text-white mb-4" data-testid="post-title">
              {data().post!.title || 'Untitled Post'}
            </h1>

            <h2 class="text-2xl font-semibold text-white mb-2">About this post</h2>
            <div class="h-1 bg-green-500 mt-2 w-16 mb-6"></div>

            <Show when={data().post!.description}>
              <p class="text-gray-300 whitespace-pre-line mb-6" data-testid="post-description">
                {data().post!.description}
              </p>
            </Show>

            {/* Info Section */}
            <div class="space-y-3">
              {/* Tags */}
              <Show when={(data().post!.tags?.length || 0) > 0}>
                <div class="flex flex-wrap items-center gap-2">
                  <span class="text-gray-400 text-sm pr-2">Tags:</span>
                  <For each={data().post!.tags}>
                    {(tag) => (
                      <Show when={typeof tag !== 'string' && tag._ref}>
                        <a
                          href={`/tags/${tag._ref}`}
                          class="text-gray-400 hover:text-green-400 transition-colors text-sm"
                        >
                          #{tag._ref}
                        </a>
                      </Show>
                    )}
                  </For>
                </div>
              </Show>

              {/* Created Date */}
              <div class="flex items-center gap-2">
                <svg
                  class="w-4 h-4 text-gray-400"
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
                <time
                  datetime={data().post!.createdAt}
                  class="text-sm text-gray-400"
                  data-testid="post-date"
                >
                  {formatDate(data().post!.createdAt)}
                </time>
              </div>

              {/* Read Time */}
              <div class="text-sm text-gray-400">
                {calculateReadTime(data().post!.content)}
              </div>
            </div>
          </div>

          {/* Content Card */}
          <div class="bg-gray-800 rounded-lg p-6 mb-8 shadow-lg" data-testid="post-content">
            <article class="text-gray-300 leading-relaxed max-w-none prose prose-invert">
              <PortableText content={data().post!.content} />
            </article>
          </div>

          {/* Related Projects Section */}
          <Show when={(data().post!.projects?.length || 0) > 0}>
            <div class="pt-8 mt-12 border-t border-gray-700">
              <h2 class="text-2xl font-bold text-white mb-4">Related Projects</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4" data-testid="post-projects">
                <For each={data().post!.projects}>
                  {(project) => (
                    <Show when={typeof project !== 'string' && project._ref}>
                      <a
                        href={`/projects/${project._ref}`}
                        class="bg-gray-800 rounded-lg p-4 hover:bg-gray-700 transition-colors"
                      >
                        <span class="text-green-400 hover:text-green-300 font-medium">
                          {project._ref}
                        </span>
                      </a>
                    </Show>
                  )}
                </For>
              </div>
            </div>
          </Show>

          {/* Related Posts Section */}
          <Show when={data().relatedPosts && data().relatedPosts.length > 0}>
            <div class="pt-8 mt-12 border-t border-gray-700">
              <h2 class="text-2xl font-bold text-white mb-4">Related Posts</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4" data-testid="related-posts">
                <For each={data().relatedPosts}>
                  {(relatedPost) => (
                    <a
                      href={`/posts/${relatedPost.slug?.current}`}
                      class="bg-gray-800 rounded-lg p-4 hover:bg-gray-700 transition-colors"
                    >
                      <h4 class="font-semibold text-white mb-2">
                        {relatedPost.title || 'Untitled Post'}
                      </h4>
                      <Show when={relatedPost.description}>
                        <p class="text-gray-300 text-sm line-clamp-2">
                          {relatedPost.description}
                        </p>
                      </Show>
                    </a>
                  )}
                </For>
              </div>
            </div>
          </Show>

          {/* Back to Posts */}
          <div class="mt-12 text-center">
            <a
              href="/posts"
              class="inline-flex items-center gap-2 px-6 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
              data-testid="back-to-posts"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M10 19l-7-7m0 0l7-7m-7 7h18"
                />
              </svg>
              Back to All Posts
            </a>
          </div>
        </div>
      </article>
    </Show>
  );
}
