import "htmx.org";
import "htmx-ext-preload";
import Alpine from "alpinejs";
import intersect from "@alpinejs/intersect";

declare global {
  interface Window {
    Alpine: typeof Alpine;
  }
}

window.Alpine = Alpine;

Alpine.plugin(intersect);
Alpine.start();
