// Package gen contains the data structures used to generate the site's content.
package gen

import assets "github.com/conneroisu/conneroh.com/internal/assets"

func init() {
	// Maps to store references by slug
	tagMap := make(map[string]*assets.Tag)
	projectMap := make(map[string]*assets.Project)
	postMap := make(map[string]*assets.Post)

	// Populate maps for quick lookup
	for _, tag := range AllTags {
		tagMap[tag.Slug] = tag
	}

	for _, project := range AllProjects {
		projectMap[project.Slug] = project
	}

	for _, post := range AllPosts {
		postMap[post.Slug] = post
	}

	// Process Posts
	for _, post := range AllPosts {
		// Add Tag references
		for _, tagSlug := range post.TagSlugs {
			if tag, ok := tagMap[tagSlug]; ok {
				post.Tags = append(post.Tags, tag)
			}
		}

		// Add Project references
		for _, projectSlug := range post.ProjectSlugs {
			if project, ok := projectMap[projectSlug]; ok {
				post.Projects = append(post.Projects, project)
			}
		}

		// Add Post references (if needed)
		for _, postSlug := range post.PostSlugs {
			if referredPost, ok := postMap[postSlug]; ok {
				post.Posts = append(post.Posts, referredPost)
			}
		}
	}

	// Process Projects
	for _, project := range AllProjects {
		// Add Tag references
		for _, tagSlug := range project.TagSlugs {
			if tag, ok := tagMap[tagSlug]; ok {
				project.Tags = append(project.Tags, tag)
			}
		}

		// Add Project references (if needed)
		for _, projectSlug := range project.ProjectSlugs {
			if referredProject, ok := projectMap[projectSlug]; ok {
				project.Projects = append(project.Projects, referredProject)
			}
		}

		// Add Post references
		for _, postSlug := range project.PostSlugs {
			if post, ok := postMap[postSlug]; ok {
				project.Posts = append(project.Posts, post)
			}
		}
	}

	// Process Tags
	for _, tag := range AllTags {
		// Add Tag references
		for _, tagSlug := range tag.TagSlugs {
			if referredTag, ok := tagMap[tagSlug]; ok {
				tag.Tags = append(tag.Tags, referredTag)
			}
		}

		// Add Project references
		for _, projectSlug := range tag.ProjectSlugs {
			if project, ok := projectMap[projectSlug]; ok {
				tag.Projects = append(tag.Projects, project)
			}
		}

		// Add Post references
		for _, postSlug := range tag.PostSlugs {
			if post, ok := postMap[postSlug]; ok {
				tag.Posts = append(tag.Posts, post)
			}
		}
	}
}
