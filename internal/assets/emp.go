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

	// EmpCache is a pointer to a Cache.
	EmpCache = new(Cache)

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
)

// InitDB initializes the database.
func InitDB(
	ctx context.Context,
	db *bun.DB,
) error {
	// Create database tables
	err := CreateTables(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

// CreateTables creates all the necessary tables for the application.
func CreateTables(ctx context.Context, db *bun.DB) error {

	// Create M2M relationship tables after main tables with their specific foreign keys

	// PostToTag
	_, err := db.NewCreateTable().
		Model(EmpPostToTag).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	// PostToPost
	_, err = db.NewCreateTable().
		Model(EmpPostToPost).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	// PostToProject
	_, err = db.NewCreateTable().
		Model(EmpPostToProject).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	// ProjectToTag
	_, err = db.NewCreateTable().
		Model(EmpProjectToTag).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	// ProjectToProject
	_, err = db.NewCreateTable().
		Model(EmpProjectToProject).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	// TagToTag
	_, err = db.NewCreateTable().
		Model(EmpTagToTag).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	// Create main entity tables first (they are referenced by relationship tables)
	models := []any{
		EmpPost,
		EmpTag,
		EmpProject,
		EmpCache,
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
	// Register all models at once to avoid ordering issues
	db.RegisterModel(
		(*Post)(nil),
		(*Tag)(nil),
		(*Project)(nil),
		(*Cache)(nil),
		(*PostToTag)(nil),
		(*PostToPost)(nil),
		(*PostToProject)(nil),
		(*ProjectToTag)(nil),
		(*ProjectToProject)(nil),
		(*TagToTag)(nil),
	)
}
