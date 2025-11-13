# Sanity CMS Setup Guide

Complete guide for setting up and configuring Sanity CMS in your TanStack Start + Solid.js project.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Environment Configuration](#environment-configuration)
- [Sanity Project Creation](#sanity-project-creation)
- [Schema Setup](#schema-setup)
- [API Token Configuration](#api-token-configuration)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

## Prerequisites

Before you begin, ensure you have:

- Node.js 18+ or Bun 1.0+
- A Sanity.io account (free tier available at [sanity.io](https://www.sanity.io))
- Basic understanding of TypeScript and GROQ queries

## Quick Start

For existing projects with Sanity already configured:

1. **Copy environment variables** from `.env.example` to `.env.local`
2. **Add your Sanity project credentials** (see [Environment Configuration](#environment-configuration))
3. **Create an API token** (see [API Token Configuration](#api-token-configuration))
4. **Start developing** with the [usage examples](./sanity-usage-examples.md)

## Environment Configuration

### Required Environment Variables

Create or update your `.env.local` file in the project root:

```env
# Sanity CMS Configuration
VITE_SANITY_PROJECT_ID=z6sd4zry
VITE_SANITY_DATASET=production

# Server-side API token (create in Sanity dashboard)
SANITY_API_TOKEN=your_token_here
```

### Variable Breakdown

| Variable | Type | Purpose | Required |
|----------|------|---------|----------|
| `VITE_SANITY_PROJECT_ID` | Public | Sanity project identifier | Yes |
| `VITE_SANITY_DATASET` | Public | Dataset name (production/development) | Yes |
| `SANITY_API_TOKEN` | Secret | Server-side authentication token | Server operations only |

**Important Notes:**
- Variables prefixed with `VITE_` are exposed to the browser (client-side)
- `SANITY_API_TOKEN` is server-only and never exposed to the browser
- Never commit `.env.local` to version control (already in `.gitignore`)

## Sanity Project Creation

If you're starting from scratch, follow these steps to create a new Sanity project:

### Step 1: Install Sanity CLI

```bash
# Using Bun
bun add -D sanity@latest

# Or using npm
npm install -g @sanity/cli
```

### Step 2: Create Sanity Project

```bash
# Initialize Sanity Studio in your project
bunx sanity init

# Follow the prompts:
# - Create new project? Yes
# - Project name: My Project
# - Dataset: production
# - Project output path: ./sanity
# - Schema template: Clean project with no predefined schemas
```

This creates:
- `/sanity` directory with Sanity Studio
- `sanity.config.ts` with your project configuration
- `sanity.cli.ts` for CLI commands

### Step 3: Note Your Project ID

After initialization, your project ID appears in the terminal:

```
Success! Created project "My Project" [abc123xyz]
```

Add this project ID to your `.env.local`:

```env
VITE_SANITY_PROJECT_ID=abc123xyz
```

## Schema Setup

Sanity uses schema definitions to structure your content. Here's how to add content types:

### Creating a Schema Type

Create schema files in `/sanity/schemaTypes/`:

**Example: Post Schema** (`/sanity/schemaTypes/post.ts`)

```typescript
import { defineType, defineField } from 'sanity';

export const post = defineType({
  name: 'post',
  title: 'Post',
  type: 'document',
  fields: [
    defineField({
      name: 'title',
      title: 'Title',
      type: 'string',
      validation: Rule => Rule.required(),
    }),
    defineField({
      name: 'slug',
      title: 'Slug',
      type: 'slug',
      options: {
        source: 'title',
        maxLength: 96,
      },
      validation: Rule => Rule.required(),
    }),
    defineField({
      name: 'excerpt',
      title: 'Excerpt',
      type: 'text',
      rows: 4,
    }),
    defineField({
      name: 'content',
      title: 'Content',
      type: 'array',
      of: [{ type: 'block' }],
    }),
    defineField({
      name: 'publishedAt',
      title: 'Published at',
      type: 'datetime',
      initialValue: () => new Date().toISOString(),
    }),
    defineField({
      name: 'tags',
      title: 'Tags',
      type: 'array',
      of: [{ type: 'reference', to: [{ type: 'tag' }] }],
    }),
  ],
  preview: {
    select: {
      title: 'title',
      subtitle: 'publishedAt',
    },
  },
});
```

### Registering Schema Types

Add your schema to `/sanity/schemaTypes/index.ts`:

```typescript
import { post } from './post';
import { project } from './project';
import { tag } from './tag';

export const schemaTypes = [post, project, tag];
```

### Deploying Schema

After creating/updating schemas, deploy them:

```bash
cd sanity
bun sanity deploy
```

Verify deployment:
1. Open Sanity Studio: `http://localhost:3333` (or your studio URL)
2. Check that your content types appear in the sidebar
3. Try creating a sample document

## API Token Configuration

API tokens enable server-side operations (create, update, delete documents).

### Step 1: Navigate to Sanity Dashboard

1. Visit [sanity.io/manage](https://sanity.io/manage)
2. Select your project
3. Click "API" in the left sidebar

### Step 2: Create Token

1. Navigate to the "Tokens" section
2. Click "Add API token"
3. Configure the token:
   - **Name**: "Website Backend" (or any descriptive name)
   - **Permissions**:
     - **Viewer**: Read-only access (recommended for public websites)
     - **Editor**: Read + write access (for CMS features)
     - **Administrator**: Full access (use cautiously)

4. Click "Add token"

### Step 3: Copy and Secure Token

**IMPORTANT**: The token is shown only once!

1. Copy the token immediately
2. Add to `.env.local`:

```env
SANITY_API_TOKEN=sk_your_actual_token_here_abc123xyz
```

3. Store securely (password manager, secure notes)

### Token Security Best Practices

- Never commit tokens to version control
- Use different tokens for development/staging/production
- Rotate tokens periodically (every 90 days recommended)
- Use minimum required permissions (prefer "Viewer" for read-only sites)
- Revoke tokens immediately if compromised

## Troubleshooting

### Environment Variables Not Loading

**Symptom**: `Error: Sanity project ID is required`

**Solutions**:
1. Verify `.env.local` exists in project root (not `/sanity/`)
2. Restart dev server after modifying `.env.local`
3. Check variable names match exactly (case-sensitive)
4. For Vite, only `VITE_` prefixed variables are client-accessible

**Verify environment variables are loaded**:
```typescript
console.log('Project ID:', import.meta.env.VITE_SANITY_PROJECT_ID);
console.log('Dataset:', import.meta.env.VITE_SANITY_DATASET);
```

### Import Path Errors

**Symptom**: `Cannot find module '@/lib/sanity-client'`

**Solutions**:
1. Verify path aliases in `tsconfig.json`:

```json
{
  "compilerOptions": {
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

2. Restart TypeScript server in your IDE
3. Check file exists: `src/lib/sanity-client.ts`

### Server Token Not Working

**Symptom**: `Error: SANITY_API_TOKEN is required for server-side operations`

**Solutions**:
1. Token must NOT have `VITE_` prefix (it's server-only)
2. Verify token is valid in Sanity dashboard
3. Check token permissions match your operation (viewer vs. editor)
4. Ensure you're calling from server context, not browser

### CORS Errors

**Symptom**: `Access to fetch at 'https://z6sd4zry.api.sanity.io/...' has been blocked by CORS`

**Solutions**:
1. Add your domain to CORS origins in Sanity dashboard:
   - Go to Project Settings → API → CORS Origins
   - Add `http://localhost:3000` for development
   - Add your production domain

2. For authenticated requests, ensure token is sent from server, not browser

### Type Generation Fails

**Symptom**: `Error: Could not connect to Sanity project`

**Solutions**:
1. Verify `sanity.config.ts` has correct project ID
2. Check internet connection
3. Ensure Sanity CLI is authenticated: `bunx sanity login`
4. Try regenerating: `bunx sanity typegen generate`

### Schema Not Updating

**Symptom**: Changes to schema not reflected in Studio

**Solutions**:
1. Deploy schema: `cd sanity && bun sanity deploy`
2. Clear browser cache and hard reload Studio
3. Restart Sanity Studio: `bun sanity dev`
4. Check for syntax errors in schema files

### Content Not Appearing

**Symptom**: GROQ queries return empty arrays

**Solutions**:
1. Verify content is published (not in draft state)
2. Check dataset matches: queries use same dataset as content
3. Test query in Sanity Studio Vision tool
4. Ensure no filters exclude your content (e.g., `!(_id in path("drafts.**"))`)

### TypeScript Errors After Schema Changes

**Symptom**: Type errors after updating Sanity schema

**Solutions**:
1. Regenerate types:
   ```bash
   cd sanity
   bun sanity typegen generate
   cp sanity.types.ts ../src/lib/sanity-types.ts
   ```

2. Restart TypeScript server
3. Update components using changed types
4. Check for breaking changes in field names

## Next Steps

Now that Sanity is configured, explore these resources:

1. **[Usage Examples](./sanity-usage-examples.md)** - Practical code examples for common use cases
2. **[API Reference](./sanity-api-reference.md)** - Complete API documentation
3. **[GROQ Query Examples](https://www.sanity.io/docs/query-cheat-sheet)** - Official GROQ cheat sheet
4. **[Sanity Documentation](https://www.sanity.io/docs)** - Official Sanity documentation

### Development Workflow

1. **Design schema** in `/sanity/schemaTypes/`
2. **Deploy schema**: `cd sanity && bun sanity deploy`
3. **Generate types**: `bun sanity typegen generate`
4. **Create content** in Sanity Studio
5. **Fetch content** using `fetchSanityContent` or `createSanityQuery`
6. **Display content** in your Solid.js components

### Recommended Project Structure

```
project-root/
├── sanity/                    # Sanity Studio
│   ├── schemaTypes/          # Content type definitions
│   ├── sanity.config.ts      # Sanity configuration
│   └── sanity.types.ts       # Auto-generated types
├── src/
│   ├── lib/
│   │   ├── sanity-client.ts  # Client configuration
│   │   ├── sanity-types.ts   # Copied from /sanity
│   │   ├── sanity-queries.ts # GROQ query builders
│   │   ├── sanity-server.ts  # Server-side utilities
│   │   └── sanity-cache.ts   # Cache layer
│   ├── hooks/
│   │   └── use-sanity-content.ts  # Solid.js hooks
│   └── routes/
│       └── posts/
│           ├── index.tsx     # Post listing
│           └── [slug].tsx    # Post detail
└── .env.local                # Environment configuration
```

## Additional Resources

- **Sanity Documentation**: [https://www.sanity.io/docs](https://www.sanity.io/docs)
- **GROQ Reference**: [https://www.sanity.io/docs/groq](https://www.sanity.io/docs/groq)
- **Sanity Community**: [https://slack.sanity.io](https://slack.sanity.io)
- **TypeScript Types**: [Sanity TypeGen Guide](./SANITY_TYPEGEN.md)
- **Project Dashboard**: [https://sanity.io/manage/project/z6sd4zry](https://sanity.io/manage/project/z6sd4zry)

## Support

If you encounter issues not covered in this guide:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review [Usage Examples](./sanity-usage-examples.md)
3. Consult [API Reference](./sanity-api-reference.md)
4. Search [Sanity Community Slack](https://slack.sanity.io)
5. Open an issue in the project repository
