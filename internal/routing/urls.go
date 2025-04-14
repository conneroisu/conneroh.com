package routing

import (
	"fmt"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
)

func ComputeAllURLs(base string) []string {
	var urls []string
	for _, post := range gen.AllPosts {
		url := fmt.Sprintf("%s/%s", base, post.PagePath())
		urls = append(urls, url)
	}
	for _, project := range gen.AllProjects {
		url := fmt.Sprintf("%s/%s", base, project.PagePath())
		urls = append(urls, url)
	}
	for _, tag := range gen.AllTags {
		url := fmt.Sprintf("%s/%s", base, tag.PagePath())
		urls = append(urls, url)
	}
	return urls
}
