package funcs

import (
	"strings"
)

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

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
		}
	}

	return result
}
