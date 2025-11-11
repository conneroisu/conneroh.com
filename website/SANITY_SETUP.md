# Sanity CMS Integration Setup

## Phase 1: Installation & Configuration - COMPLETED

This document provides instructions for completing the Sanity CMS setup, specifically the API token creation step.

## Configuration Files

### Environment Variables (.env.local)

The following environment variables have been configured in `/home/connerohnesorge/Documents/001Repos/conneroh.com/website/.env.local`:

```env
# Sanity CMS Configuration
VITE_SANITY_PROJECT_ID=z6sd4zry
VITE_SANITY_DATASET=production

# IMPORTANT: Create an API token in Sanity Project Settings → API → Add Token
# Use "Read" permissions for read-only operations, "Write" for mutations
# Then uncomment and add your token below:
# SANITY_API_TOKEN=your_token_here
```

## REQUIRED: Creating a Sanity API Token

To enable server-side operations (mutations, authenticated queries), you need to create an API token:

### Step-by-Step Instructions:

1. **Navigate to Sanity Project Dashboard**
   - Go to https://sanity.io/manage
   - Select your project: `z6sd4zry`

2. **Access API Settings**
   - Click on "API" in the left sidebar
   - Navigate to the "Tokens" section
   - Click "Add API token"

3. **Configure Token**
   - **Name**: "Website Backend" (or any descriptive name)
   - **Permissions**:
     - For read-only operations: Select "Viewer"
     - For read/write operations: Select "Editor"
   - Click "Add token"

4. **Copy Token**
   - The token will only be shown ONCE
   - Copy it immediately

5. **Add to .env.local**
   - Open `/home/connerohnesorge/Documents/001Repos/conneroh.com/website/.env.local`
   - Uncomment the `SANITY_API_TOKEN` line
   - Replace `your_token_here` with your actual token:
   ```env
   SANITY_API_TOKEN=sk_your_actual_token_here
   ```

6. **Security Notes**
   - NEVER commit `.env.local` to version control (it's already in .gitignore)
   - Store the token securely (password manager, secure notes)
   - Rotate tokens periodically for security
   - Use different tokens for development/staging/production environments

## Installed Dependencies

The following packages have been installed:

- `@sanity/client@^7.12.1` - Sanity JavaScript client for querying content

## Created Files

### 1. Sanity Client Wrapper
**Location**: `/home/connerohnesorge/Documents/001Repos/conneroh.com/website/src/lib/sanity-client.ts`

**Purpose**: Provides factory functions for creating Sanity client instances

**Exports**:
- `createSanityClient(options?)` - Create a custom Sanity client
- `getSanityClient()` - Get singleton browser client (CDN-enabled)
- `createServerClient()` - Create authenticated server-side client
- TypeScript types: `SanityClientOptions`, `SanityQueryResult<T>`, `SanityQueryParams`

**Usage Examples**:

```typescript
// Browser client (public, read-only, uses CDN)
import { getSanityClient } from '@/lib/sanity-client';

const client = getSanityClient();
const posts = await client.fetch('*[_type == "post"]');

// Server client (authenticated, for mutations)
import { createServerClient } from '@/lib/sanity-client';

export async function createPost(title: string) {
  const client = createServerClient();
  return await client.create({
    _type: 'post',
    title,
    publishedAt: new Date().toISOString(),
  });
}

// Custom client with specific configuration
import { createSanityClient } from '@/lib/sanity-client';

const customClient = createSanityClient({
  projectId: 'custom-project',
  dataset: 'development',
  useCdn: false,
  perspective: 'previewDrafts',
});
```

### 2. Error Handling Utilities
**Location**: `/home/connerohnesorge/Documents/001Repos/conneroh.com/website/src/lib/sanity-errors.ts`

**Purpose**: Comprehensive error handling, retry logic, and logging for Sanity operations

**Exports**:
- Error classes: `SanityError`, `SanityValidationError`, `SanityNetworkError`, `SanityAuthError`, `SanityRateLimitError`
- `parseSanityError(error)` - Normalize errors from Sanity API
- `withRetry(fn, options?)` - Execute with exponential backoff retry
- `safeFetch(fn, options?)` - Safe wrapper with fallback support
- `logger` - Logging utilities (debug, info, warn, error)

**Usage Examples**:

```typescript
import { withRetry, safeFetch, parseSanityError } from '@/lib/sanity-errors';

// Automatic retry on network errors
const posts = await withRetry(
  async () => client.fetch('*[_type == "post"]'),
  { maxAttempts: 3, baseDelay: 1000 }
);

// Safe fetch with fallback
const posts = await safeFetch(
  async () => client.fetch('*[_type == "post"]'),
  { fallback: [] }
);

// Custom error handling
try {
  await client.create({ _type: 'post' });
} catch (error) {
  const sanityError = parseSanityError(error);
  console.error(`Failed with code: ${sanityError.code}`, sanityError);
}
```

### 3. Tests
**Location**:
- `/home/connerohnesorge/Documents/001Repos/conneroh.com/website/src/lib/__tests__/sanity-client.test.ts`
- `/home/connerohnesorge/Documents/001Repos/conneroh.com/website/src/lib/__tests__/sanity-errors.test.ts`

**Test Coverage**:
- Client configuration validation
- Environment variable handling
- CDN usage based on token presence
- Perspective configuration
- Error classification and parsing
- Retry logic with exponential backoff
- Safe fetch with fallbacks

**Run Tests**:
```bash
bun test src/lib/__tests__/sanity-client.test.ts
bun test src/lib/__tests__/sanity-errors.test.ts
```

**Test Results**: All 28 tests passing ✓

## Validation Checklist

- [x] T1.1: Install `@sanity/client` dependency
- [x] T1.2: Create environment configuration (.env.local, .env.example)
- [x] T1.3: Create Sanity client wrapper (src/lib/sanity-client.ts)
- [x] T1.4: Add error handling utilities (src/lib/sanity-errors.ts)
- [x] Build verification (production build succeeds)
- [x] Test coverage (28 tests passing)
- [ ] **MANUAL STEP REQUIRED**: Create API token in Sanity dashboard

## Next Steps

1. **Complete API Token Setup** (see instructions above)
2. **Verify Token Works**:
   ```typescript
   import { createServerClient } from '@/lib/sanity-client';

   try {
     const client = createServerClient();
     console.log('✓ Server client created successfully');
   } catch (error) {
     console.error('✗ Token not configured:', error.message);
   }
   ```

3. **Proceed to Phase 2**: Query Helpers & Data Fetching

## Troubleshooting

### Error: "Sanity project ID is required"
**Solution**: Ensure `VITE_SANITY_PROJECT_ID=z6sd4zry` is in `.env.local`

### Error: "Sanity dataset is required"
**Solution**: Ensure `VITE_SANITY_DATASET=production` is in `.env.local`

### Error: "SANITY_API_TOKEN is required for server-side operations"
**Solution**: Create API token (see instructions above) and add to `.env.local`

### Environment variables not loading
**Solution**:
- Restart dev server after modifying `.env.local`
- Ensure `.env.local` is in the website root directory (not `/sanity/`)
- For Vite, only variables prefixed with `VITE_` are exposed to client-side code

### Import errors
**Solution**: Ensure path aliases are configured correctly in `tsconfig.json`:
```json
{
  "compilerOptions": {
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

## Support

For issues specific to Sanity:
- Documentation: https://www.sanity.io/docs
- Community: https://slack.sanity.io
- Project Dashboard: https://sanity.io/manage/project/z6sd4zry
