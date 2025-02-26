package master

import "context"

// FullProject is a project with all its posts and tags.
type FullProject struct {
	Project
	Posts []Post
	Tags  []Tag
}

// FullProjectsList returns all projects with all their posts and tags.
func (q *Queries) FullProjectsList(
	ctx context.Context,
) (*[]FullProject, error) {
	var fullProjects []FullProject
	projects, err := q.ProjectsList(ctx)
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		projectPosts, err := q.PostsListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		projectTags, err := q.TagsListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		fullProjects = append(fullProjects, FullProject{
			Project: project,
			Posts:   projectPosts,
			Tags:    projectTags,
		})
	}
	return &fullProjects, nil
}

// FullProjectsSlugMapGet returns a map of projects by their slugs.
func (q *Queries) FullProjectsSlugMapGet(
	ctx context.Context,
) (*map[string]FullProject, error) {
	var (
		fullProjects = make(map[string]FullProject)
		projectTags  []Tag
		projectPosts []Post
		project      Project
		projectList  []Project
		err          error
	)
	projectList, err = q.ProjectsList(ctx)
	if err != nil {
		return nil, err
	}
	for _, project = range projectList {
		projectPosts, err = q.PostsListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		projectTags, err = q.TagsListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		fullProjects[project.Slug] = FullProject{
			Project: project,
			Posts:   projectPosts,
			Tags:    projectTags,
		}
	}
	return &fullProjects, nil
}

// FullProjectCreate creates a project with all its posts and tags.
func (q *Queries) FullProjectCreate(
	ctx context.Context,
	project Project,
	projectPosts []Post,
	projectTags []Tag,
) (*FullProject, error) {
	fullProject, err := q.ProjectCreate(ctx, ProjectCreateParams{
		Name:        project.Name,
		Slug:        project.Slug,
		Description: project.Description,
	})
	if err != nil {
		return nil, err
	}
	for _, post := range projectPosts {
		err := q.PostProjectCreate(ctx, post.ID, fullProject.ID)
		if err != nil {
			return nil, err
		}
	}
	for _, tag := range projectTags {
		err := q.ProjectTagCreate(ctx, fullProject.ID, tag.ID)
		if err != nil {
			return nil, err
		}
	}
	return &FullProject{
		Project: fullProject,
		Posts:   projectPosts,
		Tags:    projectTags,
	}, nil
}
