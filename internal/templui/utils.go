package templui

import (
	"fmt"

	"math/rand"

	"github.com/a-h/templ"

	"maps"

	twmerge "github.com/Oudwins/tailwind-merge-go"
)

// TwMerge combines Tailwind classes and handles conflicts.
// Later classes override earlier ones with the same base.
// Example: "bg-red-500 hover:bg-blue-500", "bg-green-500" → "hover:bg-blue-500 bg-green-500"
func TwMerge(classes ...string) string {
	return twmerge.Merge(classes...)
}

// TwIf returns a class if a condition is true, otherwise an empty string
// Example: "bg-red-500", true → "bg-red-500", false → ""
func TwIf(class string, condition bool) string {
	result := templ.KV(class, condition)
	if result.Value {
		return result.Key
	}
	return ""
}

// mergeAttributes merges multiple Attributes into one
func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		maps.Copy(merged, attr)
	}
	return merged
}

// RandomID returns a random ID string
// Example: "id-123456"
func RandomID() string {
	return fmt.Sprintf("id-%d", rand.Intn(1000000))
}
