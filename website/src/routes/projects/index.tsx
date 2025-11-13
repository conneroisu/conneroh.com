import { createFileRoute } from '@tanstack/solid-router';
import { For, Show } from 'solid-js';
import { fetchSanityContent } from '@/lib/sanity-server';
import { listProjects } from '@/lib/sanity-queries';
import type { Project } from '@/lib/sanity-types';

export const Route = createFileRoute('/projects/')({
  component: ProjectsPage,
  loader: async () => {
    const projects = await fetchSanityContent<Project[]>(listProjects());
    return { projects };
  },
});

function ProjectsPage() {
  const data = Route.useLoaderData();

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'No date';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
    });
  };

  return (
    <div class="min-h-screen bg-gray-900 py-12 px-4 sm:px-6 lg:px-8" data-testid="projects-page">
      <div class="max-w-6xl mx-auto">
        <header class="mb-12">
          <h1 class="text-5xl font-bold text-white mb-4">Projects</h1>
          <p class="text-xl text-gray-300">
            A collection of projects I've built, contributed to, or experimented with
          </p>
        </header>

        <Show
          when={data().projects && data().projects.length > 0}
          fallback={
            <div class="text-center py-12 bg-gray-800 rounded-lg" data-testid="projects-empty">
              <p class="text-gray-400 text-lg">No projects found</p>
              <p class="text-gray-500 mt-2">Check back later for new projects!</p>
            </div>
          }
        >
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8" data-testid="projects-list">
            <For each={data().projects}>
                    {(project) => (
                      <article
                        class="bg-gray-800 rounded-lg border border-gray-700 hover:border-green-500 transition-all duration-300 overflow-hidden flex flex-col"
                        data-testid={`project-item-${project.slug?.current}`}
                      >
                        <a
                          href={`/projects/${project.slug?.current}`}
                          class="block p-6 flex-grow hover:bg-gray-750 transition-colors"
                        >
                          <h2 class="text-2xl font-bold text-white mb-3 hover:text-green-400 transition-colors">
                            {project.title || 'Untitled Project'}
                          </h2>

                          <Show when={project.description}>
                            <p class="text-gray-300 mb-4 line-clamp-3">
                              {project.description}
                            </p>
                          </Show>

                          <div class="flex flex-wrap items-center gap-4 text-sm text-gray-400 mb-4">
                            <time
                              datetime={project.createdAt}
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
                              {formatDate(project.createdAt)}
                            </time>

                            <Show when={project.tags && project.tags.length > 0}>
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
                                  <For each={project.tags?.slice(0, 3)}>
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
    </Show>
      </div>
    </div>
  );
}
