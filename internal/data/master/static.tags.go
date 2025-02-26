package master

import "context"

// FullTag is a tag with all its posts and projects.
type FullTag struct {
	Tag
	Posts    []Post
	Projects []Project
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

// FullTagCreate creates a tag with all its posts and projects.
func (q *Queries) FullTagCreate(
	ctx context.Context,
	tag Tag,
	tagProjects []Project,
	tagPosts []Post,
) (*FullTag, error) {
	fullTag, err := q.TagCreate(ctx, TagCreateParams{
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
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
