package conneroh

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/stretchr/testify/assert"
	// "github.com/uptrace/bun" // Not strictly needed if we don't mock db and set caches directly
)

// Mock data setup
var (
	mockPost1 = &assets.Post{
		Title:       "Test Post Alpha",
		Slug:        "test-post-alpha",
		Description: "Contains keyword_one and a bit about Go",
		Content:     "More details on Go programming.",
		Tags:        []*assets.Tag{mockTag2}, // Link to another tag
	}
	mockPost2 = &assets.Post{
		Title:       "Another Example",
		Slug:        "another-example",
		Description: "No special keywords here.",
		Content:     "Some generic content.",
	}
	mockProject1 = &assets.Project{
		Title:       "Test Project Beta",
		Slug:        "test-project-beta",
		Description: "A project that has keyword_two.",
		Content:     "More about keyword_two.",
	}
	mockTag1 = &assets.Tag{
		Title:       "Test Tag Gamma",
		Slug:        "test-tag-gamma",
		Description: "Includes keyword_one and is about general tech.",
		Content:     "Content for Gamma tag",
	}
	mockTag2 = &assets.Tag{
		Title:       "Go Programming",
		Slug:        "go-programming",
		Description: "Tag for Go related posts",
	}
	mockEmployment1 = &assets.Employment{
		Title:       "Job Epsilon",
		Slug:        "job-epsilon",
		Description: "Role focused on keyword_three development.",
		Content:     "Details about keyword_three.",
	}
)

func setupGlobalSearchTestData() {
	allPosts = []*assets.Post{mockPost1, mockPost2}
	allProjects = []*assets.Project{mockProject1}
	allTags = []*assets.Tag{mockTag1, mockTag2}
	allEmployments = []*assets.Employment{mockEmployment1}
}

func clearGlobalSearchTestData() {
	allPosts = []*assets.Post{}
	allProjects = []*assets.Project{}
	allTags = []*assets.Tag{}
	allEmployments = []*assets.Employment{}
}

func TestHandleGlobalSearch(t *testing.T) {
	// The DB instance can be nil because we are manually setting the global cache slices
	// and the handler checks if len(slice) == 0 before DB access.
	handler := HandleGlobalSearch(nil)

	tests := []struct {
		name               string
		query              string
		expectedStatusCode int
		expectedInBody     []string
		notExpectedInBody  []string
		setupFunc          func()
		cleanupFunc        func()
	}{
		{
			name:               "Query matching multiple types",
			query:              "keyword_one",
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"Test Post Alpha", "Test Tag Gamma"},
			notExpectedInBody:  []string{"Test Project Beta"},
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
		{
			name:               "Query matching specific title",
			query:              "Test Project Beta",
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"Test Project Beta"},
			notExpectedInBody:  []string{"Test Post Alpha"},
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
		{
			name:               "Query matching content (description)",
			query:              "keyword_two", // in project description
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"Test Project Beta"},
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
		{
			name:               "Query matching content (post content)",
			query:              "Go programming", // in post content and linked tag
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"Test Post Alpha", "Go Programming"}, // Post title and Tag title
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
		{
			name:               "Case-insensitive search",
			query:              "KEYWORD_THREE", // in employment description
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"Job Epsilon"},
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
		{
			name:               "Case-insensitive search for title",
			query:              "test post alpha",
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"Test Post Alpha"},
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
		{
			name:               "No results",
			query:              "nonexistentqueryxyz123",
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"No results found for: nonexistentqueryxyz123"},
			notExpectedInBody:  []string{"Test Post Alpha"},
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
		{
			name:               "Empty query",
			query:              "",
			expectedStatusCode: http.StatusOK,
			expectedInBody:     []string{"Global Search", "Enter a search query to find posts, projects, tags, or employments."},
			notExpectedInBody:  []string{"Test Post Alpha"}, // Should not list items by default
			setupFunc:          setupGlobalSearchTestData,
			cleanupFunc:        clearGlobalSearchTestData,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
			}
			if tc.cleanupFunc != nil {
				t.Cleanup(tc.cleanupFunc)
			}

			req, err := http.NewRequest("GET", "/search?q="+tc.query, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := handler(w, r)
				assert.NoError(t, err)
			})

			httpHandler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatusCode, rr.Code)

			body := rr.Body.String()
			for _, expected := range tc.expectedInBody {
				assert.True(t, strings.Contains(body, expected), "Expected to find '"+expected+"' in body")
			}
			for _, notExpected := range tc.notExpectedInBody {
				assert.False(t, strings.Contains(body, notExpected), "Did not expect to find '"+notExpected+"' in body")
			}
		})
	}
}
