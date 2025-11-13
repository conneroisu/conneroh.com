import { createFileRoute, Navigate } from '@tanstack/solid-router';
import { Show, For } from 'solid-js';
import { fetchSanityContent } from '@/lib/sanity-server';
import { getProjectBySlug } from '@/lib/sanity-queries';
import type { Project } from '@/lib/sanity-types';
import { PortableText } from '@/components/sanity/PortableText';

export const Route = createFileRoute('/projects/$slug')({
  component: ProjectDetailPage,
  loader: async ({ params }) => {
    const project = await fetchSanityContent<Project | null>(
      getProjectBySlug(),
      { slug: params.slug }
    );

    return { project };
  },
});

function ProjectDetailPage() {
  const data = Route.useLoaderData();

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'No date';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
    });
  };

  return (
    <Show
      when={data().project}
      fallback={<Navigate to="/projects" />}
    >
      <article class="min-h-screen bg-gray-900" data-testid="project-detail">
        {/* Main Container */}
        <div class="max-w-4xl mx-auto px-4 py-12">
          {/* Banner Image (if exists) */}
          <Show when={data().project!.bannerImage}>
            <img
              src={data().project!.bannerImage}
              alt={data().project!.title || 'Project banner'}
              class="w-full h-64 md:h-96 object-cover rounded-lg mb-8"
            />
          </Show>

          {/* Header Card */}
          <div class="bg-gray-800 rounded-lg p-6 mb-8 shadow-lg">
            <h1 class="text-4xl font-bold text-white mb-4" data-testid="project-title">
              {data().project!.title || 'Untitled Project'}
            </h1>

            <h2 class="text-2xl font-semibold text-white mb-2">About this project</h2>
            <div class="h-1 bg-green-500 mt-2 w-16 mb-6"></div>

            <Show when={data().project!.description}>
              <p class="text-gray-300 whitespace-pre-line mb-6" data-testid="project-description">
                {data().project!.description}
              </p>
            </Show>

            {/* Info Section */}
            <div class="space-y-3">
              {/* Tags */}
              <Show when={(data().project!.tags?.length || 0) > 0}>
                <div class="flex flex-wrap items-center gap-2">
                  <span class="text-gray-400 text-sm pr-2">Technologies:</span>
                  <For each={data().project!.tags}>
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
              <Show when={data().project!.createdAt}>
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
                    datetime={data().project!.createdAt}
                    class="text-sm text-gray-400"
                    data-testid="project-date"
                  >
                    {formatDate(data().project!.createdAt)}
                  </time>
                </div>
              </Show>
            </div>
          </div>

          {/* Content Card */}
          <div class="bg-gray-800 rounded-lg p-6 mb-8 shadow-lg" data-testid="project-content">
            <article class="text-gray-300 leading-relaxed max-w-none prose prose-invert">
              <PortableText content={data().project!.content} />
            </article>
          </div>

          {/* Related Posts Section */}
          <Show when={(data().project!.posts?.length || 0) > 0}>
            <div class="pt-8 mt-12 border-t border-gray-700">
              <h2 class="text-2xl font-bold text-white mb-4">Related Posts</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4" data-testid="project-posts">
                <For each={data().project!.posts}>
                  {(post) => (
                    <Show when={typeof post !== 'string' && post._ref}>
                      <a
                        href={`/posts/${post._ref}`}
                        class="bg-gray-800 rounded-lg p-4 hover:bg-gray-700 transition-colors"
                      >
                        <span class="text-green-400 hover:text-green-300 font-medium">
                          {post._ref}
                        </span>
                      </a>
                    </Show>
                  )}
                </For>
              </div>
            </div>
          </Show>

          {/* Back to Projects */}
          <div class="mt-12 text-center">
            <a
              href="/projects"
              class="inline-flex items-center gap-2 px-6 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
              data-testid="back-to-projects"
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
              Back to All Projects
            </a>
          </div>
        </div>
      </article>
    </Show>
  );
}
