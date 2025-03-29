---
name: New Ideology Category/Folder
about: Request or remember to create a new folder for ideologies
title: "[IDEA-Category]"
labels: ''
assignees: conneroisu

---

## Category Name
<!-- Enter the name of the ideology category/folder -->

## Description
<!-- Provide a brief description of what this category encompasses -->

## Examples of Ideologies to Include
<!-- List some example ideologies that would belong in this category -->

## Parent Category (if applicable)
<!-- If this should be a subfolder of an existing category, specify which one -->

## Folder Structure
The folder should be created at: `internal/data/docs/ideologies/[category-slug]/`

## Index File
Consider creating an index.md file with the following frontmatter:
```yaml
---
title: [Category Name]
slug: [category-slug]
description: [Brief category description]
created_at: [Current date]
updated_at: [Current date]
icon: "folder" # or another appropriate icon
tags:
  - [related-tag-1] 
  - [related-tag-2]
---

# [Category Name] Ideologies

This folder contains documents related to [purpose/description of the category].

## Included Ideologies

- [List will be populated as content is added]
