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

func SearchArtists(query string, data ArtistsWithLocation) []Artists {
	var result []Artists
	query = strings.ToLower(strings.TrimSpace(query))

	// If query is empty, return full list
	if query == "" {
		return data.Artist
	}

	for _, artist := range data.Artist {
		// Search in artist name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			result = append(result, artist)
			continue
		}

		// Search in members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
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
	}

	// Search by location (match with artist ID)
	for _, temp := range data.Location.Index {
		for _, loc := range temp.Locations {
			if strings.Contains(strings.ToLower(loc), query) {
				for _, artist := range data.Artist {
					if artist.ID == temp.ID {
						result = append(result, artist)
						break
					}
				}
				break
			}
		}
	}

	return result
}
