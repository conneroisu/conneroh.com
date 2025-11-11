# Schema Design: Sanity Content Types

**Project**: conneroh.com (z6sd4zry)
**Created**: 2025-11-11
**Content Source**: Published documents
**Update Strategy**: Polling-based

## Overview

This document defines the Sanity schema for the conneroh.com website, mapping from the existing Go data models to Sanity document types. The schema supports:

- Four main content types: Post, Project, Tag, Employment
- Complex many-to-many relationships
- Rich text content with portable text
- Image assets with metadata
- Publishing workflow (published/draft states)
- SEO metadata

## Design Principles

1. **Field Normalization**: Use Sanity's native field types (slug, reference, array)
2. **Relationships**: Leverage Sanity references instead of storing IDs as strings
3. **Rich Content**: Use Portable Text for flexible content formatting
4. **Assets**: Store images as Sanity Image type with alt text and metadata
5. **Validation**: Strict validation rules in schema for data integrity
6. **Single Source of Truth**: All relationships defined in schema, no manual sync needed

## Content Type Schemas

### 1. Post Document Type

```typescript
// sanity/schemaTypes/post.ts
export const postSchema = {
  name: 'post',
  title: 'Post',
  type: 'document',
  fields: [
    {
      name: 'title',
      title: 'Title',
      type: 'string',
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'slug',
      title: 'Slug',
      type: 'slug',
      options: {
        source: 'title',
        maxLength: 96,
      },
      validation: (Rule) => Rule.required().unique(),
    },
    {
      name: 'description',
      title: 'Description',
      type: 'text',
      rows: 3,
      validation: (Rule) => Rule.required().max(500),
    },
    {
      name: 'content',
      title: 'Content',
      type: 'blockContent',
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'banner',
      title: 'Banner Image',
      type: 'image',
      options: { hotspot: true },
      fields: [
        {
          name: 'alt',
          title: 'Alt Text',
          type: 'string',
        },
      ],
    },
    {
      name: 'createdAt',
      title: 'Created At',
      type: 'datetime',
      initialValue: () => new Date().toISOString(),
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'updatedAt',
      title: 'Updated At',
      type: 'datetime',
      initialValue: () => new Date().toISOString(),
    },
    {
      name: 'tags',
      title: 'Tags',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'tag' } }],
    },
    {
      name: 'projects',
      title: 'Related Projects',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'project' } }],
    },
    {
      name: 'relatedPosts',
      title: 'Related Posts',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'post' } }],
      description: 'Link to related/recommended posts',
    },
    {
      name: 'employments',
      title: 'Related Employments',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'employment' } }],
    },
  ],
  preview: {
    select: {
      title: 'title',
      date: 'createdAt',
    },
    prepare(selection) {
      const { date } = selection
      return {
        title: selection.title,
        subtitle: date && new Date(date).toLocaleDateString(),
      }
    },
  },
}
```

### 2. Project Document Type

```typescript
// sanity/schemaTypes/project.ts
export const projectSchema = {
  name: 'project',
  title: 'Project',
  type: 'document',
  fields: [
    {
      name: 'title',
      title: 'Title',
      type: 'string',
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'slug',
      title: 'Slug',
      type: 'slug',
      options: {
        source: 'title',
        maxLength: 96,
      },
      validation: (Rule) => Rule.required().unique(),
    },
    {
      name: 'description',
      title: 'Description',
      type: 'text',
      rows: 3,
      validation: (Rule) => Rule.required().max(500),
    },
    {
      name: 'content',
      title: 'Content',
      type: 'blockContent',
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'banner',
      title: 'Banner Image',
      type: 'image',
      options: { hotspot: true },
      fields: [
        {
          name: 'alt',
          title: 'Alt Text',
          type: 'string',
        },
      ],
    },
    {
      name: 'createdAt',
      title: 'Created At',
      type: 'datetime',
      initialValue: () => new Date().toISOString(),
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'tags',
      title: 'Tags',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'tag' } }],
    },
    {
      name: 'posts',
      title: 'Related Posts',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'post' } }],
    },
    {
      name: 'relatedProjects',
      title: 'Related Projects',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'project' } }],
    },
    {
      name: 'employments',
      title: 'Related Employments',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'employment' } }],
    },
  ],
  preview: {
    select: {
      title: 'title',
      date: 'createdAt',
    },
    prepare(selection) {
      const { date } = selection
      return {
        title: selection.title,
        subtitle: date && new Date(date).toLocaleDateString(),
      }
    },
  },
}
```

