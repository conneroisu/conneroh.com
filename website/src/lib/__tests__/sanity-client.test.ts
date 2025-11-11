import { describe, it, expect, beforeEach, vi } from 'vitest';
import { createSanityClient, getSanityClient, createServerClient } from '../sanity-client';

describe('Sanity Client', () => {
  beforeEach(() => {
    // Clear any mocked environment variables before each test
    vi.unstubAllEnvs();
  });

  it('should create a client with custom configuration', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
    });
    expect(client).toBeDefined();
    expect(client.config).toBeDefined();
  });

  it('should throw error if project ID is missing', () => {
    // Mock empty environment variables to ensure they're not set
    vi.stubEnv('VITE_SANITY_PROJECT_ID', '');
    vi.stubEnv('VITE_SANITY_DATASET', 'production');

    expect(() =>
      createSanityClient({
        dataset: 'production',
      })
    ).toThrow('Sanity project ID is required');
  });

  it('should throw error if dataset is missing', () => {
    // Mock empty environment variables to ensure they're not set
    vi.stubEnv('VITE_SANITY_PROJECT_ID', 'test-project');
    vi.stubEnv('VITE_SANITY_DATASET', '');

    expect(() =>
      createSanityClient({
        projectId: 'test-project',
      })
    ).toThrow('Sanity dataset is required');
  });

  it('should accept full custom configuration', () => {
    const client = createSanityClient({
      projectId: 'custom-project',
      dataset: 'custom-dataset',
      apiVersion: '2025-01-01',
    });
    expect(client).toBeDefined();
    expect(client.config().apiVersion).toBe('2025-01-01');
  });

  it('should use CDN by default for browser clients', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
    });
    expect(client.config().useCdn).toBe(true);
  });

  it('should disable CDN when token is provided', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
      token: 'test-token',
    });
    expect(client.config().useCdn).toBe(false);
  });

  it('should create client with published perspective by default', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
    });
    expect(client.config().perspective).toBe('published');
  });

  it('should accept custom perspective', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
      perspective: 'previewDrafts',
    });
    expect(client.config().perspective).toBe('previewDrafts');
  });

  it('should return singleton client from getSanityClient with env vars', () => {
    // Mock environment variables
    vi.stubEnv('VITE_SANITY_PROJECT_ID', 'test-project');
    vi.stubEnv('VITE_SANITY_DATASET', 'production');

    // Reset the singleton before testing
    // Note: This is a workaround since we can't directly reset the module state
    // In a real scenario, the singleton would persist across tests
    const client1 = getSanityClient();
    const client2 = getSanityClient();
    expect(client1).toBeDefined();
    expect(client2).toBeDefined();
    // Both should have the same configuration
    expect(client1.config().projectId).toBe(client2.config().projectId);
  });

  it('should use environment variables when options are not provided', () => {
    // Mock environment variables
    vi.stubEnv('VITE_SANITY_PROJECT_ID', 'env-project-id');
    vi.stubEnv('VITE_SANITY_DATASET', 'env-dataset');

    const client = createSanityClient();
    expect(client.config().projectId).toBe('env-project-id');
    expect(client.config().dataset).toBe('env-dataset');
  });

  it('should prefer explicit options over environment variables', () => {
    // Mock environment variables
    vi.stubEnv('VITE_SANITY_PROJECT_ID', 'env-project-id');
    vi.stubEnv('VITE_SANITY_DATASET', 'env-dataset');

    const client = createSanityClient({
      projectId: 'explicit-project-id',
      dataset: 'explicit-dataset',
    });
    expect(client.config().projectId).toBe('explicit-project-id');
    expect(client.config().dataset).toBe('explicit-dataset');
  });

  it('should set correct API version', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
      apiVersion: '2024-01-01',
    });
    expect(client.config().apiVersion).toBe('2024-01-01');
  });

  it('should use default API version when not provided', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
    });
    expect(client.config().apiVersion).toBe('2025-01-11');
  });

  it('should configure token for authenticated requests', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
      token: 'sk-test-token-123',
    });
    // Token should be set (we can't directly check it, but CDN should be disabled)
    expect(client.config().useCdn).toBe(false);
  });

  describe('createServerClient', () => {
    it('should throw error if SANITY_API_TOKEN is not set', () => {
      // Ensure process.env.SANITY_API_TOKEN is not set
      const originalToken = process.env.SANITY_API_TOKEN;
      delete process.env.SANITY_API_TOKEN;

      expect(() => createServerClient()).toThrow(
        'SANITY_API_TOKEN is required for server-side operations'
      );

      // Restore original value
      if (originalToken) {
        process.env.SANITY_API_TOKEN = originalToken;
      }
    });

    it('should create authenticated client with server token', () => {
      // Mock process.env for server-side
      const originalToken = process.env.SANITY_API_TOKEN;
      process.env.SANITY_API_TOKEN = 'sk-server-token-123';
      vi.stubEnv('VITE_SANITY_PROJECT_ID', 'test-project');
      vi.stubEnv('VITE_SANITY_DATASET', 'production');

      const client = createServerClient();

      expect(client).toBeDefined();
      expect(client.config().useCdn).toBe(false);
      expect(client.config().perspective).toBe('previewDrafts');

      // Clean up
      if (originalToken) {
        process.env.SANITY_API_TOKEN = originalToken;
      } else {
        delete process.env.SANITY_API_TOKEN;
      }
    });
  });
});
