package md

// TagFrontMatter is the frontmatter of a tag markdown document.
type TagFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
}
