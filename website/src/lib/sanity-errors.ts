/**
 * Base class for all Sanity-related errors
 */
export class SanityError extends Error {
  constructor(
    message: string,
    public readonly code: string,
    public readonly statusCode?: number,
    public readonly details?: unknown
  ) {
    super(message);
    this.name = 'SanityError';
    Object.setPrototypeOf(this, SanityError.prototype);
  }
}

/**
 * Validation error for invalid queries or data
 */
export class SanityValidationError extends SanityError {
  constructor(message: string, details?: unknown) {
    super(message, 'VALIDATION_ERROR', 400, details);
    this.name = 'SanityValidationError';
    Object.setPrototypeOf(this, SanityValidationError.prototype);
  }
}

/**
 * Network error for connection issues
 */
export class SanityNetworkError extends SanityError {
  constructor(message: string, details?: unknown) {
    super(message, 'NETWORK_ERROR', 503, details);
    this.name = 'SanityNetworkError';
    Object.setPrototypeOf(this, SanityNetworkError.prototype);
  }
}

/**
 * Authentication error for invalid or missing credentials
 */
export class SanityAuthError extends SanityError {
  constructor(message: string, details?: unknown) {
    super(message, 'AUTH_ERROR', 401, details);
    this.name = 'SanityAuthError';
    Object.setPrototypeOf(this, SanityAuthError.prototype);
  }
}

/**
 * Rate limit error when API limits are exceeded
 */
export class SanityRateLimitError extends SanityError {
  constructor(message: string, details?: unknown) {
    super(message, 'RATE_LIMIT_ERROR', 429, details);
    this.name = 'SanityRateLimitError';
    Object.setPrototypeOf(this, SanityRateLimitError.prototype);
  }
}

/**
 * Options for retry logic
 */
export interface RetryOptions {
  /**
   * Maximum number of retry attempts (default: 3)
   */
  maxAttempts?: number;
  /**
   * Base delay in milliseconds (default: 1000)
   */
  baseDelay?: number;
  /**
   * Maximum delay in milliseconds (default: 10000)
   */
  maxDelay?: number;
  /**
   * Exponential backoff multiplier (default: 2)
   */
  backoffMultiplier?: number;
  /**
   * Should retry this error? (default: retries network errors only)
   */
  shouldRetry?: (error: Error, attempt: number) => boolean;
}

/**
 * Default retry logic: only retry network errors
 */
function defaultShouldRetry(error: Error, attempt: number): boolean {
  if (attempt >= 3) return false;
  // Retry on network errors, rate limits, or generic errors with "network" in the message
  return (
    error instanceof SanityNetworkError ||
    error instanceof SanityRateLimitError ||
    error.message.toLowerCase().includes('network')
  );
}

/**
 * Execute a function with automatic retry on failure
 *
 * @example
 * ```typescript
 * const result = await withRetry(
 *   async () => client.fetch('*[_type == "post"]'),
 *   { maxAttempts: 3, baseDelay: 1000 }
 * );
 * ```
 *
 * @param fn - Async function to execute
 * @param options - Retry configuration
 * @returns Result of the function
 * @throws Last error if all retries fail
 */
export async function withRetry<T>(
  fn: () => Promise<T>,
  options: RetryOptions = {}
): Promise<T> {
  const {
    maxAttempts = 3,
    baseDelay = 1000,
    maxDelay = 10000,
    backoffMultiplier = 2,
    shouldRetry = defaultShouldRetry,
  } = options;

  let lastError: Error | undefined;

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error instanceof Error ? error : new Error(String(error));

      if (attempt < maxAttempts && shouldRetry(lastError, attempt)) {
        const delay = Math.min(baseDelay * Math.pow(backoffMultiplier, attempt - 1), maxDelay);
        logRetry(attempt, maxAttempts, delay, lastError);
        await sleep(delay);
      } else {
        break;
      }
    }
  }

  throw lastError;
}

/**
 * Sleep utility for retry delays
 */
function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

/**
 * Parse and normalize Sanity API errors
 *
 * @param error - Raw error from Sanity client
 * @returns Normalized SanityError instance
 */
export function parseSanityError(error: unknown): SanityError {
  if (error instanceof SanityError) {
    return error;
  }

  if (error instanceof Error) {
    const message = error.message;

    // Check for specific error patterns
    if (message.includes('Unauthorized') || message.includes('Authentication')) {
      return new SanityAuthError(message, error);
    }

    if (message.includes('validation') || message.includes('invalid')) {
      return new SanityValidationError(message, error);
    }

    if (
      message.toLowerCase().includes('rate limit') ||
      message.toLowerCase().includes('too many requests')
    ) {
      return new SanityRateLimitError(message, error);
    }

    if (
      message.includes('network') ||
      message.includes('timeout') ||
      message.includes('ECONNREFUSED')
    ) {
      return new SanityNetworkError(message, error);
    }

    // Generic Sanity error
    return new SanityError(message, 'UNKNOWN_ERROR', undefined, error);
  }

  // Unknown error type
  return new SanityError('An unknown error occurred', 'UNKNOWN_ERROR', undefined, error);
}

/**
 * Log levels for debugging
 */
export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

/**
 * Simple logger for Sanity operations
 */
export const logger = {
  debug: (message: string, ...args: unknown[]) => {
    if (import.meta.env.DEV) {
      console.debug(`[Sanity Debug] ${message}`, ...args);
    }
  },

  info: (message: string, ...args: unknown[]) => {
    if (import.meta.env.DEV) {
      console.info(`[Sanity Info] ${message}`, ...args);
    }
  },

  warn: (message: string, ...args: unknown[]) => {
    console.warn(`[Sanity Warning] ${message}`, ...args);
  },

  error: (message: string, error?: unknown, ...args: unknown[]) => {
    console.error(`[Sanity Error] ${message}`, error, ...args);
  },
};

/**
 * Log retry attempt
 */
function logRetry(attempt: number, maxAttempts: number, delay: number, error: Error): void {
  logger.warn(
    `Retry attempt ${attempt}/${maxAttempts} after ${delay}ms`,
    error.message
  );
}

/**
 * Helper to safely execute Sanity operations with error handling
 *
 * @example
 * ```typescript
 * const posts = await safeFetch(
 *   async () => client.fetch('*[_type == "post"]'),
 *   { fallback: [] }
 * );
 * ```
 */
export async function safeFetch<T>(
  fn: () => Promise<T>,
  options: {
    fallback?: T;
    logError?: boolean;
    retry?: boolean | RetryOptions;
  } = {}
): Promise<T> {
  const { fallback, logError = true, retry = true } = options;

  try {
    if (retry) {
      const retryOptions = typeof retry === 'boolean' ? {} : retry;
      return await withRetry(fn, retryOptions);
    }
    return await fn();
  } catch (error) {
    const sanityError = parseSanityError(error);

    if (logError) {
      logger.error('Sanity operation failed', sanityError);
    }

    if (fallback !== undefined) {
      return fallback;
    }

    throw sanityError;
  }
}
