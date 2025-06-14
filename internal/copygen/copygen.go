// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.

// Package copygen is a code generator for the assets package.
package copygen

import "github.com/conneroisu/conneroh.com/internal/assets"

// ToPost copies a *assets.Doc to a *assets.Post.
func ToPost(tP *assets.Post, fD *assets.Doc) {
	// *assets.Post fields
	tP.Title = fD.Title
	tP.Slug = fD.Slug
	tP.Description = fD.Description
	tP.Content = fD.Content
	tP.BannerPath = fD.BannerPath
	tP.CreatedAt = fD.CreatedAt
	tP.TagSlugs = fD.TagSlugs
	tP.PostSlugs = fD.PostSlugs
	tP.ProjectSlugs = fD.ProjectSlugs
	tP.EmploymentSlugs = fD.EmploymentSlugs
	tP.Tags = fD.Tags
	tP.Projects = fD.Projects
	tP.Employments = fD.Employments
}

// ToProject copies a *assets.Doc to a *assets.Project.
func ToProject(tP *assets.Project, fD *assets.Doc) {
	// *assets.Project fields
	tP.Title = fD.Title
	tP.Slug = fD.Slug
	tP.Description = fD.Description
	tP.Content = fD.Content
	tP.BannerPath = fD.BannerPath
	tP.CreatedAt = fD.CreatedAt
	tP.TagSlugs = fD.TagSlugs
	tP.PostSlugs = fD.PostSlugs
	tP.ProjectSlugs = fD.ProjectSlugs
	tP.EmploymentSlugs = fD.EmploymentSlugs
	tP.Employments = fD.Employments
}

// ToTag copies a *assets.Doc to a *assets.Tag.
func ToTag(tT *assets.Tag, fD *assets.Doc) {
	// *assets.Tag fields
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
	tT.EmploymentSlugs = fD.EmploymentSlugs
	tT.Projects = fD.Projects
	tT.Employments = fD.Employments
}

// ToEmployment copies a *assets.Doc to a *assets.Employment.
func ToEmployment(tE *assets.Employment, fD *assets.Doc) {
	// *assets.Employment fields
	tE.Title = fD.Title
	tE.Slug = fD.Slug
	tE.Description = fD.Description
	tE.Content = fD.Content
	tE.BannerPath = fD.BannerPath
	tE.CreatedAt = fD.CreatedAt
	tE.EndDate = fD.EndDate
	tE.TagSlugs = fD.TagSlugs
	tE.PostSlugs = fD.PostSlugs
	tE.ProjectSlugs = fD.ProjectSlugs
	tE.EmploymentSlugs = fD.EmploymentSlugs
}
