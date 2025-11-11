import {defineType} from 'sanity'

export const tagSchema = defineType({
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
      validation: (Rule) => Rule.required(),
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
      options: {hotspot: true},
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
      of: [{type: 'reference', to: {type: 'post'}}],
      readOnly: true,
      description: 'Auto-populated from post tags',
    },
    {
      name: 'projects',
      title: 'Projects',
      type: 'array',
      of: [{type: 'reference', to: {type: 'project'}}],
      readOnly: true,
      description: 'Auto-populated from project tags',
    },
    {
      name: 'relatedTags',
      title: 'Related Tags',
      type: 'array',
      of: [{type: 'reference', to: {type: 'tag'}}],
    },
    {
      name: 'employments',
      title: 'Employments',
      type: 'array',
      of: [{type: 'reference', to: {type: 'employment'}}],
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
})
