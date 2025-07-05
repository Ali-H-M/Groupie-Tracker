package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artists/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "error.html", nil)
		return
	}

	artistAPI := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id)
	resp, err := http.Get(artistAPI)
	if err != nil || resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusNotFound)
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
		w.WriteHeader(http.StatusNotFound)
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
