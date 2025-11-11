import {defineType} from 'sanity'

export const projectSchema = defineType({
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
      validation: (Rule) => Rule.required(),
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
      of: [{type: 'reference', to: {type: 'tag'}}],
    },
    {
      name: 'posts',
      title: 'Related Posts',
      type: 'array',
      of: [{type: 'reference', to: {type: 'post'}}],
    },
    {
      name: 'relatedProjects',
      title: 'Related Projects',
      type: 'array',
      of: [{type: 'reference', to: {type: 'project'}}],
    },
    {
      name: 'employments',
      title: 'Related Employments',
      type: 'array',
      of: [{type: 'reference', to: {type: 'employment'}}],
    },
  ],
  preview: {
    select: {
      title: 'title',
      date: 'createdAt',
    },
    prepare(selection) {
      const {date} = selection
      return {
        title: selection.title,
        subtitle: date && new Date(date).toLocaleDateString(),
      }
    },
  },
})
