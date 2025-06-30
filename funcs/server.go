package funcs

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Call the API
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to get API data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON into struct
	var apiData []Artists
	if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}

	excludeIDs := []int{11, 12, 21, 22, 49}
	filtered := FilterArtists(apiData, excludeIDs)
	RenderTemplate(w, "index.html", filtered)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		RenderTemplate(w, "error.html", nil)
		return
	}

	query := r.URL.Query().Get("searchQuary")
	if query != "" {
		ArtistSearchHandler(w, r) // call search handler
	} else {
		Handler(w, r) // call normal home handler
	}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		// Render custom 404 page
		w.WriteHeader(http.StatusNotFound)
		t404, _ := template.ParseFiles("templates/error.html")
		t404.Execute(w, nil)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError) // Template execution fails 500
	}
}
