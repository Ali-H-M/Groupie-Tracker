package funcs

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

var ExcludeIDs = []int{11, 12, 21, 22, 49}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
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

	filtered := FilterArtists(apiData, ExcludeIDs)
	RenderTemplate(w, "index.html", filtered)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", nil)
		return
	}

	// Query parameter check
	for key := range r.URL.Query() {
		if key != "searchQuary" {
			w.WriteHeader(http.StatusBadRequest)
			RenderTemplate(w, "error.html", nil)
			return
		}
	}

	query, ok := r.URL.Query()["searchQuary"]

	if ok && strings.TrimSpace(query[0]) != "" {
		ArtistSearchHandler(w, r) //  Valid query
	} else if ok && strings.TrimSpace(query[0]) == "" {
		w.WriteHeader(http.StatusBadRequest) // searchQuary is present but empty
		Handler(w, r)
	} else {
		Handler(w, r) // No searchQuary param at all: (load normal home page)
	}
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {

	joinMembersList := template.FuncMap{
		"join": strings.Join,
	}

	t, err := template.New(tmpl).Funcs(joinMembersList).ParseFiles("templates/" + tmpl)
	if err != nil {
		// Render custom 404 page
		w.WriteHeader(http.StatusNotFound)
		t404, _ := template.ParseFiles("templates/error.html")
		t404.Execute(w, nil)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		// Template execution error
		w.WriteHeader(http.StatusInternalServerError)
		t500, _ := template.ParseFiles("templates/error.html")
		t500.Execute(w, nil)
		return
	}
}