### 3. Tag Document Type

```typescript
// sanity/schemaTypes/tag.ts
export const tagSchema = {
  name: 'tag',
  title: 'Tag',
  type: 'document',
  fields: [
    {
      name: 'title',
      title: 'Title',
      type: 'string',
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'slug',
      title: 'Slug',
      type: 'slug',
      options: {
        source: 'title',
        maxLength: 96,
      },
      validation: (Rule) => Rule.required().unique(),
    },
    {
      name: 'description',
      title: 'Description',
      type: 'text',
      rows: 3,
      validation: (Rule) => Rule.required().max(500),
    },
    {
      name: 'content',
      title: 'Content',
      type: 'blockContent',
    },
    {
      name: 'banner',
      title: 'Banner Image',
      type: 'image',
      options: { hotspot: true },
      fields: [
        {
          name: 'alt',
          title: 'Alt Text',
          type: 'string',
        },
      ],
    },
    {
      name: 'icon',
      title: 'Icon',
      type: 'string',
      description: 'Icon name or emoji (e.g., "code", "🎨")',
    },
    {
      name: 'createdAt',
      title: 'Created At',
      type: 'datetime',
      initialValue: () => new Date().toISOString(),
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'posts',
      title: 'Posts',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'post' } }],
      readOnly: true,
      description: 'Auto-populated from post tags',
    },
    {
      name: 'projects',
      title: 'Projects',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'project' } }],
      readOnly: true,
      description: 'Auto-populated from project tags',
    },
    {
      name: 'relatedTags',
      title: 'Related Tags',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'tag' } }],
    },
    {
      name: 'employments',
      title: 'Employments',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'employment' } }],
      readOnly: true,
      description: 'Auto-populated from employment tags',
    },
  ],
  preview: {
    select: {
      title: 'title',
      icon: 'icon',
    },
    prepare(selection) {
      return {
        title: selection.title,
        media: selection.icon || '🏷️',
      }
    },
  },
}
```

### 4. Employment Document Type

```typescript
// sanity/schemaTypes/employment.ts
export const employmentSchema = {
  name: 'employment',
  title: 'Employment',
  type: 'document',
  fields: [
    {
      name: 'title',
      title: 'Title/Position',
      type: 'string',
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'slug',
      title: 'Slug',
      type: 'slug',
      options: {
        source: 'title',
        maxLength: 96,
      },
      validation: (Rule) => Rule.required().unique(),
    },
    {
      name: 'description',
      title: 'Description',
      type: 'text',
      rows: 3,
      validation: (Rule) => Rule.required().max(500),
    },
    {
      name: 'content',
      title: 'Content/Responsibilities',
      type: 'blockContent',
    },
    {
      name: 'banner',
      title: 'Company Logo/Banner',
      type: 'image',
      options: { hotspot: true },
      fields: [
        {
          name: 'alt',
          title: 'Alt Text',
          type: 'string',
        },
      ],
    },
    {
      name: 'createdAt',
      title: 'Start Date',
      type: 'datetime',
      validation: (Rule) => Rule.required(),
    },
    {
      name: 'endDate',
      title: 'End Date',
      type: 'datetime',
      description: 'Leave empty if currently employed',
    },
    {
      name: 'tags',
      title: 'Tags/Skills',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'tag' } }],
    },
    {
      name: 'posts',
      title: 'Related Posts',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'post' } }],
    },
    {
      name: 'projects',
      title: 'Related Projects',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'project' } }],
    },
    {
      name: 'relatedEmployments',
      title: 'Related Positions',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'employment' } }],
    },
  ],
  preview: {
    select: {
      title: 'title',
      date: 'createdAt',
      endDate: 'endDate',
    },
    prepare(selection) {
      const { date, endDate } = selection
      const startStr = date && new Date(date).getFullYear()
      const endStr = endDate ? new Date(endDate).getFullYear() : 'Present'
      return {
        title: selection.title,
        subtitle: `${startStr} - ${endStr}`,
      }
    },
  },
}
```

