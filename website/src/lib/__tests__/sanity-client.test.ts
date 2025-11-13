import { describe, it, expect, beforeEach } from 'vitest';
import { createSanityClient } from '../sanity-client';
import { env } from '../../env';

describe('Sanity Client', () => {
  beforeEach(() => {
    // Tests use explicit configuration to avoid environment variable dependencies
  });

  it('should create a client with custom configuration', () => {
    const client = createSanityClient({
      projectId: 'test-project',
      dataset: 'production',
    });
    expect(client).toBeDefined();
    expect(client.config).toBeDefined();
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
});
