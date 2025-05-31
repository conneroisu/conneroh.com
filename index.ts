import "htmx.org";
import "htmx-ext-preload";
import Alpine from "alpinejs";
import intersect from "@alpinejs/intersect";
import anchor from "@alpinejs/anchor";

declare global {
  interface Window {
    Alpine: typeof Alpine;
    MathJax: typeof MathJax;
  }
}

window.Alpine = Alpine;

Alpine.plugin(intersect);
Alpine.plugin(anchor);
Alpine.start();

htmx.config.globalViewTransitions = true;

// Scroll to top when HTMX replaces main content
document.addEventListener('htmx:afterSwap', function(event) {
  // Check if the swapped element is the main content area
  if ((event as any).detail.target.id === 'bodiody') {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
});

// Enhance anchor link scrolling for smooth behavior
document.addEventListener('DOMContentLoaded', function() {
  // Handle all anchor links that start with #
  document.addEventListener('click', function(event) {
    const target = event.target as HTMLElement;
    if (target.tagName === 'A') {
      const link = target as HTMLAnchorElement;
      const href = link.getAttribute('href');

      // If it's a hash link to an element on the same page
      if (href && href.startsWith('#') && href.length > 1) {
        const targetElement = document.querySelector(href);
        if (targetElement) {
          event.preventDefault();
          targetElement.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
          });

          // Update URL without jumping
          if (history.pushState) {
            history.pushState(null, '', href);
          }
        }
      }
    }
  });
});