### 5. Block Content Type (Portable Text)

```typescript
// sanity/schemaTypes/blockContent.ts
export const blockContentSchema = {
  name: 'blockContent',
  title: 'Block Content',
  type: 'array',
  of: [
    {
      title: 'Block',
      type: 'block',
      styles: [
        { title: 'Normal', value: 'normal' },
        { title: 'H1', value: 'h1' },
        { title: 'H2', value: 'h2' },
        { title: 'H3', value: 'h3' },
        { title: 'Quote', value: 'blockquote' },
      ],
      lists: [
        { title: 'Bullet', value: 'bullet' },
        { title: 'Numbered', value: 'number' },
      ],
      marks: {
        decorators: [
          { title: 'Strong', value: 'strong' },
          { title: 'Emphasis', value: 'em' },
          { title: 'Code', value: 'code' },
          { title: 'Underline', value: 'underline' },
          { title: 'Strike', value: 'strike-through' },
        ],
        annotations: [
          {
            title: 'URL',
            name: 'link',
            type: 'object',
            fields: [
              {
                title: 'URL',
                name: 'href',
                type: 'url',
                validation: (Rule) => Rule.uri({ allowRelative: true }),
              },
            ],
          },
          {
            title: 'Internal Link',
            name: 'internalLink',
            type: 'object',
            fields: [
              {
                title: 'Document',
                name: 'reference',
                type: 'reference',
                to: [
                  { type: 'post' },
                  { type: 'project' },
                  { type: 'employment' },
                ],
              },
            ],
          },
        ],
      },
    },
    {
      type: 'image',
      options: { hotspot: true },
      fields: [
        {
          name: 'alt',
          title: 'Alt Text',
          type: 'string',
        },
        {
          name: 'caption',
          title: 'Caption',
          type: 'string',
        },
      ],
    },
  ],
}
```

## GROQ Query Examples

```typescript
// Fetch all posts with populated references
*[_type == "post"] | order(createdAt desc) {
  ...,
  tags[]-> { title, slug },
  projects[]-> { title, slug },
  relatedPosts[]-> { title, slug },
  employments[]-> { title, slug }
}

// Fetch single post by slug
*[_type == "post" && slug.current == $slug][0] {
  ...,
  tags[]-> { title, slug, icon },
  projects[]-> { title, slug },
  relatedPosts[]-> { title, slug },
  employments[]-> { title, slug }
}

// Fetch all projects with tags
*[_type == "project"] | order(createdAt desc) {
  ...,
  tags[]-> { title, slug, icon }
}

// Fetch posts by tag
*[_type == "post" && references(^._id)] | order(createdAt desc)
```

## Implementation Checklist

- [ ] Create Sanity Studio directory structure
- [ ] Define all schema types (post, project, tag, employment)
- [ ] Create blockContent portable text schema
- [ ] Configure CORS for website domain
- [ ] Create API token for content fetching
- [ ] Test GROQ queries against schema
- [ ] Generate TypeScript types with `sanity typegen`
- [ ] Document schema in README
- [ ] Test migrations from Go data models (if applicable)

## Future Enhancements

- Custom input components for rich validation
- Workflow states (draft → review → published)
- Content scheduling for future publication
- Asset optimization and transformations
- Webhook integration for cache invalidation
- Visual editing mode support
