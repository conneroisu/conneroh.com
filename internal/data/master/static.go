// Package master contains the master schema for the database.
package master

import (
	"context"
	_ "embed"

	"github.com/ollama/ollama/api"
)

//go:generate sqlcquash combine
//go:generate sqlc generate
//go:generate rm -f db.go
//go:generate gomarkdoc -o README.md -e .
//go:generate sh -c "cat combined/schema.sql | sqlite3 -batch static.db"
//go:generate rm -f static.db

// Schema is the schema for the database.
//
//go:embed combined/schema.sql
var Schema string

// Seed is the seed for the database.
//
//go:embed combined/seeds.sql
var Seed string

func createEmbed(ctx context.Context, q *Queries, input string) (int64, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return 0, err
	}
	res, err := client.Embed(ctx, &api.EmbedRequest{
		Model: "text-embedding-ada-002",
		Input: input,
	})
	if err != nil {
		return 0, err
	}
	return q.EmbeddingsCreate(ctx, res.Embeddings)
}

type (
	// FullPost is a post with all its projects and tags.
	FullPost struct {
		Post
		Projects []Project
		Tags     []Tag
	}
	// FullProject is a project with all its posts and tags.
	FullProject struct {
		Project
		Posts []Post
		Tags  []Tag
	}
	// FullTag is a tag with all its posts and projects.
	FullTag struct {
		Tag
		Posts    []Post
		Projects []Project
	}
)

// FullProjectCreate creates a project with all its posts and tags.
func (q *Queries) FullProjectCreate(
	ctx context.Context,
	project Project,
	projectPosts []Post,
	projectTags []Tag,
) (*FullProject, error) {
	id, err := createEmbed(ctx, q, project.Name)
	if err != nil {
		return nil, err
	}
	fullProject, err := q.ProjectCreate(ctx, ProjectCreateParams{
		Name:        project.Name,
		Slug:        project.Slug,
		Description: project.Description,
		EmbeddingID: id,
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

// FullPostsCreate creates a new post with all the associated projects and tags.
func (q *Queries) FullPostsCreate(
	ctx context.Context,
	posts Post,
	postProjects []Project,
	postTags []Tag,
) (*FullPost, error) {
	var (
		fullPost Post
		project  Project
		tag      Tag
		err      error
	)
	id, err := createEmbed(ctx, q, posts.Content)
	if err != nil {
		return nil, err
	}
	fullPost, err = q.PostCreate(ctx, PostCreateParams{
		Title:       posts.Title,
		Slug:        posts.Slug,
		Content:     posts.Content,
		Description: posts.Description,
		EmbeddingID: id,
	})
	if err != nil {
		return nil, err
	}
	for _, project = range postProjects {
		err := q.ProjectPostCreate(ctx, fullPost.ID, project.ID)
		if err != nil {
			return nil, err
		}
	}
	for _, tag = range postTags {
		err := q.PostTagCreate(ctx, fullPost.ID, tag.ID)
		if err != nil {
			return nil, err
		}
	}
	return &FullPost{
		Post:     fullPost,
		Projects: postProjects,
		Tags:     postTags,
	}, nil
}

// FullTagCreate creates a tag with all its posts and projects.
func (q *Queries) FullTagCreate(
	ctx context.Context,
	tag Tag,
	tagProjects []Project,
	tagPosts []Post,
) (*FullTag, error) {
	id, err := createEmbed(ctx, q, tag.Name)
	if err != nil {
		return nil, err
	}
	fullTag, err := q.TagCreate(ctx, TagCreateParams{
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
		EmbeddingID: id,
	})
	if err != nil {
		return nil, err
	}
	for _, project := range tagProjects {
		err := q.ProjectTagCreate(ctx, project.ID, fullTag.ID)
		if err != nil {
			return nil, err
		}
	}
	for _, post := range tagPosts {
		err := q.PostTagCreate(ctx, post.ID, fullTag.ID)
		if err != nil {
			return nil, err
		}
	}
	return &FullTag{
		Tag:      fullTag,
		Posts:    tagPosts,
		Projects: tagProjects,
	}, nil
}

// FullPostsList returns all posts with all their projects and tags.
func (q *Queries) FullPostsList(
	ctx context.Context,
) (*[]FullPost, error) {
	var (
		fullPosts    []FullPost
		posts        []Post
		postProjects []Project
		postTags     []Tag
		post         Post
		err          error
	)
	posts, err = q.PostsList(ctx)
	if err != nil {
		return nil, err
	}
	for _, post = range posts {
		postProjects, err = q.ProjectsListByPost(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		postTags, err = q.TagsListByPost(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		fullPosts = append(fullPosts, FullPost{
			Post:     post,
			Projects: postProjects,
			Tags:     postTags,
		})
	}
	return &fullPosts, nil
}

// FullPostsSlugMapGet returns all posts with all their projects and tags
// as a map of slugs to FullPosts.
func (q *Queries) FullPostsSlugMapGet(
	ctx context.Context,
) (*map[string]FullPost, error) {
	var (
		postsMap     = make(map[string]FullPost)
		postProjects []Project
		postTags     []Tag
		postsList    []Post
		post         Post
		err          error
	)
	postsList, err = q.PostsList(ctx)
	if err != nil {
		return nil, err
	}
	for _, post = range postsList {
		postProjects, err = q.ProjectsListByPost(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		postTags, err = q.TagsListByPost(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		postsMap[post.Slug] = FullPost{
			Post:     post,
			Projects: postProjects,
			Tags:     postTags,
		}
	}
	return &postsMap, nil
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

// FullTagsList returns all tags with all their posts and projects.
func (q *Queries) FullTagsList(
	ctx context.Context,
) (*[]FullTag, error) {
	var fullTags []FullTag
	tags, err := q.TagsListAlphabetical(ctx)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		tagProjects, err := q.ProjectsListByTag(ctx, tag.ID)
		if err != nil {
			return nil, err
		}
		tagPosts, err := q.PostsListByTag(ctx, tag.ID)
		if err != nil {
			return nil, err
		}
		fullTags = append(fullTags, FullTag{
			Tag:      tag,
			Posts:    tagPosts,
			Projects: tagProjects,
		})
	}
	return &fullTags, nil
}

// FullTagsSlugMapGet returns a map of tags by their slugs.
func (q *Queries) FullTagsSlugMapGet(
	ctx context.Context,
) (*map[string]FullTag, error) {
	var (
		fullTags    = make(map[string]FullTag)
		tagProjects []Project
		tagPosts    []Post
		tagsList    []Tag
		tag         Tag
		err         error
	)
	tagsList, err = q.TagsListAlphabetical(ctx)
	if err != nil {
		return nil, err
	}
	for _, tag = range tagsList {
		tagProjects, err = q.ProjectsListByTag(ctx, tag.ID)
		if err != nil {
			return nil, err
		}
		tagPosts, err = q.PostsListByTag(ctx, tag.ID)
		if err != nil {
			return nil, err
		}
		fullTags[tag.Slug] = FullTag{
			Tag:      tag,
			Posts:    tagPosts,
			Projects: tagProjects,
		}
	}
	return &fullTags, nil
}
