import {defineType} from 'sanity'

export const employmentSchema = defineType({
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
      title: 'Content/Responsibilities',
      type: 'blockContent',
    },
    {
      name: 'banner',
      title: 'Company Logo/Banner',
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
      of: [{type: 'reference', to: {type: 'tag'}}],
    },
    {
      name: 'posts',
      title: 'Related Posts',
      type: 'array',
      of: [{type: 'reference', to: {type: 'post'}}],
    },
    {
      name: 'projects',
      title: 'Related Projects',
      type: 'array',
      of: [{type: 'reference', to: {type: 'project'}}],
    },
    {
      name: 'relatedEmployments',
      title: 'Related Positions',
      type: 'array',
      of: [{type: 'reference', to: {type: 'employment'}}],
    },
  ],
  preview: {
    select: {
      title: 'title',
      date: 'createdAt',
      endDate: 'endDate',
    },
    prepare(selection) {
      const {date, endDate} = selection
      const startStr = date && new Date(date).getFullYear()
      const endStr = endDate ? new Date(endDate).getFullYear() : 'Present'
      return {
        title: selection.title,
        subtitle: `${startStr} - ${endStr}`,
      }
    },
  },
})
