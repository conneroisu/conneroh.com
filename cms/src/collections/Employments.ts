import type { CollectionConfig } from 'payload'

export const Employments: CollectionConfig = {
  slug: 'employments',
  admin: {
    useAsTitle: 'title',
    defaultColumns: ['title', 'slug', 'createdAt', 'endDate'],
  },
  fields: [
    {
      name: 'title',
      type: 'text',
      required: true,
    },
    {
      name: 'slug',
      type: 'text',
      required: true,
      unique: true,
      admin: {
        position: 'sidebar',
      },
    },
    {
      name: 'description',
      type: 'textarea',
      admin: {
        description: 'Brief description of the employment',
      },
    },
    {
      name: 'content',
      type: 'richText',
      required: true,
    },
    {
      name: 'bannerPath',
      type: 'relationship',
      relationTo: 'media',
      admin: {
        description: 'Company logo or banner image',
      },
    },
    {
      name: 'tags',
      type: 'relationship',
      relationTo: 'tags',
      hasMany: true,
      admin: {
        position: 'sidebar',
      },
    },
    {
      name: 'posts',
      type: 'relationship',
      relationTo: 'posts',
      hasMany: true,
      admin: {
        position: 'sidebar',
      },
    },
    {
      name: 'projects',
      type: 'relationship',
      relationTo: 'projects',
      hasMany: true,
      admin: {
        position: 'sidebar',
      },
    },
    {
      name: 'employments',
      type: 'relationship',
      relationTo: 'employments',
      hasMany: true,
      admin: {
        position: 'sidebar',
        description: 'Related employments',
      },
    },
    {
      name: 'createdAt',
      type: 'date',
      required: true,
      admin: {
        position: 'sidebar',
        date: {
          pickerAppearance: 'dayOnly',
        },
      },
    },
    {
      name: 'endDate',
      type: 'date',
      admin: {
        position: 'sidebar',
        date: {
          pickerAppearance: 'dayOnly',
        },
        description: 'End date (leave empty if current employment)',
      },
    },
  ],
  hooks: {
    beforeChange: [
      ({ data, operation }) => {
        if (operation === 'create') {
          const now = new Date().toISOString()
          data.createdAt = data.createdAt || now
        }
        return data
      },
    ],
  },
}