// Package master contains the master schema for the database.
package master

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
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

// FullPostsList returns all posts with all their projects and tags.
func (q *Queries) FullPostsList(
	ctx context.Context,
	posts []Post,
) (*[]FullPost, error) {
	var (
		fullPosts    []FullPost
		postProjects []Project
		postTags     []Tag
		post         Post
		err          error
	)
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
	posts []Post,
) (*map[string]FullPost, error) {
	var (
		postsMap     = make(map[string]FullPost)
		postProjects []Project
		postTags     []Tag
		post         Post
		err          error
	)
	for _, post = range posts {
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
	projects []Project,
) (*[]FullProject, error) {
	var fullProjects []FullProject
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
	projects []Project,
) (*map[string]FullProject, error) {
	var (
		fullProjects = make(map[string]FullProject)
		projectTags  []Tag
		projectPosts []Post
		project      Project
		err          error
	)
	for _, project = range projects {
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
	tags []Tag,
) (*[]FullTag, error) {
	var fullTags []FullTag
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
	tags []Tag,
) (*map[string]FullTag, error) {
	var (
		fullTags    = make(map[string]FullTag)
		tagProjects []Project
		tagPosts    []Post
		tag         Tag
		err         error
	)
	for _, tag = range tags {
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

// UpsertProjectTags upserts the project tags for a project.
func (q *Queries) UpsertProjectTags(
	ctx context.Context,
	tags []string, // slugs
	id int64,
) error {
	for _, tag := range tags {
		t, err := q.TagGetBySlug(ctx, tag)
		if err != nil {
			return fmt.Errorf("failed to get tag with slug %s: %w", tag, err)
		}
		_, err = q.ProjectTagsGet(ctx, id, t.ID)
		if err == nil {
			return nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get project tag: %w", err)
		}
		err = q.ProjectTagCreate(ctx, id, t.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpsertPostTags upserts the post tags for a post.
func (q *Queries) UpsertPostTags(
	ctx context.Context,
	tags []string, // slugs
	id int64,
) (err error) {
	var tag string
	var t Tag
	for _, tag = range tags {
		t, err = q.TagGetBySlug(ctx, tag)
		if err != nil {
			return fmt.Errorf("failed to get by slug post tags: %s:	%w", tag, err)
		}
		_, err = q.PostTagGet(ctx, id, t.ID)
		if err == nil {
			return nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to upsert post tags (not no rows): %w", err)
		}
		err = q.PostTagCreate(ctx, id, t.ID)
		if err != nil {
			return fmt.Errorf("failed to create post tags: %w", err)
		}
	}
	return nil
}
