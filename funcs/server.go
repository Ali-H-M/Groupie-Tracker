package funcs

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

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

	// Take form locations API , (needed for location Auto complete)
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

	filtered := FilterArtists(apiData, ExcludeIDs)

	// Pagentation Logic
	homeData, ok := Pagentation(r, filtered, Items_Per_Page)
	if !ok { // Invalid page number
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", ErrorData{ErrorName: "404 Page not found", ErrorCode: 404})
		return
	}

	// Auto Complete logic
	suggestions := AutoComplete(filtered, location.Index)
	homeData.Suggestions = suggestions

	RenderTemplate(w, "index.html", homeData)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", ErrorData{ErrorName: "404 Page not found", ErrorCode: 404})
		return
	}

	// Get query parameters
	queryParams := r.URL.Query()
	if len(queryParams) == 0 {
		Handler(w, r)
		return
	}

	// Check if query parameters is allowed
	for key := range r.URL.Query() {
		if !validQueries[key] {
			w.WriteHeader(http.StatusBadRequest)
			RenderTemplate(w, "error.html", ErrorData{ErrorName: "400 Bad Request", ErrorCode: 400})
			return
		}
	}

	// Handle search queries if they exist
	HandleQueries(w, r)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {

	joinMembersList := template.FuncMap{
		"join": strings.Join,
	}

	t, err := template.New(tmpl).Funcs(joinMembersList).ParseFiles("templates/" + tmpl)
	if err != nil {
		// Render custom 404 page
		w.WriteHeader(http.StatusInternalServerError)
		t404, _ := template.ParseFiles("templates/error.html")
		t404.Execute(w, ErrorData{ErrorName: "500 Internal Server Error", ErrorCode: 500})
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		// Template execution error
		w.WriteHeader(http.StatusNotFound)
		t500, _ := template.ParseFiles("templates/error.html")
		t500.Execute(w, ErrorData{ErrorName: "404 Page not found", ErrorCode: 400})
		return
	}
}
