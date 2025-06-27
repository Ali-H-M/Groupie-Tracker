package funcs

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r) // 404
		return
	}

	if r.Method == http.MethodGet {

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
		filtered := filterArtists(apiData, excludeIDs)
		RenderTemplate(w, "index.html", filtered)
		return
	}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound) // 404
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError) // Template execution fails 500
	}
}
