import ajax from "@imacrayon/alpine-ajax";
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
Alpine.plugin(ajax);
Alpine.start();

// Scroll to top when Alpine AJAX replaces main content
document.addEventListener('ajax:after', function(event) {
  // Check if any of the rendered targets is the main content area
  const targets = (event as any).detail.render || [];
  const hasMainContent = targets.some((target: any) => target?.id === 'bodiody');
  if (hasMainContent) {
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
