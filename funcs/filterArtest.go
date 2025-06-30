package funcs

import (
	"strconv"
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
		// Search in artist name
		if strings.Contains(artist.Name, query) {
			result = append(result, artist)
			continue // Skip the member check to avoid duplication
		}

		for _, member := range artist.Members {
			// Search in members
			if strings.Contains(member, query) {
				result = append(result, artist)
				break
			}
		}

		// Search by creation date
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			result = append(result, artist)
			continue
		}

		// Search by first album date
		if strings.Contains(artist.FirstAlbum, query) {
			result = append(result, artist)
			continue
		}

		// TODO: Search by locations
	
	}

	return result
}
