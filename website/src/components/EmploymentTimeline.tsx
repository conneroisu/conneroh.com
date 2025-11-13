import { For, Show } from "solid-js";
import { Link } from "@tanstack/solid-router";

interface Employment {
  title: string;
  description?: string;
  startDate: string;
  endDate?: string;
  tags?: Array<{ _ref?: string; title?: string; slug?: { current?: string } }>;
  href: string;
}

interface EmploymentTimelineProps {
  items: Employment[];
}

// Date formatting helper
const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", { month: "short", year: "numeric" });
};

// Format date range helper
const formatDateRange = (startDate: string, endDate?: string) => {
  const start = formatDate(startDate);
  const end = endDate ? formatDate(endDate) : "Present";
  return `${start} - ${end}`;
};

// Desktop Timeline Component
const EmploymentTimelineDesktop = (props: EmploymentTimelineProps) => {
  return (
    <div class="relative hidden md:block">
      {/* Central vertical timeline line */}
      <div class="absolute left-1/2 -translate-x-1/2 w-0.5 h-full bg-gradient-to-b from-green-500 via-green-400 to-transparent" />

      <div class="space-y-12">
        <For each={props.items}>
          {(item, index) => {
            const isEven = index() % 2 === 0;
            const maxTags = 4;
            const displayTags = item.tags?.slice(0, maxTags) || [];
            const extraTagsCount = (item.tags?.length || 0) - maxTags;

            return (
              <div
                class={`relative flex ${
                  isEven ? "justify-start" : "justify-end"
                }`}
              >
                {/* Timeline circle */}
                <div class="absolute left-1/2 -translate-x-1/2 w-4 h-4 bg-green-500 rounded-full border-4 border-gray-900 z-10" />

                {/* Card */}
                <Link
                  to={item.href}
                  class={`w-5/12 bg-gray-800 rounded-lg p-6 hover:shadow-xl transition-all hover:-translate-y-1 hover:bg-gray-750 ${
                    isEven ? "text-right" : "text-left"
                  }`}
                >
                  {/* Date range */}
                  <div class="text-green-400 font-semibold text-sm mb-2">
                    {formatDateRange(item.startDate, item.endDate)}
                  </div>

                  {/* Title */}
                  <h3 class="text-white font-bold text-xl mb-3">
                    {item.title}
                  </h3>

                  {/* Description */}
                  <Show when={item.description}>
                    <p class="text-gray-300 text-sm mb-4 line-clamp-3">
                      {item.description}
                    </p>
                  </Show>

                  {/* Tags */}
                  <Show when={displayTags.length > 0}>
                    <div
                      class={`flex gap-2 flex-wrap ${
                        isEven ? "justify-end" : "justify-start"
                      }`}
                    >
                      <For each={displayTags}>
                        {(tag) => (
                          <Show when={tag.slug?.current}>
                            <Link
                              to={`/tags/${tag.slug!.current}`}
                              class="text-xs bg-gray-700 text-gray-300 px-2 py-1 rounded hover:bg-green-600 hover:text-white transition-colors"
                              onClick={(e) => e.stopPropagation()}
                            >
                              {tag.title || tag.slug!.current}
                            </Link>
                          </Show>
                        )}
                      </For>
                      <Show when={extraTagsCount > 0}>
                        <span class="text-xs text-gray-500 self-center">
                          +{extraTagsCount}
                        </span>
                      </Show>
                    </div>
                  </Show>
                </Link>
              </div>
            );
          }}
        </For>
      </div>
    </div>
  );
};

// Mobile Timeline Component
const EmploymentTimelineMobile = (props: EmploymentTimelineProps) => {
  return (
    <div class="relative md:hidden">
      {/* Left vertical timeline line */}
      <div class="absolute left-4 w-0.5 h-full bg-gradient-to-b from-green-500 via-green-400 to-transparent" />

      <div class="space-y-8 pl-12">
        <For each={props.items}>
          {(item) => {
            const maxTags = 4;
            const displayTags = item.tags?.slice(0, maxTags) || [];
            const extraTagsCount = (item.tags?.length || 0) - maxTags;

            return (
              <div class="relative">
                {/* Timeline circle */}
                <div class="absolute -left-8 top-6 w-4 h-4 bg-green-500 rounded-full border-4 border-gray-900" />

                {/* Card */}
                <Link
                  to={item.href}
                  class="block bg-gray-800 rounded-lg p-4 hover:shadow-xl transition-all hover:-translate-y-1 hover:bg-gray-750"
                >
                  {/* Date range */}
                  <div class="text-green-400 font-semibold text-sm mb-2">
                    {formatDateRange(item.startDate, item.endDate)}
                  </div>

                  {/* Title */}
                  <h3 class="text-white font-bold text-lg mb-3">
                    {item.title}
                  </h3>

                  {/* Description */}
                  <Show when={item.description}>
                    <p class="text-gray-300 text-sm mb-4 line-clamp-2">
                      {item.description}
                    </p>
                  </Show>

                  {/* Tags */}
                  <Show when={displayTags.length > 0}>
                    <div class="flex gap-2 flex-wrap">
                      <For each={displayTags}>
                        {(tag) => (
                          <Show when={tag.slug?.current}>
                            <Link
                              to={`/tags/${tag.slug!.current}`}
                              class="text-xs bg-gray-700 text-gray-300 px-2 py-1 rounded hover:bg-green-600 hover:text-white transition-colors"
                              onClick={(e) => e.stopPropagation()}
                            >
                              {tag.title || tag.slug!.current}
                            </Link>
                          </Show>
                        )}
                      </For>
                      <Show when={extraTagsCount > 0}>
                        <span class="text-xs text-gray-500 self-center">
                          +{extraTagsCount}
                        </span>
                      </Show>
                    </div>
                  </Show>
                </Link>
              </div>
            );
          }}
        </For>
      </div>
    </div>
  );
};

// Main Export Component
export const EmploymentTimeline = (props: EmploymentTimelineProps) => {
  return (
    <div class="w-full">
      <EmploymentTimelineDesktop items={props.items} />
      <EmploymentTimelineMobile items={props.items} />
    </div>
  );
};

export default EmploymentTimeline;
