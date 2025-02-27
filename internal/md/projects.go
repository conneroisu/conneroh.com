package md

// ProjectFrontMatter is the frontmatter of a project markdown document.
type ProjectFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
	Slug        string   `yaml:"slug" validate:"required"`
}
