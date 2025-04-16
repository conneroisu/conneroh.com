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

import htmx from "htmx.org";

htmx.logAll();
