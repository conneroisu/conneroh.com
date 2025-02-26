package tags

import (
	"strings"

	"github.com/conneroisu/conneroh.com/internal/data/master"
)

// NestedSort is a helper function to sort tags by section.
func NestedSort(tags []master.Tag) map[string][]master.Tag {
	nested := make(map[string][]master.Tag)
	var sec string
	for _, tag := range tags {
		if strings.Contains(tag.Name, "/") {
			sec = strings.Split(tag.Name, "/")[0]
			if _, ok := nested[sec]; !ok {
				nested[sec] = []master.Tag{}
			}
			nested[sec] = append(nested[sec], tag)
			continue
		}
		sec = "misc"
		if _, ok := nested[sec]; !ok {
			nested[sec] = []master.Tag{}
		}
		nested[sec] = append(nested[sec], tag)
	}
	return nested
}
