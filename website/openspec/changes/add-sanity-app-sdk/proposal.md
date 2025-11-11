# Add Sanity App SDK Integration

**Change ID**: `add-sanity-app-sdk`

**Status**: Draft

**Created**: 2025-11-11

## Why

Content management is currently hardcoded in the application, which makes it difficult to update without code changes. Integrating Sanity App SDK enables:

1. **Separation of Concerns**: Content (CMS) separate from presentation (frontend code)
2. **Editor-Friendly**: Non-technical editors can manage content via Sanity Studio without touching code
3. **Scalability**: Easy to add new content types and manage growing content volumes
4. **Real-Time Sync**: Updates published in Sanity instantly appear on the website (optional feature)
5. **Ecosystem Integration**: Access to webhooks, transformations, and other Sanity platform features
6. **Type Safety**: Auto-generated TypeScript types ensure compile-time content schema validation

This change establishes the foundation for a modern, CMS-driven website architecture while maintaining full backward compatibility with existing hardcoded content.

## Scope

Integrate the Sanity App SDK into the TanStack Start + Solid.js website to enable:
- Custom content management applications
- Real-time content synchronization
- Headless CMS capabilities with Sanity as the content backend
- Document handles and content fetching patterns
- Type-safe content queries

## Motivation

The website can benefit from Sanity's content management capabilities to:
1. Separate content from presentation logic
2. Enable editors to manage structured content without code changes
3. Provide real-time content updates across the application
4. Leverage Sanity's ecosystem for content operations, webhooks, and integrations

## Capabilities Affected

- **New**: Sanity client setup and initialization
- **New**: Sanity content fetching hooks and utilities
- **New**: Authentication integration with Sanity
- **New**: TypeScript types generation from Sanity schema
- **Modified**: API layer to support content queries
- **Modified**: Environment configuration for Sanity credentials

## Key Design Decisions

1. **Client Library Choice**: Use official `@sanity/client` for server/client content fetching
2. **React Hooks vs Solid.js**: Evaluate compatibility - likely will use vanilla client with Solid.js wrappers
3. **Type Safety**: Leverage Sanity's TypeGen for automatic TypeScript types from schema
4. **Caching Strategy**: Integrate with Solid.js and TanStack infrastructure for efficient caching
5. **SSR Support**: Ensure server-side content fetching works with TanStack Start SSR

## Implementation Phases

1. **Phase 1**: Install and configure Sanity client, set up environment variables
2. **Phase 2**: Create Sanity client wrapper and utilities for content fetching
3. **Phase 3**: Develop Solid.js-compatible hooks for content queries
4. **Phase 4**: Integrate TypeGen for schema-based type generation
5. **Phase 5**: Implement authentication and authorization patterns

## Risks and Mitigations

| Risk | Severity | Mitigation |
|------|----------|-----------|
| Solid.js incompatibility with React-focused SDK | Medium | Wrap SDK in Solid.js primitives; use vanilla client directly |
| Performance regression from additional network calls | Medium | Implement caching, use Solid.js fine-grained reactivity, optimize queries |
| TypeScript type conflicts | Low | Use generated types separately; avoid import conflicts |
| Deployment complexity | Low | Keep configuration in environment variables; document setup steps |

## Acceptance Criteria

- [ ] Sanity client installed and configured
- [ ] Environment variables properly set (project ID, dataset, API token)
- [ ] Content fetching utilities created and tested
- [ ] TypeScript types generated from Sanity schema
- [ ] Integration tests pass for content queries
- [ ] Documentation updated with Sanity setup instructions
- [ ] No breaking changes to existing code

## Related Specs

- TBD (will reference specs created during implementation)

## Dependencies

- External: Sanity account, project, and credentials
- Internal: None (new integration)

## Project Context (Discovered)

**Sanity Project ID**: `z6sd4zry`
**Project Name**: conneroh.com
**Content Source**: Published documents (not drafts)
**Polling Strategy**: Acceptable (real-time sync not needed)
**Performance Constraints**: None specified

## Content Types (Based on Current Data Models)

The following content types will be created/optimized in Sanity:

1. **Post** - Blog post or article
   - title, slug, description, content (rich text)
   - banner_path (image), createdAt, updatedAt
   - Relationships: tags, projects, related posts, employments

2. **Project** - Portfolio or project entry
   - title, slug, description, content (rich text)
   - banner_path (image), createdAt
   - Relationships: tags, posts, related projects, employments

3. **Tag** - Categorization/tagging
   - title, slug, description, content (rich text)
   - banner_path (image), icon
   - createdAt
   - Relationships: posts, projects, related tags, employments

4. **Employment** - Work experience entry
   - title, slug, description, content (rich text)
   - banner_path (image), createdAt, endDate (optional)
   - Relationships: tags, posts, projects, related employments

All types include many-to-many relationships and cross-references as defined in the current Go data models.
