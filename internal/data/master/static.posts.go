package master

import "context"

// FullPost is a post with all its projects and tags.
type FullPost struct {
	Post
	Projects []Project
	Tags     []Tag
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
	fullPost, err = q.PostCreate(ctx, PostCreateParams{
		Title:       posts.Title,
		Slug:        posts.Slug,
		Content:     posts.Content,
		Description: posts.Description,
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
