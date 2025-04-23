package assets

import (
	"context"

	"github.com/uptrace/bun"
)

var (
	// EmpPost is a pointer to a Post.
	EmpPost = new(Post)
	// EmpTag is a pointer to a Tag.
	EmpTag = new(Tag)
	// EmpProject is a pointer to a Project.
	EmpProject = new(Project)

	// EmpAsset is a pointer to an Asset.
	EmpAsset = new(Asset)

	// EmpPostToTag is a pointer to a PostToTag.
	EmpPostToTag = new(PostToTag)
	// EmpPostToPost is a pointer to a PostToPost.
	EmpPostToPost = new(PostToPost)
	// EmpPostToProject is a pointer to a PostToProject.
	EmpPostToProject = new(PostToProject)
	// EmpProjectToTag is a pointer to a ProjectToTag.
	EmpProjectToTag = new(ProjectToTag)
	// EmpProjectToProject is a pointer to a ProjectToProject.
	EmpProjectToProject = new(ProjectToProject)
	// EmpTagToTag is a pointer to a TagToTag.
	EmpTagToTag = new(TagToTag)
	// EmpTagToPost is a pointer to a TagToPost.
	EmpTagToPost = new(TagToPost)
	// EmpTagToProject is a pointer to a TagToProject.
	EmpTagToProject = new(TagToProject)
)

// InitDB initializes the database.
func InitDB(
	ctx context.Context,
	db *bun.DB,
) error {
	err := CreateTables(ctx, db)
	if err != nil {
		return err
	}
	RegisterModels(db)
	return nil
}

// CreateTables creates all the necessary tables for the application.
func CreateTables(ctx context.Context, db *bun.DB) error {
	// Create M2M relationship tables
	relationModels := []any{
		EmpPostToTag,
		EmpPostToPost,
		EmpPostToProject,
		EmpProjectToTag,
		EmpProjectToProject,
		EmpTagToTag,
		EmpTagToPost,
		EmpTagToProject,
	}

	for _, model := range relationModels {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	// Create main entity tables
	models := []any{
		EmpPost,
		EmpTag,
		EmpProject,
		EmpAsset,
	}

	for _, model := range models {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterModels registers all the M2M relationship models with Bun.
func RegisterModels(db *bun.DB) {
	db.RegisterModel(
		EmpPostToTag,
		EmpPostToPost,
		EmpPostToProject,
		EmpProjectToTag,
		EmpProjectToProject,
		EmpTagToTag,
		EmpTagToPost,
		EmpTagToProject,
	)
}
