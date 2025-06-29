package funcs

import (
	"strings"
)

// To remove unwanted API objects
func FilterArtists(artists []Artists, excludeIDs []int) []Artists {
	var result []Artists

	for _, artist := range artists {
		exclude := false
		for _, id := range excludeIDs {
			if artist.ID == id {
				exclude = true
				break
			}
		}
		if !exclude {
			result = append(result, artist)
		}
	}
	return result
}

// Filters artists by name
func SearchArtists(query string, data []Artists) []Artists {
	var result []Artists

	// If query is empty, return full list
	if strings.TrimSpace(query) == "" {
		return data
	}

	for _, artist := range data {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			result = append(result, artist)
			continue // Skip the member check to avoid duplication
		}
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
				result = append(result, artist)
				break
			}
		}
	}

	return result
}
