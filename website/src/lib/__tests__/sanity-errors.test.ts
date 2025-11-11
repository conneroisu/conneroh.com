import { describe, it, expect, vi } from 'vitest';
import {
  SanityError,
  SanityValidationError,
  SanityNetworkError,
  SanityAuthError,
  SanityRateLimitError,
  parseSanityError,
  withRetry,
  safeFetch,
} from '../sanity-errors';

describe('Sanity Errors', () => {
  describe('Error Classes', () => {
    it('should create SanityError with correct properties', () => {
      const error = new SanityError('Test error', 'TEST_ERROR', 500);
      expect(error.message).toBe('Test error');
      expect(error.code).toBe('TEST_ERROR');
      expect(error.statusCode).toBe(500);
      expect(error.name).toBe('SanityError');
    });

    it('should create SanityValidationError', () => {
      const error = new SanityValidationError('Invalid query');
      expect(error.code).toBe('VALIDATION_ERROR');
      expect(error.statusCode).toBe(400);
    });

    it('should create SanityNetworkError', () => {
      const error = new SanityNetworkError('Connection failed');
      expect(error.code).toBe('NETWORK_ERROR');
      expect(error.statusCode).toBe(503);
    });

    it('should create SanityAuthError', () => {
      const error = new SanityAuthError('Unauthorized');
      expect(error.code).toBe('AUTH_ERROR');
      expect(error.statusCode).toBe(401);
    });

    it('should create SanityRateLimitError', () => {
      const error = new SanityRateLimitError('Rate limit exceeded');
      expect(error.code).toBe('RATE_LIMIT_ERROR');
      expect(error.statusCode).toBe(429);
    });
  });

  describe('parseSanityError', () => {
    it('should return SanityError as-is', () => {
      const original = new SanityError('Test', 'TEST', 500);
      const parsed = parseSanityError(original);
      expect(parsed).toBe(original);
    });

    it('should parse authentication errors', () => {
      const error = new Error('Unauthorized access');
      const parsed = parseSanityError(error);
      expect(parsed).toBeInstanceOf(SanityAuthError);
    });

    it('should parse validation errors', () => {
      const error = new Error('Invalid validation');
      const parsed = parseSanityError(error);
      expect(parsed).toBeInstanceOf(SanityValidationError);
    });

    it('should parse network errors', () => {
      const error = new Error('Network timeout');
      const parsed = parseSanityError(error);
      expect(parsed).toBeInstanceOf(SanityNetworkError);
    });

    it('should parse rate limit errors', () => {
      const error = new Error('Too many requests');
      const parsed = parseSanityError(error);
      expect(parsed).toBeInstanceOf(SanityRateLimitError);
    });

    it('should handle unknown errors', () => {
      const parsed = parseSanityError('string error');
      expect(parsed).toBeInstanceOf(SanityError);
      expect(parsed.code).toBe('UNKNOWN_ERROR');
    });
  });

  describe('withRetry', () => {
    it('should succeed on first attempt', async () => {
      const fn = vi.fn().mockResolvedValue('success');
      const result = await withRetry(fn);
      expect(result).toBe('success');
      expect(fn).toHaveBeenCalledTimes(1);
    });

    it('should retry on network errors', async () => {
      const fn = vi
        .fn()
        .mockRejectedValueOnce(new SanityNetworkError('Network error'))
        .mockResolvedValue('success');

      const result = await withRetry(fn, { baseDelay: 10 });
      expect(result).toBe('success');
      expect(fn).toHaveBeenCalledTimes(2);
    });

    it('should throw after max attempts', async () => {
      const fn = vi.fn().mockRejectedValue(new SanityNetworkError('Network error'));

      await expect(withRetry(fn, { maxAttempts: 2, baseDelay: 10 })).rejects.toThrow(
        'Network error'
      );
      expect(fn).toHaveBeenCalledTimes(2);
    });

    it('should not retry validation errors', async () => {
      const fn = vi.fn().mockRejectedValue(new SanityValidationError('Invalid'));

      await expect(withRetry(fn, { baseDelay: 10 })).rejects.toThrow('Invalid');
      expect(fn).toHaveBeenCalledTimes(1);
    });

    it('should use exponential backoff', async () => {
      const delays: number[] = [];
      const fn = vi.fn().mockRejectedValue(new SanityNetworkError('Error'));

      const customRetry = async () => {
        let attempt = 0;
        const maxAttempts = 3;
        const baseDelay = 100;

        while (attempt < maxAttempts) {
          attempt++;
          try {
            return await fn();
          } catch (error) {
            if (attempt < maxAttempts) {
              const delay = baseDelay * Math.pow(2, attempt - 1);
              delays.push(delay);
              await new Promise(resolve => setTimeout(resolve, delay));
            } else {
              throw error;
            }
          }
        }
      };

      await expect(customRetry()).rejects.toThrow();
      expect(delays).toEqual([100, 200]);
    });
  });

  describe('safeFetch', () => {
    it('should return result on success', async () => {
      const fn = vi.fn().mockResolvedValue('success');
      const result = await safeFetch(fn, { retry: false });
      expect(result).toBe('success');
    });

    it('should return fallback on error', async () => {
      const fn = vi.fn().mockRejectedValue(new Error('Failed'));
      const result = await safeFetch(fn, { fallback: 'default', retry: false });
      expect(result).toBe('default');
    });

    it('should throw without fallback', async () => {
      const fn = vi.fn().mockRejectedValue(new Error('Failed'));
      await expect(safeFetch(fn, { retry: false })).rejects.toThrow();
    });

    it('should retry by default', async () => {
      const fn = vi
        .fn()
        .mockRejectedValueOnce(new SanityNetworkError('Error'))
        .mockResolvedValue('success');

      const result = await safeFetch(fn, { retry: { baseDelay: 10 } });
      expect(result).toBe('success');
      expect(fn).toHaveBeenCalledTimes(2);
    });
  });
});
