package funcs

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func ArtistSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Fetch all artists
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var artists []Artists
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		http.Error(w, "Failed to parse artist data", http.StatusInternalServerError)
		return
	}

	// Take form locations API
	locResp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil || locResp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", nil)
		return
	}
	defer locResp.Body.Close()

	var location Locations
	if err := json.NewDecoder(locResp.Body).Decode(&location); err != nil {
		http.Error(w, "Failed to fetch Locations", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Artist:   artists,
		Location: location,
	}

	// Get search query from URL
	query := r.URL.Query().Get("searchQuary")
	data.SearchQuery = query

	// Remove unwanted objects
	filtered := FilterArtists(data.Artist, ExcludeIDs)
	data.Artist = filtered

	// When nothing match user quary search
	if len(filtered) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	// Save Auto Complete logic
	suggestions := AutoComplete(data.Artist, location.Index)

	// Call search function
	data.Artist = SearchArtists(query, data)

	// Apply pagentation to search result
	data, _ = Pagentation(r, data.Artist, 16)

	// Aplay autocomplete
	data.Suggestions = suggestions

	RenderTemplate(w, "index.html", data)
}

func HandleQueries(w http.ResponseWriter, r *http.Request) {
	// Check for searchQuery URL parameter
	searchQuery, ok := r.URL.Query()["searchQuary"]

	if ok && len(searchQuery) > 0 && strings.TrimSpace(searchQuery[0]) != "" {
		ArtistSearchHandler(w, r) //  Valid query
		return
	} else if ok && len(searchQuery) > 0 && strings.TrimSpace(searchQuery[0]) == "" {
		w.WriteHeader(http.StatusBadRequest) // searchQuary is present but empty
		Handler(w, r)
		return
	}

	Handler(w, r) // No searchQuary param at all: (load normal home page)
}

func Pagentation(r *http.Request, items []Artists, itemsPerPage int) (PageData, bool) {
	pageStr := r.URL.Query().Get("page")
	page := 1

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return PageData{}, false // invalid query
		}
		page = p
	}

	totalItems := len(items)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	if page > totalPages {
		return PageData{}, false // page not found
	}

	start := (page - 1) * itemsPerPage // Calculate the first item (index) in the page
	end := start + itemsPerPage
	if end > totalItems { // When cant fill the hole page
		end = totalItems
	}

	data := PageData{}

	// save serch query if exist {
	if _, exists := r.URL.Query()["searchQuary"]; exists {
		query := r.URL.Query().Get("searchQuary")
		data.SearchQuery = query
	}

	data.Artist = items[start:end]
	data.Page = page
	data.TotalPages = totalPages
	data.HasNext = page < totalPages
	data.HasPrev = page > 1
	data.NextPage = page + 1
	data.PrevPage = page - 1

	// Auto Complete logic
	suggestions := AutoComplete(items, nil)
	data.Suggestions = suggestions

	return data, true
}


func AutoComplete(filtered []Artists, locationIndex []LocationIndex) []Suggestion {
    var suggestions []Suggestion
    
    // Add artist-related suggestions
    for _, art := range filtered {
        suggestions = append(suggestions, Suggestion{Value: art.Name, Label: art.Name + " - Brand Name"})
        suggestions = append(suggestions, Suggestion{Value: strconv.Itoa(art.CreationDate), Label: strconv.Itoa(art.CreationDate) + " - Creation Date"})
        suggestions = append(suggestions, Suggestion{Value: art.FirstAlbum, Label: art.FirstAlbum + " - First Album"})
       
        for _, m := range art.Members {
            suggestions = append(suggestions, Suggestion{Value: m, Label: m + " - Member"})
        }
    }
   
    // Add location suggestions
    for _, indexItem := range locationIndex {
        for _, location := range indexItem.Locations {
            suggestions = append(suggestions, Suggestion{Value: location, Label: location + " - Location"})
        }
    }
    
    return suggestions
}