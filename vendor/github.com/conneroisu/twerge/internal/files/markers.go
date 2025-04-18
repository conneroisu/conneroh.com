package files

import (
	"bytes"
	"fmt"
)

// ReplaceBetweenMarkers replaces the content between two markers with the given replacement.
// It takes the content and the replacement content as byte slices. Additionally, it takes
// the markers as strings.
//
// The markers are expected to be at the beginning and end of a line.
func ReplaceBetweenMarkers(
	content, replacement []byte,
	beginMarker, endMarker string,
) ([]byte, error) {
	// Find begin marker
	beginMarkerBytes := []byte(beginMarker)
	beginIdx := bytes.Index(content, beginMarkerBytes)
	if beginIdx == -1 {
		// Markers don't exist, append content with markers
		suffix := append([]byte("\n\n"), beginMarkerBytes...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, replacement...)
		suffix = append(suffix, '\n')
		suffix = append(suffix, []byte(endMarker)...)
		return append(content, suffix...), nil
	}

	// Find the end of the line containing the begin marker
	beginLineEnd := beginIdx + len(beginMarkerBytes)
	for beginLineEnd < len(content) && content[beginLineEnd] != '\n' && content[beginLineEnd] != '\r' {
		beginLineEnd++
	}
	if beginLineEnd < len(content) {
		beginLineEnd++ // Include the newline character
	}

	// Find end marker
	endMarkerBytes := []byte(endMarker)
	endIdx := bytes.Index(content[beginLineEnd:], endMarkerBytes)
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
