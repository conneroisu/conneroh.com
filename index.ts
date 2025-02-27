import "htmx.org";
import htmx from "htmx.org";
import Alpine from "alpinejs";

declare global {
  interface Window {
    Alpine: typeof Alpine;
  }
}

type Post = {
  id: number;
  title: string;
  description: string;
  slug: string;
  date: number;
  banner: string;
  tags: Tag[];
};

type Tag = {
  id: number;
  name: string;
  slug: string;
};

window.Alpine = Alpine;

Alpine.start();

// This code will implement posts filtering by +tag and -tag
// Function to parse URL query params and extract tag filters
function parseTagFilters(queryString: string) {
  const params = new URLSearchParams(queryString);
  const tags = params.get("tags") || "";

  const includeTags: string[] = [];
  const excludeTags: string[] = [];

  // Parse the tags parameter for +tag and -tag patterns
  const tagPatterns = tags.split(/\s+/).filter((tag) => tag.trim().length > 0);

  tagPatterns.forEach((pattern) => {
    if (pattern.startsWith("+")) {
      const tag = pattern.substring(1).trim();
      if (tag) includeTags.push(tag);
    } else if (pattern.startsWith("-")) {
      const tag = pattern.substring(1).trim();
      if (tag) excludeTags.push(tag);
    } else {
      // If no prefix, assume include
      if (pattern.trim()) includeTags.push(pattern.trim());
    }
  });

  return { includeTags, excludeTags };
}

// Function to filter posts based on tags
function filterPostsByTags(
  posts: Post[],
  includeTags: string[],
  excludeTags: string[],
) {
  if (includeTags.length === 0 && excludeTags.length === 0) {
    return posts; // No filtering needed
  }

  return posts.filter((post) => {
    const postTags = post.tags.map((tag) => tag.name.toLowerCase());

    // Check if post should be excluded
    if (excludeTags.length > 0) {
      const shouldExclude = excludeTags.some((tag) =>
        postTags.includes(tag.toLowerCase()),
      );
      if (shouldExclude) return false;
    }

    // Check if post should be included
    if (includeTags.length > 0) {
      const shouldInclude = includeTags.every((tag) =>
        postTags.includes(tag.toLowerCase()),
      );
      return shouldInclude;
    }

    // If we have excludes but no includes, keep all posts that weren't excluded
    return true;
  });
}

// Helper to reconstruct current tags parameter
function currentTagsParam(includeTags: string[], excludeTags: string[]) {
  let tagsParam = "";

  includeTags.forEach((tag) => {
    tagsParam += " +" + tag;
  });

  excludeTags.forEach((tag) => {
    tagsParam += " -" + tag;
  });

  return tagsParam.trim();
}

// Main code to integrate with the posts view
document.addEventListener("alpine:init", () => {
  Alpine.data("postFilters", () => ({
    queryString: window.location.search,
    allTags: Array<string>(),
    includeTags: Array<string>(),
    excludeTags: Array<string>(),

    init() {
      // Extract all unique tags from posts
      this.allTags = [
        ...new Set(
          Array.from(document.querySelectorAll("[data-post-tags]"))
            .map((el) => JSON.parse(el.getAttribute("data-post-tags") || "[]"))
            .flat()
            .map((tag) => tag.name),
        ),
      ];

      // Parse current filters
      const { includeTags, excludeTags } = parseTagFilters(this.queryString);
      this.includeTags = includeTags;
      this.excludeTags = excludeTags;

      // Apply filters to Alpine.js data
      this.$watch("filteredPosts", (posts) => {
        const filtered = filterPostsByTags(
          posts,
          this.includeTags,
          this.excludeTags,
        );
        return filtered;
      });
    },

    toggleIncludeTag(tag: string) {
      if (this.includeTags.includes(tag)) {
        this.includeTags = this.includeTags.filter((t) => t !== tag);
      } else {
        this.includeTags.push(tag);
        // Remove from exclude if it's there
        this.excludeTags = this.excludeTags.filter((t) => t !== tag);
      }
      this.updateURL();
    },

    toggleExcludeTag(tag: string) {
      if (this.excludeTags.includes(tag)) {
        this.excludeTags = this.excludeTags.filter((t) => t !== tag);
      } else {
        this.excludeTags.push(tag);
        // Remove from include if it's there
        this.includeTags = this.includeTags.filter((t) => t !== tag);
      }
      this.updateURL();
    },

    updateURL() {
      let tagsParam = currentTagsParam(this.includeTags, this.excludeTags);

      // Update the URL
      const url = new URL(window.URL.toString());
      if (tagsParam) {
        url.searchParams.set("tags", tagsParam);
      } else {
        url.searchParams.delete("tags");
      }

      window.history.pushState({}, "", url);
    },
  }));
});
