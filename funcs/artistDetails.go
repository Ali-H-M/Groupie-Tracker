package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type RelationsData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistDetailPage struct {
	Artist   Artists
	Relation RelationsData
}

type ArtistsWithLocation struct {
	Artist   []Artists
	Location Locations
}

type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artists/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artistAPI := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id)
	resp, err := http.Get(artistAPI)
	if err != nil || resp.StatusCode != http.StatusOK {
		RenderTemplate(w, "error.html", nil)
		return
	}
	defer resp.Body.Close()

	var artist Artists
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		http.Error(w, "Failed to get artist relation data", http.StatusInternalServerError)
		return
	}

	// Get the location data (selected.Location is an API itself)
	relResp, err := http.Get(artist.Relations)
	if err != nil || relResp.StatusCode != http.StatusOK {
		RenderTemplate(w, "error.html", nil)
		return
	}
	defer relResp.Body.Close()

	var relData RelationsData
	if err := json.NewDecoder(relResp.Body).Decode(&relData); err != nil {
		http.Error(w, "Failed to get artist relation data", http.StatusInternalServerError)
		return
	}

	data := ArtistDetailPage{
		Artist:   artist,
		Relation: relData,
	}

	// Check if struct id empty
	if data.Artist.ID != 0 {
		RenderTemplate(w, "artistDetails.html", data)
	} else {
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", nil)
	}
}

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

	// Take form relations API
	locResp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil || locResp.StatusCode != http.StatusOK {
		RenderTemplate(w, "error.html", nil)
		return
	}
	defer locResp.Body.Close()

	var location Locations
	if err := json.NewDecoder(locResp.Body).Decode(&location); err != nil {
		http.Error(w, "Failed to fetch Locations", http.StatusInternalServerError)
		return
	}

	data := ArtistsWithLocation{
		Artist:   artists,
		Location: location,
	}

	// Get search query from URL
	query := r.URL.Query().Get("searchQuary")
	// Call search function
	filtered := SearchArtists(query, data)
	// Remove unwanted objects
	filtered = FilterArtists(filtered, ExcludeIDs)

	RenderTemplate(w, "index.html", filtered)
}
