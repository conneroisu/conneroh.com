import { createFileRoute } from '@tanstack/solid-router';
import { For, Show } from 'solid-js';
import { fetchSanityContent } from '@/lib/sanity-server';
import { listTags } from '@/lib/sanity-queries';
import type { Tag } from '@/lib/sanity-types';

export const Route = createFileRoute('/tags/')({
  component: TagsPage,
  loader: async () => {
    const tags = await fetchSanityContent<Tag[]>(listTags());
    return { tags };
  },
});

function TagsPage() {
  const data = Route.useLoaderData();

  return (
    <div class="min-h-screen bg-gray-900 py-12 px-4 sm:px-6 lg:px-8" data-testid="tags-page">
      <div class="max-w-6xl mx-auto">
        <header class="mb-12">
          <h1 class="text-5xl font-bold text-white mb-4">Skills & Technologies</h1>
          <p class="text-xl text-gray-300">
            Explore my technical skills and tools
          </p>
        </header>

        <Show
          when={data().tags && data().tags.length > 0}
          fallback={
            <div class="text-center py-12 bg-gray-800 rounded-lg" data-testid="tags-empty">
              <p class="text-gray-400 text-lg">No tags found</p>
              <p class="text-gray-500 mt-2">Tags will appear as content is added</p>
            </div>
          }
        >
          <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8" data-testid="tags-list">
            <For each={data().tags}>
                    {(tag) => (
                      <a
                        href={`/tags/${tag.slug?.current}`}
                        class="bg-gray-800 rounded-lg border border-gray-700 hover:border-green-500 transition-all duration-300 p-6 text-center group"
                        data-testid={`tag-item-${tag.slug?.current}`}
                      >
                        <Show when={tag.icon}>
                          <div class="text-4xl mb-2">{tag.icon}</div>
                        </Show>
                        <h2 class="text-lg font-bold text-white group-hover:text-green-400 transition-colors mb-2">
                          {tag.title || 'Untitled Tag'}
                        </h2>
                        <Show when={tag.description}>
                          <p class="text-sm text-gray-400 line-clamp-2">
                            {tag.description}
                          </p>
                        </Show>
            </a>
          )}
        </For>
      </div>
    </Show>
      </div>
    </div>
  );
}
