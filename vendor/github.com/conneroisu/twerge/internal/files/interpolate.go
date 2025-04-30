package files

import (
	"bytes"
	"fmt"
)

// InteroplateMarkers replaces content between two markers.
func InteroplateMarkers(
	content, replacement, beginMarker, endMarker []byte,
) ([]byte, error) {
	// Find begin marker
	beginIdx := bytes.Index(content, beginMarker)
	if beginIdx == -1 {
		// Markers don't exist, append content with markers
		suffix := append([]byte("\n\n"), beginMarker...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, replacement...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, endMarker...)
		content = append(content, suffix...)
	}

	// Find the end of the line containing the begin marker
	beginLineEnd := beginIdx + len(beginMarker)
	for beginLineEnd < len(content) &&
		content[beginLineEnd] != '\n' &&
		content[beginLineEnd] != '\r' {
		beginLineEnd++
	}
	if beginLineEnd < len(content) {
		beginLineEnd++ // Include the newline character
	}

	// Find end marker
	endIdx := bytes.Index(content[beginLineEnd:], endMarker)
	if endIdx == -1 {
		return nil, fmt.Errorf("found begin marker but no end marker")
	}

	// Adjust end marker index to be relative to the whole content
	endIdx += beginLineEnd

	// Create new content with replacement
	result := make([]byte, 0, len(content)-(endIdx-beginLineEnd)+len(replacement)+1)
	result = append(result, content[:beginLineEnd]...)
	result = append(result, replacement...)
	result = append(result, '\n')
	result = append(result, content[endIdx:]...)

	return result, nil
}
