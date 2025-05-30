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
	// EmpEmployment is a pointer to an Employment.
	EmpEmployment = new(Employment)
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
	// EmpEmploymentToTag is a pointer to an EmploymentToTag.
	EmpEmploymentToTag = new(EmploymentToTag)
	// EmpEmploymentToPost is a pointer to an EmploymentToPost.
	EmpEmploymentToPost = new(EmploymentToPost)
	// EmpEmploymentToProject is a pointer to an EmploymentToProject.
	EmpEmploymentToProject = new(EmploymentToProject)
	// EmpEmploymentToEmployment is a pointer to an EmploymentToEmployment.
	EmpEmploymentToEmployment = new(EmploymentToEmployment)
)

var models = []any{
	EmpPostToTag,
	EmpPostToPost,
	EmpPostToProject,
	EmpProjectToTag,
	EmpProjectToProject,
	EmpTagToTag,
	EmpEmploymentToTag,
	EmpEmploymentToPost,
	EmpEmploymentToProject,
	EmpEmploymentToEmployment,
	EmpPost,
	EmpTag,
	EmpProject,
	EmpEmployment,
	EmpCache,
}

// InitDB initializes the database.
func InitDB(
	ctx context.Context,
	db *bun.DB,
) error {
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
		(*PostToTag)(nil),
		(*PostToPost)(nil),
		(*PostToProject)(nil),
		(*ProjectToTag)(nil),
		(*ProjectToProject)(nil),
		(*TagToTag)(nil),
		(*EmploymentToTag)(nil),
		(*EmploymentToPost)(nil),
		(*EmploymentToProject)(nil),
		(*EmploymentToEmployment)(nil),
		(*Post)(nil),
		(*Tag)(nil),
		(*Project)(nil),
		(*Employment)(nil),
		(*Cache)(nil),
	)
}
