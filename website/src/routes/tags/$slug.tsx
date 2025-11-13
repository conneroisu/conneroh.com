import { createFileRoute, Navigate } from '@tanstack/solid-router';
import { Show, For, createSignal } from 'solid-js';
import { getPostsByTag, getTagBySlug, getProjectsByTag } from '@/lib/sanity-server-funcs';
import type { Tag, Post, Project } from '@/lib/sanity-types';

export const Route = createFileRoute('/tags/$slug')({
  component: TagDetailPage,
  loader: async ({ params }) => {
    const tag = await getTagBySlug(params.slug);

    if (!tag) {
      return { tag: null, posts: [], projects: [] };
    }

    const [posts, projects] = await Promise.all([
      getPostsByTag(params.slug),
      getProjectsByTag(params.slug),
    ]);

    return { tag, posts, projects };
  },
});

function TagDetailPage() {
  const data = Route.useLoaderData();
  const [activeTab, setActiveTab] = createSignal<'posts' | 'projects' | 'tags'>('posts');

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'No date';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <Show
      when={data().tag}
      fallback={<Navigate to="/tags" />}
    >
      <div class="min-h-screen bg-gray-900" data-testid="tag-detail">
        {/* Main Container */}
        <div class="max-w-4xl mx-auto px-4 py-8">
          {/* Header Section */}
          <div class="text-center mb-8">
            <Show when={data().tag!.icon}>
              <div class="text-8xl mb-4">{data().tag!.icon}</div>
            </Show>
            <h1 class="text-4xl font-bold text-white mb-4" data-testid="tag-title">
              {data().tag!.title || 'Untitled Tag'}
            </h1>
            <Show when={data().tag!.description}>
              <p class="text-xl text-gray-300 max-w-2xl mx-auto" data-testid="tag-description">
                {data().tag!.description}
              </p>
            </Show>
          </div>

          {/* Content Card */}
          <Show when={data().tag!.content}>
            <div class="bg-gray-800 rounded-lg p-6 mb-8 shadow-lg">
              <div class="text-gray-300 leading-relaxed">
                {data().tag!.content}
              </div>
            </div>
          </Show>

          {/* Tabs Section */}
          <div class="mb-6">
            <div class="border-b border-gray-700">
              <nav class="flex gap-8" aria-label="Tabs">
                <button
                  onClick={() => setActiveTab('posts')}
                  class={`px-1 sm:text-base text-sm font-medium py-4 border-b-2 transition-colors ${
                    activeTab() === 'posts'
                      ? 'text-green-500 border-green-500'
                      : 'text-gray-400 hover:text-gray-300 border-transparent'
                  }`}
                >
                  Posts
                </button>
                <button
                  onClick={() => setActiveTab('projects')}
                  class={`px-1 sm:text-base text-sm font-medium py-4 border-b-2 transition-colors ${
                    activeTab() === 'projects'
                      ? 'text-green-500 border-green-500'
                      : 'text-gray-400 hover:text-gray-300 border-transparent'
                  }`}
                >
                  Projects
                </button>
                <button
                  onClick={() => setActiveTab('tags')}
                  class={`px-1 sm:text-base text-sm font-medium py-4 border-b-2 transition-colors ${
                    activeTab() === 'tags'
                      ? 'text-green-500 border-green-500'
                      : 'text-gray-400 hover:text-gray-300 border-transparent'
                  }`}
                >
                  Related Tags
                </button>
              </nav>
            </div>
          </div>

          {/* Tab Content */}
          <div class="mb-8">
            {/* Posts Tab */}
            <Show when={activeTab() === 'posts'}>
              <Show
                when={data().posts && data().posts.length > 0}
                fallback={
                  <div class="bg-gray-800 rounded-lg p-8 text-center">
                    <p class="text-gray-400">No blog posts with this tag yet</p>
                  </div>
                }
              >
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4" data-testid="tag-posts">
                  <For each={data().posts}>
                    {(post) => (
                      <a
                        href={`/posts/${post.slug?.current}`}
                        class="bg-gray-800 rounded-lg p-6 hover:bg-gray-700 transition-colors shadow-lg"
                      >
                        <h3 class="text-xl font-bold text-white mb-2">
                          {post.title || 'Untitled Post'}
                        </h3>
                        <Show when={post.description}>
                          <p class="text-gray-300 mb-3 line-clamp-2">
                            {post.description}
                          </p>
                        </Show>
                        <time datetime={post.createdAt} class="text-sm text-gray-400">
                          {formatDate(post.createdAt)}
                        </time>
                      </a>
                    )}
                  </For>
                </div>
              </Show>
            </Show>

            {/* Projects Tab */}
            <Show when={activeTab() === 'projects'}>
              <Show
                when={data().projects && data().projects.length > 0}
                fallback={
                  <div class="bg-gray-800 rounded-lg p-8 text-center">
                    <p class="text-gray-400">No projects with this tag yet</p>
                  </div>
                }
              >
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4" data-testid="tag-projects">
                  <For each={data().projects}>
                    {(project) => (
                      <a
                        href={`/projects/${project.slug?.current}`}
                        class="bg-gray-800 rounded-lg p-6 hover:bg-gray-700 transition-colors shadow-lg"
                      >
                        <h3 class="text-xl font-bold text-white mb-2">
                          {project.title || 'Untitled Project'}
                        </h3>
                        <Show when={project.description}>
                          <p class="text-gray-300 mb-3 line-clamp-2">
                            {project.description}
                          </p>
                        </Show>
                        <time datetime={project.createdAt} class="text-sm text-gray-400">
                          {formatDate(project.createdAt)}
                        </time>
                      </a>
                    )}
                  </For>
                </div>
              </Show>
            </Show>

            {/* Related Tags Tab */}
            <Show when={activeTab() === 'tags'}>
              <div class="bg-gray-800 rounded-lg p-8 text-center">
                <p class="text-gray-400">Related tags coming soon</p>
              </div>
            </Show>
          </div>

          {/* Back to Tags */}
          <div class="text-center">
            <a
              href="/tags"
              class="inline-flex items-center gap-2 px-6 py-3 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
              data-testid="back-to-tags"
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
              Back to All Tags
            </a>
          </div>
        </div>
      </div>
    </Show>
  );
}
