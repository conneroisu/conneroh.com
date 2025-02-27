package md

// PostFrontMatter is the frontmatter of a post markdown document.
type PostFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
}
