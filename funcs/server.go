package funcs

import (
	"encoding/json"
	"html/template"
	"net/http"
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

		RenderTemplate(w, "index.html", apiData)
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
