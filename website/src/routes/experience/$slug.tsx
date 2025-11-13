import { createFileRoute, Navigate } from '@tanstack/solid-router';
import { Show, For } from 'solid-js';
import { fetchSanityContent } from '@/lib/sanity-server';
import { getEmploymentBySlug } from '@/lib/sanity-queries';
import type { Employment } from '@/lib/sanity-types';
import { PortableText } from '@/components/sanity/PortableText';

export const Route = createFileRoute('/experience/$slug')({
  component: ExperienceDetailPage,
  loader: async ({ params }) => {
    const employment = await fetchSanityContent<Employment | null>(
      getEmploymentBySlug(),
      { slug: params.slug }
    );

    return { employment };
  },
});

function ExperienceDetailPage() {
  const data = Route.useLoaderData();

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'Present';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
    });
  };

  const calculateDuration = (startDate: string | undefined, endDate: string | undefined) => {
    if (!startDate) return '';

    const start = new Date(startDate);
    const end = endDate ? new Date(endDate) : new Date();

    const months = (end.getFullYear() - start.getFullYear()) * 12 + (end.getMonth() - start.getMonth());
    const years = Math.floor(months / 12);
    const remainingMonths = months % 12;

    if (years === 0) {
      return `${remainingMonths} month${remainingMonths !== 1 ? 's' : ''}`;
    }
    if (remainingMonths === 0) {
      return `${years} year${years !== 1 ? 's' : ''}`;
    }
    return `${years} year${years !== 1 ? 's' : ''}, ${remainingMonths} month${remainingMonths !== 1 ? 's' : ''}`;
  };

  return (
    <Show
      when={data().employment}
      fallback={<Navigate to="/experience" />}
    >
      <article class="min-h-screen bg-gray-900" data-testid="employment-detail">
        {/* Main Container */}
        <div class="max-w-4xl mx-auto px-4 py-12">
          {/* Banner Image (if exists) */}
          <Show when={data().employment!.bannerImage}>
            <img
              src={data().employment!.bannerImage}
              alt={data().employment!.title || 'Employment banner'}
              class="w-full h-64 md:h-96 object-cover rounded-lg mb-8"
            />
          </Show>

          {/* Header Card */}
          <div class="bg-gray-800 rounded-lg p-6 mb-8 shadow-lg">
            <h1 class="text-4xl font-bold text-white mb-4" data-testid="employment-title">
              {data().employment!.title || 'Untitled Position'}
            </h1>

            <h2 class="text-2xl font-semibold text-white mb-2">About this employment</h2>
            <div class="h-1 bg-green-500 mt-2 w-16 mb-6"></div>

            <Show when={data().employment!.description}>
              <p class="text-gray-300 whitespace-pre-line mb-6" data-testid="employment-description">
                {data().employment!.description}
              </p>
            </Show>

            {/* Info Section */}
            <div class="space-y-3">
              {/* Tags */}
              <Show when={(data().employment!.tags?.length || 0) > 0}>
                <div class="flex flex-wrap items-center gap-2">
                  <span class="text-gray-400 text-sm pr-2">Skills & Technologies:</span>
                  <For each={data().employment!.tags}>
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

              {/* Start Date */}
              <Show when={data().employment!.createdAt}>
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
                  <span class="text-sm text-gray-400" data-testid="employment-dates">
                    Started {formatDate(data().employment!.createdAt)}
                  </span>
                </div>
              </Show>

              {/* End Date (if exists) */}
              <Show when={data().employment!.endDate}>
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
                  <span class="text-sm text-gray-400">
                    Ended {formatDate(data().employment!.endDate)}
                  </span>
                </div>
              </Show>

              {/* Duration */}
              <Show when={data().employment!.createdAt}>
                <div class="text-sm text-gray-400">
                  Duration: {calculateDuration(data().employment!.createdAt, data().employment!.endDate)}
                </div>
              </Show>
            </div>
          </div>

          {/* Content Card */}
          <div class="bg-gray-800 rounded-lg p-6 mb-8 shadow-lg" data-testid="employment-content">
            <article class="text-gray-300 leading-relaxed max-w-none prose prose-invert">
              <PortableText content={data().employment!.content} />
            </article>
          </div>

          {/* Related Projects Section */}
          <Show when={(data().employment!.projects?.length || 0) > 0}>
            <div class="pt-8 mt-12 border-t border-gray-700">
              <h2 class="text-2xl font-bold text-white mb-4">Related Projects</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4" data-testid="employment-projects">
                <For each={data().employment!.projects}>
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
          <Show when={(data().employment!.posts?.length || 0) > 0}>
            <div class="pt-8 mt-12 border-t border-gray-700">
              <h2 class="text-2xl font-bold text-white mb-4">Related Posts</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4" data-testid="employment-posts">
                <For each={data().employment!.posts}>
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

          {/* Related Employments Section (if applicable) */}
          <Show when={(data().employment!.relatedEmployments?.length || 0) > 0}>
            <div class="pt-8 mt-12 border-t border-gray-700">
              <h2 class="text-2xl font-bold text-white mb-4">Related Employments</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <For each={data().employment!.relatedEmployments}>
                  {(employment) => (
                    <Show when={typeof employment !== 'string' && employment._ref}>
                      <a
                        href={`/experience/${employment._ref}`}
                        class="bg-gray-800 rounded-lg p-4 hover:bg-gray-700 transition-colors"
                      >
                        <span class="text-green-400 hover:text-green-300 font-medium">
                          {employment._ref}
                        </span>
                      </a>
                    </Show>
                  )}
                </For>
              </div>
            </div>
          </Show>

          {/* Back to Experience */}
          <div class="mt-12 text-center">
            <a
              href="/experience"
              class="inline-flex items-center gap-2 px-6 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
              data-testid="back-to-experience"
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
              Back to Experience
            </a>
          </div>
        </div>
      </article>
    </Show>
  );
}
