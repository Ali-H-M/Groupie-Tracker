package funcs

import (
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

	ArtistsAPIData, err := FetchAPI(w, r)
	if err != nil {
		return
	}

	erro := r.ParseForm()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "Error parsing form", BadRequest)
		return
	}

	homeData := PageData{
		Artist: ArtistsAPIData.Artist,
	}

	// Pagentation Logic
	homeData, ok := Pagentation(r, ArtistsAPIData.Artist, Items_Per_Page)

	if !ok { // Invalid page number
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", ErrorData{ErrorName: "404 Page not found", ErrorCode: 404})
		return
	}

	// Auto Complete logic
	suggestions := AutoComplete(ArtistsAPIData.Artist, ArtistsAPIData.Location.Index)
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
	for key := range queryParams {
		if !validQueries[key] {
			w.WriteHeader(http.StatusBadRequest)
			RenderTemplate(w, "error.html", ErrorData{ErrorName: "400 Bad Request", ErrorCode: 400})
			return
		}
	}

	// Handle search queries if they exist
	for key := range queryParams {
		switch key {
		case "searchQuary":
			HandleMainSearchQuerie(w, r)
			return
		case "filterQuerys":
			HandlerFilltersQuerie(w, r)
			return // RETURN to prevent multiple handlers
		}
	}

	// If no handler matched, go to default
	Handler(w, r)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {

	joinMembersList := template.FuncMap{
		"join": strings.Join,
	}

	t, err := template.New(tmpl).Funcs(joinMembersList).ParseFiles("templates/" + tmpl)
	if err != nil {
		// Render custom error page
		w.WriteHeader(http.StatusInternalServerError)
		tError, _ := template.ParseFiles("templates/error.html")
		if tError != nil {
			tError.Execute(w, InternalServerError)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		// Template execution error
		w.WriteHeader(http.StatusInternalServerError)
		t500, _ := template.ParseFiles("templates/error.html")
		if t500 != nil {
			t500.Execute(w, InternalServerError)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
}
