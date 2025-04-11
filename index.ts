import "htmx.org";
import "htmx-ext-preload";
import Alpine from "alpinejs";
import intersect from "@alpinejs/intersect";

declare global {
  interface Window {
    Alpine: typeof Alpine;
    MathJax: typeof MathJax;
  }
}

window.Alpine = Alpine;

Alpine.plugin(intersect);
Alpine.start();
//
// import htmx from "htmx.org";
//
// htmx.logAll();
