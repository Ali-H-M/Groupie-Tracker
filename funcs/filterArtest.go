package funcs

import (
	"strconv"
	"strings"
	"time"
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

func SearchArtists(query string, data PageData) []Artists {
	var result []Artists
	query = strings.ToLower(strings.TrimSpace(query))

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

func filterByCreationDate(artists []Artists, start, end int) []Artists {
	var result []Artists

	for _, artist := range artists {
		if (start == 0 || artist.CreationDate >= start) &&
			(end == 0 || artist.CreationDate <= end) {
			result = append(result, artist)
		}
	}
	return result
}

func filterByFirstAlbum(artists []Artists, startStr, endStr string) []Artists {
	var result []Artists

	var start, end time.Time
	var hasStart, hasEnd bool

	if startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			start = t
			hasStart = true
		}
	}

	if endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			end = t
			hasEnd = true
		}
	}

	for _, artist := range artists {
		albumDate, err := time.Parse("2006-01-02", rearrangeDate(artist.FirstAlbum))
		if err != nil {
			continue
		}

		if (!hasStart || !albumDate.Before(start)) &&
			(!hasEnd || !albumDate.After(end)) {
			result = append(result, artist)
		}
	}
	return result
}

// fixes the api date format to match the form's one
func rearrangeDate(str string) string {
	result := str[6:]
	result += str[2:6]
	result += str[:2]
	return result
}

// validate the input date
func isValidDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// check the checkbox and return's the number of members in int
func filterCheckboxURL(key string) int {
	switch key {
	case "cb1":
		return 1
	case "cb2":
		return 2
	case "cb3":
		return 3
	case "cb4":
		return 4
	case "cb5":
		return 5
	case "cb6":
		return 6
	case "cb7":
		return 7
	case "cb8":
		return 8

	}
	return -1
}

// filters the artists based on the slice we made containing the number of artists
func filterByCheckbox(artists []Artists, slice []int) []Artists {
	var result []Artists

	for _, artist := range artists {
		for _, val := range slice {
			if len(artist.Members) == val {
				result = append(result, artist)
			}
		}
	}
	return result
}

func normalizeLocation(s string) string {
	// Lowercase everything
	s = strings.ToLower(s)

	// Replace
	s = strings.ReplaceAll(s, "_", " ")

	s = strings.ReplaceAll(s, ",", "")

	s = strings.ReplaceAll(s, "-", " ")

	// Collapse multiple spaces into one
	s = strings.Join(strings.Fields(s), " ")

	return s
}

func filterByLocation(artists []Artists, locations Locations, query string) []Artists {
	var result []Artists
	query = normalizeLocation(query)

	// Build a lookup map: location ID -> locations
	locMap := make(map[int][]string)
	for _, loc := range locations.Index {
		locMap[loc.ID] = loc.Locations
	}

	for _, artist := range artists {
		if locs, ok := locMap[artist.ID]; ok {
			for _, loc := range locs {
				if strings.Contains(normalizeLocation(loc), query) {
					result = append(result, artist)
					break
				}
			}
		}
	}
	return result
}
