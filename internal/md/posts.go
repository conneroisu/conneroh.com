package md

// PostFrontMatter is the frontmatter of a post markdown document.
type PostFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
	Slug        string   `yaml:"slug" validate:"required"`
	BannerUrl   string   `yaml:"banner_url" validate:"required"`
}
