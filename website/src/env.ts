import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

/**
 * Consolidated Sanity CMS environment variable configuration.
 *
 * This module validates and exports Sanity environment variables using @t3-oss/env-core
 * with Zod schema validation. It supports both client-side (Vite) and server-side (Node)
 * runtime environments.
 */

export const env = createEnv({
  /**
   * Client-side environment variables.
   * These are prefixed with VITE_ and are exposed to the browser.
   */
  client: {
    /**
     * Sanity Project ID - Identifies the Sanity project instance.
     * Safe to expose to clients as it's used for public API endpoints.
     */
    VITE_SANITY_PROJECT_ID: z.string().min(1),

    /**
     * Sanity Dataset - The dataset name within the Sanity project (e.g., "production", "development").
     * Safe to expose to clients as it's used for public content queries.
     */
    VITE_SANITY_DATASET: z.string().min(1),
  },

  /**
   * Server-side environment variables.
   * These are NEVER exposed to the browser and remain secure on the server.
   */
  server: {
    /**
     * Sanity API Token - Authentication token for write operations and admin queries.
     * CRITICAL: Must remain server-side only. Never expose to client code.
     */
    SANITY_API_TOKEN: z.string().min(1),
  },

  /**
   * Configure the client-side environment variable prefix.
   * Vite only exposes variables prefixed with VITE_ to the browser.
   */
  clientPrefix: "VITE_",

  /**
   * Runtime environment variable sources.
   * Supports both import.meta.env (Vite/browser) and process.env (Node/build time).
   */
  runtimeEnv: {
    // Client variables (safe for browser)
    VITE_SANITY_PROJECT_ID:
      typeof import.meta !== "undefined"
        ? import.meta.env?.VITE_SANITY_PROJECT_ID
        : process.env.VITE_SANITY_PROJECT_ID,
    VITE_SANITY_DATASET:
      typeof import.meta !== "undefined"
        ? import.meta.env?.VITE_SANITY_DATASET
        : process.env.VITE_SANITY_DATASET,

    // Server-only variable (NEVER exposed to browser)
    SANITY_API_TOKEN: process.env.SANITY_API_TOKEN,
  },

  /**
   * Skip validation during build if environment variables are not yet available.
   * Set to true in CI/CD environments where variables are injected at runtime.
   */
  skipValidation: !!process.env.SKIP_ENV_VALIDATION,
});

// Export as both default and named export for flexibility
export default env;
