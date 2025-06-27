package funcs

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
func filterArtists(artists []Artists, excludeIDs []int) []Artists {
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
