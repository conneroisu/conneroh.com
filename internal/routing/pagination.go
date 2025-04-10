package routing

import (
	"strconv"
)

const (
	// MaxListLargeItems is the maximum number of items in a list view.
	MaxListLargeItems = 9
	// MaxListSmallItems is the maximum number of items in a list view.
	MaxListSmallItems = 27

	// MaxMobilePageNumber is the maximum number of pages to display before ... is shown.
	MaxMobilePageNumber = 5
	// MaxDesktopPageNumber is the maximum number of pages to display before ... is shown.
	MaxDesktopPageNumber = 10

	// Ellipsis represents pagination gaps.
	Ellipsis = "..."
)

// GeneratePagination generates a pagination list of page numbers.
func GeneratePagination(currentPage, totalPages, maxDisplay int) []string {
	result := make([]string, 0)
	if totalPages <= maxDisplay {
		// Show all pages if total fits within maxDisplay
		for i := 1; i <= totalPages; i++ {
			result = append(result, strconv.Itoa(i))
		}
		return result
	}

	// Near the start
	if currentPage <= maxDisplay/2 {
		for i := 1; i <= maxDisplay-1; i++ {
			result = append(result, strconv.Itoa(i))
		}
		result = append(result, "...")
		return result
	}

	// Near the end
	if currentPage >= totalPages-(maxDisplay/2) {
		result = append(result, "...")
		for i := totalPages - maxDisplay + 2; i <= totalPages; i++ {
			result = append(result, strconv.Itoa(i))
		}
		return result
	}

	// Middle range
	result = append(result, "...")
	midCount := maxDisplay - 2
	start := currentPage - midCount/2
	end := start + midCount - 1
	for i := start; i <= end; i++ {
		result = append(result, strconv.Itoa(i))
	}
	result = append(result, "...")

	return result
}
