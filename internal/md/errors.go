package md

// FrontMatterMissingError is returned when the front matter is missing from the markdown file.
type FrontMatterMissingError struct {
}

// Error implements the error interface on FrontMatterMissingError.
func (e FrontMatterMissingError) Error() string {
	return "front matter missing"
}
