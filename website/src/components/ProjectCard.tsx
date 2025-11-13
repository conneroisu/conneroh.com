import { Link } from "@tanstack/solid-router";
import { Show, For } from "solid-js";

interface ProjectCardProps {
  title: string;
  description?: string;
  bannerPath?: string;
  tags?: Array<{ _ref?: string; slug?: { current?: string } }>;
  href: string;
  postCount?: number;
  tagCount?: number;
}

export default function ProjectCard(props: ProjectCardProps) {
  return (
    <Link
      to={props.href}
      class="block bg-gray-800 border border-gray-700 rounded-lg overflow-hidden hover:shadow-lg hover:-translate-y-1 transition-all duration-300"
    >
      <div class="flex flex-col h-full relative">
        <Show when={props.bannerPath}>
          <div class="w-full h-48">
            <img
              src={props.bannerPath}
              alt={props.title}
              class="w-full h-full object-cover"
            />
          </div>
        </Show>

        <div class="p-6 flex-1 flex flex-col">
          <h3 class="text-white font-bold text-xl hover:underline mb-2">
            {props.title}
          </h3>

          <Show when={props.description}>
            <p class="text-gray-300 line-clamp-2 mb-4">{props.description}</p>
          </Show>

          <Show when={props.tags && props.tags.length > 0}>
            <div class="flex flex-wrap gap-2 mt-auto">
              <For each={props.tags.slice(0, 3)}>
                {(tag) => (
                  <Show when={tag.slug?.current}>
                    <span class="px-3 py-1 bg-emerald-800 text-emerald-300 rounded-full text-sm hover:bg-emerald-600 transition-colors">
                      {tag.slug?.current}
                    </span>
                  </Show>
                )}
              </For>
              <Show when={props.tags.length > 3}>
                <span class="px-3 py-1 bg-emerald-800 text-emerald-300 rounded-full text-sm">
                  +{props.tags.length - 3}
                </span>
              </Show>
            </div>
          </Show>
        </div>

        <Show when={props.postCount !== undefined || props.tagCount !== undefined}>
          <div class="relative h-0">
            <Show when={props.postCount !== undefined}>
              <div class="text-xs py-1 bottom-0 px-2 rounded-tr-md text-white left-0 absolute border-emerald-800 border-2 border-l-0 border-b-0 bg-gray-800">
                {props.postCount} posts
              </div>
            </Show>
            <Show when={props.tagCount !== undefined}>
              <div class="px-2 bottom-0 py-1 right-0 text-xs rounded-tl-md text-white absolute border-emerald-800 border-2 border-r-0 border-b-0 bg-gray-800">
                {props.tagCount} tags
              </div>
            </Show>
          </div>
        </Show>
      </div>
    </Link>
  );
}
