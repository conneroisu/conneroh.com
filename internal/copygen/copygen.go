// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.

// Package copygen is a code generator for the assets package.
package copygen

import "github.com/conneroisu/conneroh.com/internal/assets"

// ToPost copies a *assets.Doc to a *assets.Post.
func ToPost(tP *assets.Post, fD *assets.Doc) {
	// *assets.Post fields
	tP.X = fD.X
	tP.Y = fD.Y
	tP.Z = fD.Z
	tP.Title = fD.Title
	tP.Slug = fD.Slug
	tP.Description = fD.Description
	tP.Content = fD.Content
	tP.BannerPath = fD.BannerPath
	tP.CreatedAt = fD.CreatedAt
	tP.TagSlugs = fD.TagSlugs
	tP.PostSlugs = fD.PostSlugs
	tP.ProjectSlugs = fD.ProjectSlugs
	tP.Tags = fD.Tags
	tP.Projects = fD.Projects
}

// ToProject copies a *assets.Doc to a *assets.Project.
func ToProject(tP *assets.Project, fD *assets.Doc) {
	// *assets.Project fields
	tP.X = fD.X
	tP.Y = fD.Y
	tP.Z = fD.Z
	tP.Title = fD.Title
	tP.Slug = fD.Slug
	tP.Description = fD.Description
	tP.Content = fD.Content
	tP.BannerPath = fD.BannerPath
	tP.CreatedAt = fD.CreatedAt
	tP.TagSlugs = fD.TagSlugs
	tP.PostSlugs = fD.PostSlugs
	tP.ProjectSlugs = fD.ProjectSlugs
}

// ToTag copies a *assets.Doc to a *assets.Tag.
func ToTag(tT *assets.Tag, fD *assets.Doc) {
	// *assets.Tag fields
	tT.X = fD.X
	tT.Y = fD.Y
	tT.Z = fD.Z
	tT.Title = fD.Title
	tT.Slug = fD.Slug
	tT.Description = fD.Description
	tT.Content = fD.Content
	tT.BannerPath = fD.BannerPath
	tT.Icon = fD.Icon
	tT.CreatedAt = fD.CreatedAt
	tT.TagSlugs = fD.TagSlugs
	tT.PostSlugs = fD.PostSlugs
	tT.ProjectSlugs = fD.ProjectSlugs
	tT.Projects = fD.Projects
}

// ToCache copies a *assets.Doc to a *assets.Cache.
func ToCache(tC *assets.Cache, fD *assets.Doc) {
	// *assets.Cache fields
	tC.Path = fD.Path
	tC.Hash = fD.Hash
	tC.X = fD.X
	tC.Y = fD.Y
	tC.Z = fD.Z
}

// CachedToPost copies a *assets.Cache, *assets.Doc to a *assets.Post.
func CachedToPost(tP *assets.Post, fC *assets.Cache, fD *assets.Doc) {
	// *assets.Post fields
	tP.BaseModel = fC.BaseModel
	tP.ID = fC.ID
	tP.X = fC.X
	tP.Y = fC.Y
	tP.Z = fC.Z
	tP.Title = fD.Title
	tP.Slug = fD.Slug
	tP.Description = fD.Description
	tP.Content = fD.Content
	tP.BannerPath = fD.BannerPath
	tP.CreatedAt = fD.CreatedAt
	tP.TagSlugs = fD.TagSlugs
	tP.PostSlugs = fD.PostSlugs
	tP.ProjectSlugs = fD.ProjectSlugs
	tP.Tags = fD.Tags
	tP.Projects = fD.Projects
}

// CachedToProject copies a *assets.Cache, *assets.Doc to a *assets.Project.
func CachedToProject(tP *assets.Project, fC *assets.Cache, fD *assets.Doc) {
	// *assets.Project fields
	tP.BaseModel = fC.BaseModel
	tP.ID = fC.ID
	tP.X = fC.X
	tP.Y = fC.Y
	tP.Z = fC.Z
	tP.Title = fD.Title
	tP.Slug = fD.Slug
	tP.Description = fD.Description
	tP.Content = fD.Content
	tP.BannerPath = fD.BannerPath
	tP.CreatedAt = fD.CreatedAt
	tP.TagSlugs = fD.TagSlugs
	tP.PostSlugs = fD.PostSlugs
	tP.ProjectSlugs = fD.ProjectSlugs
}

// CachedToTag copies a *assets.Cache, *assets.Doc to a *assets.Tag.
func CachedToTag(tT *assets.Tag, fC *assets.Cache, fD *assets.Doc) {
	// *assets.Tag fields
	tT.BaseModel = fC.BaseModel
	tT.ID = fC.ID
	tT.X = fC.X
	tT.Y = fC.Y
	tT.Z = fC.Z
	tT.Title = fD.Title
	tT.Slug = fD.Slug
	tT.Description = fD.Description
	tT.Content = fD.Content
	tT.BannerPath = fD.BannerPath
	tT.Icon = fD.Icon
	tT.CreatedAt = fD.CreatedAt
	tT.TagSlugs = fD.TagSlugs
	tT.PostSlugs = fD.PostSlugs
	tT.ProjectSlugs = fD.ProjectSlugs
	tT.Projects = fD.Projects
}
