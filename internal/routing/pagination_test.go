package routing_test

import (
	"testing"

	"github.com/conneroisu/conneroh.com/internal/routing"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name          string
		items         []int
		page          int
		pageSize      int
		expectedItems []int
		expectedPages int
	}{
		{
			name:          "Empty slice",
			items:         []int{},
			page:          1,
			pageSize:      10,
			expectedItems: []int{},
			expectedPages: 0,
		},
		{
			name:          "Single page",
			items:         []int{1, 2, 3, 4, 5},
			page:          1,
			pageSize:      10,
			expectedItems: []int{1, 2, 3, 4, 5},
			expectedPages: 1,
		},
		{
			name:          "Exactly two pages",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			page:          1,
			pageSize:      5,
			expectedItems: []int{1, 2, 3, 4, 5},
			expectedPages: 2,
		},
		{
			name:          "Exactly two pages - page 2",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			page:          2,
			pageSize:      5,
			expectedItems: []int{6, 7, 8, 9, 10},
			expectedPages: 2,
		},
		{
			name:          "Almost two pages",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			page:          1,
			pageSize:      5,
			expectedItems: []int{1, 2, 3, 4, 5},
			expectedPages: 2,
		},
		{
			name:          "Almost two pages - page 2",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			page:          2,
			pageSize:      5,
			expectedItems: []int{6, 7, 8, 9},
			expectedPages: 2,
		},
		{
			name:          "Three full pages",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			page:          1,
			pageSize:      5,
			expectedItems: []int{1, 2, 3, 4, 5},
			expectedPages: 3,
		},
		{
			name:          "Three full pages - page 2",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			page:          2,
			pageSize:      5,
			expectedItems: []int{6, 7, 8, 9, 10},
			expectedPages: 3,
		},
		{
			name:          "Three full pages - page 3",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			page:          3,
			pageSize:      5,
			expectedItems: []int{11, 12, 13, 14, 15},
			expectedPages: 3,
		},
		{
			name:          "Almost three pages",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			page:          1,
			pageSize:      5,
			expectedItems: []int{1, 2, 3, 4, 5},
			expectedPages: 3,
		},
		{
			name:          "Almost three pages - page 3",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			page:          3,
			pageSize:      5,
			expectedItems: []int{11, 12, 13, 14},
			expectedPages: 3,
		},
		{
			name:          "Partial page with edge case - 11 items, size 5",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			page:          1,
			pageSize:      5,
			expectedItems: []int{1, 2, 3, 4, 5},
			expectedPages: 3, // Should be 3 pages (5 + 5 + 1)
		},
		{
			name:          "Zero page size",
			items:         []int{1, 2, 3, 4, 5},
			page:          1,
			pageSize:      0,
			expectedItems: []int{}, // Should return empty slice to avoid division by zero
			expectedPages: 0,       // Should handle this edge case properly
		},
		{
			name:          "Negative page size",
			items:         []int{1, 2, 3, 4, 5},
			page:          1,
			pageSize:      -5,
			expectedItems: []int{}, // Should return empty slice
			expectedPages: 0,       // Should handle this edge case properly
		},
		{
			name:          "Negative page number",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			page:          -1,
			pageSize:      5,
			expectedItems: []int{1, 2, 3, 4, 5}, // Should default to page 1
			expectedPages: 2,
		},
		{
			name:          "Page number too high",
			items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			page:          5,
			pageSize:      5,
			expectedItems: []int{6, 7, 8, 9, 10}, // Should default to the last page (2)
			expectedPages: 2,
		},
		{
			name:          "Break Test",
			items:         []int{1, 2, 3, 4, 5},
			page:          1,
			pageSize:      2,
			expectedItems: []int{1, 2},
			expectedPages: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItems, gotPages := routing.Paginate(tt.items, tt.page, tt.pageSize)

			// Check if total pages count is correct
			if gotPages != tt.expectedPages {
				t.Errorf("paginate() totalPages = %v, want %v", gotPages, tt.expectedPages)
			}

			// Check if returned items match expected
			if !slicesEqual(gotItems, tt.expectedItems) {
				t.Errorf("paginate() items = %v, want %v", gotItems, tt.expectedItems)
			}
		})
	}
}

// Helper function to compare slices
func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
