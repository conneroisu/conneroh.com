import { Link } from "@tanstack/solid-router";
import { Show } from "solid-js";

interface TagCardProps {
  title: string;
  description?: string;
  icon?: string;
  href: string;
  postCount?: number;
  projectCount?: number;
}

export default function TagCard(props: TagCardProps) {
  return (
    <Link
      to={props.href}
      class="block bg-gray-800 border border-gray-700 rounded-lg overflow-hidden hover:shadow-lg hover:-translate-y-1 transition-all duration-300"
    >
      <div class="flex flex-col h-full relative">
        <div class="p-6 flex-1 flex flex-col">
          <Show when={props.icon}>
            <div class="mb-4">
              <img
                src={props.icon}
                alt={props.title}
                class="w-12 h-12 object-contain"
              />
            </div>
          </Show>

          <h3 class="text-white font-semibold text-lg hover:underline mb-2">
            {props.title}
          </h3>

          <Show when={props.description}>
            <p class="text-gray-300 line-clamp-2">{props.description}</p>
          </Show>
        </div>

        <Show when={props.postCount !== undefined || props.projectCount !== undefined}>
          <div class="relative h-0">
            <Show when={props.postCount !== undefined}>
              <div class="text-xs py-1 bottom-0 px-2 rounded-tr-md text-white left-0 absolute border-emerald-800 border-2 border-l-0 border-b-0 bg-gray-800">
                {props.postCount} posts
              </div>
            </Show>
            <Show when={props.projectCount !== undefined}>
              <div class="px-2 bottom-0 py-1 right-0 text-xs rounded-tl-md text-white absolute border-emerald-800 border-2 border-r-0 border-b-0 bg-gray-800">
                {props.projectCount} projects
              </div>
            </Show>
          </div>
        </Show>
      </div>
    </Link>
  );
}
