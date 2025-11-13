import { createFileRoute } from '@tanstack/solid-router';
import { For, Show } from 'solid-js';
import { fetchSanityContent } from '@/lib/sanity-server';
import { listEmployments } from '@/lib/sanity-queries';
import type { Employment } from '@/lib/sanity-types';

export const Route = createFileRoute('/experience/')({
  component: ExperiencePage,
  loader: async () => {
    const employments = await fetchSanityContent<Employment[]>(listEmployments());
    return { employments };
  },
});

function ExperiencePage() {
  const data = Route.useLoaderData();

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'Present';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
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
    <div class="min-h-screen bg-gray-900 py-12 px-4 sm:px-6 lg:px-8" data-testid="experience-page">
      <div class="max-w-6xl mx-auto">
        <header class="mb-12">
          <h1 class="text-5xl font-bold text-white mb-4">Professional Experience</h1>
          <p class="text-xl text-gray-300">
            My journey through various roles in engineering and technology
          </p>
        </header>

        <Show
          when={data().employments && data().employments.length > 0}
          fallback={
            <div class="text-center py-12 bg-gray-800 rounded-lg" data-testid="experience-empty">
              <p class="text-gray-400 text-lg">No experience records found</p>
            </div>
          }
        >
          <div class="relative" data-testid="experience-list">
            {/* Timeline Line */}
            <div class="absolute left-8 top-0 bottom-0 w-0.5 bg-gray-700 hidden md:block"></div>

            <div class="space-y-8">
              <For each={data().employments}>
                      {(employment) => (
                        <article
                          class="relative bg-gray-800 rounded-lg border border-gray-700 hover:border-green-500 transition-all duration-300 overflow-hidden"
                          data-testid={`employment-item-${employment.slug?.current}`}
                        >
                          <a
                            href={`/experience/${employment.slug?.current}`}
                            class="block p-6 hover:bg-gray-750 transition-colors"
                          >
                            <div class="flex items-start gap-6">
                              {/* Timeline Dot */}
                              <div class="hidden md:block absolute left-6 w-5 h-5 bg-green-500 rounded-full border-4 border-gray-900 -ml-2"></div>

                              {/* Company Icon Placeholder */}
                              <div class="w-16 h-16 bg-gradient-to-br from-green-500 to-green-700 rounded-full flex items-center justify-center text-white text-2xl font-bold flex-shrink-0 md:ml-12">
                                {employment.title?.[0]?.toUpperCase() || 'E'}
                              </div>

                              {/* Content */}
                              <div class="flex-grow">
                                <h2 class="text-2xl font-bold text-white mb-2 hover:text-green-400 transition-colors">
                                  {employment.title || 'Untitled Position'}
                                </h2>

                                <div class="flex flex-wrap items-center gap-4 text-sm text-gray-400 mb-4">
                                  <span class="flex items-center gap-1 font-medium">
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
                                        d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                                      />
                                    </svg>
                                    Company Name
                                  </span>

                                  <span class="flex items-center gap-1">
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
                                    {formatDate(employment.createdAt)} - {formatDate(employment.endDate)}
                                  </span>

                                  <Show when={employment.createdAt}>
                                    <span class="text-gray-500">
                                      ({calculateDuration(employment.createdAt, employment.endDate)})
                                    </span>
                                  </Show>
                                </div>

                                <Show when={employment.description}>
                                  <p class="text-gray-300 mb-4 line-clamp-3">
                                    {employment.description}
                                  </p>
                                </Show>

                                <Show when={employment.tags && employment.tags.length > 0}>
                                  <div class="flex flex-wrap gap-2">
                                    <For each={employment.tags?.slice(0, 5)}>
                                      {(tag) => (
                                        <Show when={typeof tag !== 'string' && tag._ref}>
                                          <span class="px-2 py-1 bg-green-900 text-green-300 rounded text-xs">
                                            {tag._ref}
                                          </span>
                                        </Show>
                                      )}
                                    </For>
                                  </div>
                                </Show>
                              </div>
                            </div>
                          </a>
              </article>
            )}
          </For>
        </div>
      </div>
    </Show>
      </div>
    </div>
  );
}
